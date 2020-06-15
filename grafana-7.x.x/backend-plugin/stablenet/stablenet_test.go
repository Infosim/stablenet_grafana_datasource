/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package stablenet

import (
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestClientImpl_QueryStableNetVersion(t *testing.T) {
	versionXml := "<info><serverversion version=\"8.6.0\" /></info>"
	tests := []struct {
		name           string
		returnedBody   string
		returnedStatus int
		httpError      error
		wantVersion    *ServerVersion
		wantErrStr     *string
	}{
		{name: "success", returnedBody: versionXml, returnedStatus: http.StatusOK, wantVersion: &ServerVersion{Version: "8.6.0"}, wantErrStr: nil},
		{name: "connection error", httpError: errors.New("server running low on schnitzels"), wantErrStr: strPtr("Connecting to StableNet® failed: Get \"https://127.0.0.1:443/rest/info\": server running low on schnitzels")},
		{name: "authentication error", returnedBody: "Forbidden", returnedStatus: http.StatusUnauthorized, wantErrStr: strPtr("The StableNet® server could be reached, but the credentials were invalid.")},
		{name: "status error", returnedBody: "Internal Server Error", returnedStatus: http.StatusInternalServerError, wantErrStr: strPtr("Log in to StableNet® successful, but the StableNet® version could not be queried. Status Code: 500")},
		{name: "unmarshal error", returnedBody: "this is no xml", returnedStatus: http.StatusOK, wantErrStr: strPtr("Log in to StableNet® successful, but the StableNet® answer \"this is no xml\" could not be parsed: EOF")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.Deactivate()
			if tt.httpError == nil {
				httpmock.RegisterResponder("GET", "https://127.0.0.1:443/rest/info", httpmock.NewStringResponder(tt.returnedStatus, tt.returnedBody))
			} else {
				httpmock.RegisterResponder("GET", "https://127.0.0.1:443/rest/info", httpmock.NewErrorResponder(tt.httpError))
			}
			client := NewClient(&ConnectOptions{Port: 443, Host: "127.0.0.1"})
			httpmock.ActivateNonDefault(client.client.GetClient())
			actual, errStr := client.QueryStableNetVersion()
			testify.Equal(t, tt.wantVersion, actual, "queried server version wrong")
			if tt.wantErrStr != nil {
				testify.Equal(t, *tt.wantErrStr, *errStr, "returned error string wrong")
			} else {
				testify.Nil(t, errStr, "returned error string should be nil")
			}
			httpmock.Reset()
		})
	}
}

func strPtr(value string) *string {
	result := value
	return &result
}

func TestClientImpl_QueryDevices(t *testing.T) {
	devices, err := ioutil.ReadFile("./test-data/devices.json")
	require.NoError(t, err)

	tests := []struct {
		name    string
		filter  string
		mockUrl string
	}{
		{name: "no filter", filter: "", mockUrl: "https://127.0.0.1:5443/api/1/devices?$top=100&$orderBy=name"},
		{name: "one filter", filter: "lab", mockUrl: "https://127.0.0.1:5443/api/1/devices?$top=100&$orderBy=name&$filter=name+ct+%27lab%27"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.Deactivate()
			httpmock.RegisterResponder("GET", tt.mockUrl, httpmock.NewBytesResponder(200, devices))
			client := NewClient(&ConnectOptions{Port: 5443, Host: "127.0.0.1"})
			httpmock.ActivateNonDefault(client.client.GetClient())
			actual, err := client.QueryDevices(tt.filter)
			require.NoError(t, err)

			assert := testify.New(t)
			assert.Equal(1, httpmock.GetTotalCallCount())
			assert.Equal(10, len(actual.Devices))
			assert.Equal("newyork.routerlab.infosim.net", actual.Devices[7].Name)
			assert.True(actual.HasMore)
			httpmock.Reset()
		})
	}

}

func TestClientImpl_QueryDevice_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/devices?$top=100&$orderBy=name&$filter=name+ct+%27lab%27"
	shouldReturnError := func(client *Client) (interface{}, error) {
		return client.QueryDevices("lab")
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, "GET", url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, "GET", url, "devices matching query \"lab\""))
	t.Run("rest error", errorResponseTest(shouldReturnError, "GET", url, "devices matching query \"lab\""))
}

func TestClientImpl_FetchMeasurementsForDevice(t *testing.T) {
	rawData, err := ioutil.ReadFile("./test-data/measurements.json")
	require.NoError(t, err)
	tests := []struct {
		name       string
		deviceObid int
		mockUrl    string
	}{
		{name: "no filter", deviceObid: -1, mockUrl: "https://127.0.0.1:5443/api/1/measurements?$top=100&$orderBy=name&$filter=destDeviceId+eq+%27-1%27"},
		{name: "device filter", deviceObid: 1024, mockUrl: "https://127.0.0.1:5443/api/1/measurements?$top=100&$orderBy=name&$filter=destDeviceId+eq+%271024%27"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.Deactivate()
			httpmock.RegisterResponder("GET", tt.mockUrl, httpmock.NewBytesResponder(200, rawData))
			client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
			httpmock.ActivateNonDefault(client.client.GetClient())
			actual, err := client.FetchMeasurementsForDevice(tt.deviceObid)
			require.NoError(t, err)
			require.Equal(t, 10, len(actual.Measurements), "number of queried measurements wrong")
			test := testify.New(t)
			test.Equal(1587, actual.Measurements[4].Obid, "obid of fifth measurement wrong")
			test.Equal("Atomcore Processor: 1 ", actual.Measurements[4].Name, "name of fifth measurement wrong")
			test.True(actual.HasMore, "hasMore should be true")
		})
	}
}

func TestClientImpl_FetchMeasurementsForDevice_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements?$top=100&$orderBy=name&$filter=destDeviceId+eq+%271024%27"
	shouldReturnError := func(client *Client) (interface{}, error) {
		return client.FetchMeasurementsForDevice(1024)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, "GET", url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, "GET", url, "measurements for device filter \"destDeviceId eq '1024'\""))
	t.Run("rest error", errorResponseTest(shouldReturnError, "GET", url, "measurements for device filter \"destDeviceId eq '1024'\""))
}

func TestClientImpl_FetchMetricsForMeasurement(t *testing.T) {
	rawData, err := ioutil.ReadFile("./test-data/metrics.json")
	require.NoError(t, err)

	mockUrl := "https://127.0.0.1:5443/api/1/measurements/1643/metrics?$top=100"

	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder("GET", mockUrl, httpmock.NewBytesResponder(200, rawData))
	client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	httpmock.ActivateNonDefault(client.client.GetClient())
	metrics, err := client.FetchMetricsForMeasurement(1643)
	require.NoError(t, err)
	require.Equal(t, 3, len(metrics), "number of queried metrics wrong")
	test := testify.New(t)
	test.Equal("SNMP_1000", metrics[0].Key, "Key of first metric wrong")
	test.Equal("System Users", metrics[0].Name, "name of first metric wrong")
	test.Equal("SNMP_1001", metrics[1].Key, "Key of first second wrong")
	test.Equal("System Processes", metrics[1].Name, "name of second metric wrong")
	test.Equal("SNMP_1002", metrics[2].Key, "Key of third metric wrong")
	test.Equal("System Uptime", metrics[2].Name, "name of third metric wrong")
}

func TestClientImpl_FetchMeasurementName(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements?$top=100&$orderBy=name&$filter=obid+eq+%271643%27"
	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "{\"count\": 2264, \"hasMore\": false, \"data\": [{\"name\": \"ThinkStation Host\", \"obid\": 1643}]}"))
	client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	httpmock.ActivateNonDefault(client.client.GetClient())
	name, err := client.FetchMeasurementName(1643)
	require.NoError(t, err, "no error expected")
	require.Equal(t, "ThinkStation Host", *name, "name not correct")
}

func TestClientImpl_FetchMetricsForMeasurement_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements/1643/metrics?$top=100"
	shouldReturnError := func(client *Client) (i interface{}, e error) {
		return client.FetchMetricsForMeasurement(1643)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, "GET", url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, "GET", url, "metrics for measurement 1643"))
	t.Run("rest error", errorResponseTest(shouldReturnError, "GET", url, "metrics for measurement 1643"))
}

func TestClientImpl_FetchDataForMetrics(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Minute)
	url := fmt.Sprintf("https://127.0.0.1:5443/api/1/measurements/5555/data")
	httpmock.Activate()
	defer httpmock.Deactivate()

	rawData, err := ioutil.ReadFile("./test-data/measurement-raw-data.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("POST", url, httpmock.NewBytesResponder(200, rawData))
	client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	httpmock.ActivateNonDefault(client.client.GetClient())
	options := DataQueryOptions{
		MeasurementObid: 5555,
		Metrics:         []string{"System Processes", "System Users", "System Uptime"},
		Start:           start,
		End:             end,
		Average:         250,
	}
	actual, err := client.FetchDataForMetrics(options)
	require.NoError(t, err)
	systemProcesses := actual["System Processes"]
	systemUsers := actual["System Users"]
	systemUptime := actual["System Uptime"]
	assert := testify.New(t)
	assert.NotNil(systemProcesses, "systemProcesses must not be nil")
	assert.NotNil(systemUsers, "systemUsers must not be nil")
	assert.NotNil(systemUptime, "systemUptime must not be nil")
	assert.Equal(3, len(actual), "number of downloaded metrics")

	var systemUptimeAvg = [][]interface{}{
		{time.Unix(0, 1574839083813*int64(time.Millisecond)), 0.207},
		{time.Unix(0, 1574839383813*int64(time.Millisecond)), 0.210},
		{time.Unix(0, 1574839683813*int64(time.Millisecond)), 0.214},
		{time.Unix(0, 1574839983813*int64(time.Millisecond)), 0.217},
		{time.Unix(0, 1574840283813*int64(time.Millisecond)), 0.221},
		{time.Unix(0, 1574840583813*int64(time.Millisecond)), 0.224},
		{time.Unix(0, 1574840883813*int64(time.Millisecond)), 0.228}}
	assert.Equal(systemUptimeAvg, systemUptime.AsTable(false, false, true), "system uptime data")
}

func TestClientImpl_FetchDataForMetrics_Error(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Minute)
	url := fmt.Sprintf("https://127.0.0.1:5443/api/1/measurements/5555/data?top=-1")
	options := DataQueryOptions{
		MeasurementObid: 5555,
		Metrics:         []string{"1", "2", "3"},
		Start:           start,
		End:             end,
	}
	shouldReturnError := func(client *Client) (i interface{}, e error) {
		return client.FetchDataForMetrics(options)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, "POST", url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, "POST", url, "metric data for measurement 5555"))
	t.Run("rest error", errorResponseTest(shouldReturnError, "POST", url, "metric data for measurement 5555"))
}

func invalidJsonTest(shouldReturnError func(*Client) (interface{}, error), method string, url string) func(*testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder(method, url, httpmock.NewStringResponder(200, "<>"))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		httpmock.ActivateNonDefault(client.client.GetClient())
		result, err := shouldReturnError(client)
		testify.Nil(t, result, "the result should be nil")
		require.EqualError(t, err, "could not unmarshal json: invalid character '<' looking for beginning of value", "error message wrong")
	}
}

func errorResponseTest(shouldReturnError func(*Client) (interface{}, error), method string, url string, msg string) func(*testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder(method, url, httpmock.NewErrorResponder(fmt.Errorf("custom error")))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		httpmock.ActivateNonDefault(client.client.GetClient())
		_, err := shouldReturnError(client)
		capitalizedMethod := []byte(strings.ToLower(method))
		capitalizedMethod[0] = byte(method[0])
		wantErr := fmt.Sprintf("retrieving %s failed: %s \"%s\": custom error", msg, capitalizedMethod, url)
		require.EqualError(t, err, wantErr, "error message wrong")
	}
}

func wrongStatusResponseTest(shouldReturnError func(*Client) (interface{}, error), method string, url string, msg string) func(*testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder(method, url, httpmock.NewStringResponder(404, "entity not found"))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		httpmock.ActivateNonDefault(client.client.GetClient())
		_, err := shouldReturnError(client)
		wantErr := fmt.Sprintf("retrieving %s failed: status code: 404, response: entity not found", msg)
		require.EqualError(t, err, wantErr, "error message wrong")
	}
}

func TestClientImpl_FetchMeasurementName_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements?$top=100&$orderBy=name&$filter=obid+eq+%271643%27"
	shouldReturnError := func(client *Client) (i interface{}, e error) {
		return client.FetchMeasurementName(1643)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, "GET", url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, "GET", url, "name for measurement 1643"))
	t.Run("rest error", errorResponseTest(shouldReturnError, "GET", url, "name for measurement 1643"))
	t.Run("no measurement", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "{\"count\": 2264, \"hasMore\": false, \"data\": []}"))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		httpmock.ActivateNonDefault(client.client.GetClient())
		_, err := client.FetchMeasurementName(1643)
		require.EqualError(t, err, "measurement with id 1643 does not exist", "error message wrong")
	})
}

func TestClientImpl_formatMetricIds(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want string
	}{
		{name: "single value", args: []int{123}, want: "value=123"},
		{name: "three values", args: []int{1, 2, 3}, want: "value0=1&value1=2&value2=3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{}
			got := client.formatMetricIds(tt.args)
			testify.Equal(t, tt.want, got)
		})
	}
}

func TestClientImpl_buildJsonApiUrl(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		orderBy  string
		filters  []string
		want     string
	}{
		{name: "no filters", endpoint: "devices", filters: []string{}, want: "https://127.0.0.1:5443/api/1/devices?$top=100"},
		{name: "two filters", endpoint: "measurement/1234/metrics", filters: []string{"destDeviceId eq '1024'", "name ct 'ether'"}, want: "https://127.0.0.1:5443/api/1/measurement/1234/metrics?$top=100&$filter=destDeviceId+eq+%271024%27+and+name+ct+%27ether%27"},
		{name: "two filter with order by", endpoint: "measurement/1234/metrics", orderBy: "description", filters: []string{"destDeviceId eq '1024'", "name ct 'ether'"}, want: "https://127.0.0.1:5443/api/1/measurement/1234/metrics?$top=100&$orderBy=description&$filter=destDeviceId+eq+%271024%27+and+name+ct+%27ether%27"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ConnectOptions: ConnectOptions{
					Host: "127.0.0.1",
					Port: 5443,
				},
			}
			got := c.buildJsonApiUrl(tt.endpoint, tt.orderBy, tt.filters...)
			require.Equal(t, tt.want, got, "constructed url not correct")
		})
	}
}
