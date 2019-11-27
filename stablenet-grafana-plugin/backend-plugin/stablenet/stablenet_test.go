/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package stablenet

import (
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/jarcoal/httpmock"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestClientImpl_QueryDevices(t *testing.T) {
	devices, err := ioutil.ReadFile("./test-data/devices.json")
	require.NoError(t, err)

	tests := []struct {
		name    string
		filter  string
		mockUrl string
	}{
		{name: "no filter", filter: "", mockUrl: "https://127.0.0.1:5443/api/1/devices?$top=100"},
		{name: "one filter", filter: "lab", mockUrl: "https://127.0.0.1:5443/api/1/devices?$top=100&$filter=name+ct+%27lab%27"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.Deactivate()
			httpmock.RegisterResponder("GET", tt.mockUrl, httpmock.NewBytesResponder(200, devices))
			client := NewClient(&ConnectOptions{Port: 5443, Host: "127.0.0.1"})
			clientImpl := client.(*ClientImpl)
			httpmock.ActivateNonDefault(clientImpl.client.GetClient())
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
	url := "https://127.0.0.1:5443/api/1/devices?$top=100&$filter=name+ct+%27lab%27"
	shouldReturnError := func(client Client) (interface{}, error) {
		return client.QueryDevices("lab")
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, url, "devices matching query \"lab\""))
	t.Run("rest error", errorResponseTest(shouldReturnError, url, "devices matching query \"lab\""))
}

func TestClientImpl_FetchMeasurementsForDevice(t *testing.T) {
	rawData, err := ioutil.ReadFile("./test-data/measurements.json")
	require.NoError(t, err)
	deviceId := 1024
	tests := []struct{
		name       string
		deviceObid *int
		nameFilter string
		mockUrl    string
	}{
		{name: "no filter", deviceObid: nil, nameFilter: "", mockUrl: "https://127.0.0.1:5443/api/1/measurements?$top=100"},
		{name: "device filter", deviceObid: &deviceId, nameFilter: "", mockUrl:"https://127.0.0.1:5443/api/1/measurements?$top=100&$filter=destDeviceId+eq+%271024%27"},
		{name: "name filter", deviceObid:nil, nameFilter:"Host", mockUrl:"https://127.0.0.1:5443/api/1/measurements?$top=100&$filter=name+ct+%27Host%27"},
		{name: "device and name filter", deviceObid: &deviceId, nameFilter: "Host", mockUrl:"https://127.0.0.1:5443/api/1/measurements?$top=100&$filter=destDeviceId+eq+%271024%27+and+name+ct+%27Host%27"},
	}
	for _, tt := range tests{
	 t.Run(tt.name, func(t *testing.T) {
		 httpmock.Activate()
		 defer httpmock.Deactivate()
		 httpmock.RegisterResponder("GET", tt.mockUrl, httpmock.NewBytesResponder(200, rawData))
		 client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		 clientImpl := client.(*ClientImpl)
		 httpmock.ActivateNonDefault(clientImpl.client.GetClient())
		 actual, err := client.FetchMeasurementsForDevice(tt.deviceObid, tt.nameFilter)
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
	url := "https://127.0.0.1:5443/api/1/measurements?$top=100&$filter=destDeviceId+eq+%271024%27+and+name+ct+%27Host%27"
	deviceId := 1024
	shouldReturnError := func(client Client) (interface{}, error) {
		return client.FetchMeasurementsForDevice(&deviceId, "Host")
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, url, "measurements for device filter \"destDeviceId eq '1024'\" and name filter \"name ct 'Host'\""))
	t.Run("rest error", errorResponseTest(shouldReturnError, url, "measurements for device filter \"destDeviceId eq '1024'\" and name filter \"name ct 'Host'\""))
}

func TestClientImpl_FetchMetricsForMeasurement(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements/1643/metrics"
	httpmock.Activate()
	defer httpmock.Deactivate()

	rawData, err := ioutil.ReadFile("./test-data/metrics.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, rawData))
	client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	clientImpl := client.(*ClientImpl)
	httpmock.ActivateNonDefault(clientImpl.client.GetClient())
	actual, err := client.FetchMetricsForMeasurement(1643)
	require.NoError(t, err)
	require.Equal(t, 3, len(actual), "number of queried metrics wrong")
	test := testify.New(t)
	test.Equal(1000, actual[0].Id, "id of first metric wrong")
	test.Equal("System Users", actual[0].Name, "name of first metric wrong")
	test.Equal(1001, actual[1].Id, "id of first second wrong")
	test.Equal("System Processes", actual[1].Name, "name of second metric wrong")
	test.Equal(1002, actual[2].Id, "id of third metric wrong")
	test.Equal("System Uptime", actual[2].Name, "name of third metric wrong")
}

func TestClientImpl_FetchMetricsForMeasurement_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements/1643/metrics"
	shouldReturnError := func(client Client) (i interface{}, e error) {
		return client.FetchMetricsForMeasurement(1643)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, url, "metrics for measurement 1643"))
	t.Run("rest error", errorResponseTest(shouldReturnError, url, "metrics for measurement 1643"))
}

func TestClientImpl_FetchDataForMetrics(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Minute)
	url := fmt.Sprintf("https://127.0.0.1:5443/StatisticServlet?stat=1010&type=json&login=infosim,stablenet&id=5555&start=%d&end=%d&value0=1&value1=2&value2=3", start.UnixNano()/int64(time.Millisecond), end.UnixNano()/int64(time.Millisecond))
	httpmock.Activate()
	defer httpmock.Deactivate()

	rawData, err := ioutil.ReadFile("./test-data/measurement-raw-data.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, rawData))
	client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	clientImpl := client.(*ClientImpl)
	httpmock.ActivateNonDefault(clientImpl.client.GetClient())
	actual, err := client.FetchDataForMetrics(5555, []int{1, 2, 3}, start, end)
	require.NoError(t, err)
	systemProcesses := actual["System Processes"]
	systemUsers := actual["System Users"]
	systemUptime := actual["System Uptime"]
	assert := testify.New(t)
	assert.NotNil(systemProcesses, "systemProcesses must not be nil")
	assert.NotNil(systemUsers, "systemUsers must not be nil")
	assert.NotNil(systemUptime, "systemUptime must not be nil")
	assert.Equal(3, len(actual), "number of downloaded metrics")

	var systemUptimeAvg = []*datasource.Point{
		{Timestamp: 1573815483000, Value: 0.207},
		{Timestamp: 1573815783000, Value: 0.210},
		{Timestamp: 1573816083000, Value: 0.214},
		{Timestamp: 1573816383000, Value: 0.217},
		{Timestamp: 1573816683000, Value: 0.221},
		{Timestamp: 1573816983000, Value: 0.224},
		{Timestamp: 1573817283000, Value: 0.228}}
	assert.Equal(systemUptimeAvg, systemUptime.AvgValues(), "system uptime data")
}

func TestClientImpl_FetchDataForMetrics_Error(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Minute)
	url := fmt.Sprintf("https://127.0.0.1:5443/StatisticServlet?stat=1010&type=json&login=infosim,stablenet&id=5555&start=%d&end=%d&value0=1&value1=2&value2=3", start.UnixNano()/int64(time.Millisecond), end.UnixNano()/int64(time.Millisecond))
	shouldReturnError := func(client Client) (i interface{}, e error) {
		return client.FetchDataForMetrics(5555, []int{1, 2, 3}, start, end)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, url, "metric data for measurement 5555"))
	t.Run("rest error", errorResponseTest(shouldReturnError, url, "metric data for measurement 5555"))
}

func invalidJsonTest(shouldReturnError func(Client) (interface{}, error), url string) func(*testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "<>"))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		clientImpl := client.(*ClientImpl)
		httpmock.ActivateNonDefault(clientImpl.client.GetClient())
		_, err := shouldReturnError(client)
		require.EqualError(t, err, "could not unmarshal json: invalid character '<' looking for beginning of value", "error message wrong")
	}
}

func errorResponseTest(shouldReturnError func(Client) (interface{}, error), url string, msg string) func(*testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder("GET", url, httpmock.NewErrorResponder(fmt.Errorf("custom error")))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		clientImpl := client.(*ClientImpl)
		httpmock.ActivateNonDefault(clientImpl.client.GetClient())
		_, err := shouldReturnError(client)
		wantErr := fmt.Sprintf("retrieving %s failed: Get %s: custom error", msg, url)
		require.EqualError(t, err, wantErr, "error message wrong")
	}
}

func wrongStatusResponseTest(shouldReturnError func(Client) (interface{}, error), url string, msg string) func(*testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(404, "entity not found"))
		client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
		clientImpl := client.(*ClientImpl)
		httpmock.ActivateNonDefault(clientImpl.client.GetClient())
		_, err := shouldReturnError(client)
		wantErr := fmt.Sprintf("retrieving %s failed: status code: 404, response: entity not found", msg)
		require.EqualError(t, err, wantErr, "error message wrong")
	}
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
			client := ClientImpl{}
			got := client.formatMetricIds(tt.args)
			testify.Equal(t, tt.want, got)
		})
	}
}

func TestClientImpl_buildJsonApiUrl(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		filters  []string
		want     string
	}{
		{name: "no filters", endpoint: "devices", filters: []string{}, want: "https://127.0.0.1:5443/api/1/devices?$top=100"},
		{name: "two filters", endpoint: "measurement/1234/metrics", filters: []string{"destDeviceId eq '1024'", "name ct 'ether'"}, want: "https://127.0.0.1:5443/api/1/measurement/1234/metrics?$top=100&$filter=destDeviceId+eq+%271024%27+and+name+ct+%27ether%27"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientImpl{
				ConnectOptions: ConnectOptions{
					Host: "127.0.0.1",
					Port: 5443,
				},
			}
			got := c.buildJsonApiUrl(tt.endpoint, tt.filters...)
			require.Equal(t, tt.want, got, "constructed url not correct")
		})
	}
}
