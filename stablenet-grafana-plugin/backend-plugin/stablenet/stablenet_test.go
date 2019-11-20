package stablenet

import (
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)
import testify "github.com/stretchr/testify/assert"

func TestClientImpl_QueryDevices(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	devices, err := ioutil.ReadFile("./test-data/devices.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET", "https://127.0.0.1:5443/api/1/devices?$filter=name+ct+%27lab%27", httpmock.NewBytesResponder(200, devices))
	client := NewClient(&ConnectOptions{Port: 5443, Host: "127.0.0.1"})
	clientImpl := client.(*ClientImpl)
	httpmock.ActivateNonDefault(clientImpl.client.GetClient())
	actual, err := client.QueryDevices("lab")
	require.NoError(t, err)

	assert := testify.New(t)
	assert.Equal(1, httpmock.GetTotalCallCount())
	assert.Equal(10, len(actual))
	assert.Equal("newyork.routerlab.infosim.net", actual[7].Name)
}

func TestClientImpl_QueryDevice_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/devices?$filter=name+ct+%27lab%27"
	shouldReturnError := func(client Client) (interface{}, error) {
		return client.QueryDevices("lab")
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, url, "devices matching query \"lab\""))
	t.Run("rest error", errorResponseTest(shouldReturnError, url, "devices matching query \"lab\""))
}

func TestClientImpl_FetchMeasurementsForDevice(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements?$filter=destDeviceId+eq+%271024%27"
	httpmock.Activate()
	defer httpmock.Deactivate()

	rawData, err := ioutil.ReadFile("./test-data/measurements.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, rawData))
	client := NewClient(&ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	clientImpl := client.(*ClientImpl)
	httpmock.ActivateNonDefault(clientImpl.client.GetClient())
	actual, err := client.FetchMeasurementsForDevice(1024)
	require.NoError(t, err)
	require.Equal(t, 10, len(actual), "number of queried measurements wrong")
	test := assert.New(t)
	test.Equal(1587, actual[4].Obid, "obid of fifth measurement wrong")
	test.Equal("Atomcore Processor: 1 ", actual[4].Name, "name of fifth measurement wrong")
}

func TestClientImpl_FetchMeasurementsForDevice_Error(t *testing.T) {
	url := "https://127.0.0.1:5443/api/1/measurements?$filter=destDeviceId+eq+%271024%27"
	shouldReturnError := func(client Client) (interface{}, error) {
		return client.FetchMeasurementsForDevice(1024)
	}
	t.Run("json error", invalidJsonTest(shouldReturnError, url))
	t.Run("status error", wrongStatusResponseTest(shouldReturnError, url, "measurements for device 1024"))
	t.Run("rest error", errorResponseTest(shouldReturnError, url, "measurements for device 1024"))
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
	test := assert.New(t)
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
	url := fmt.Sprintf("https://127.0.0.1:5443/StatisticServlet?stat=1010&type=json&login=infosim,stablenet&id=5555&start=%d&end=%d&value=1&value=2&value=3", start.UnixNano()/int64(time.Millisecond), end.UnixNano()/int64(time.Millisecond))
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
	url := fmt.Sprintf("https://127.0.0.1:5443/StatisticServlet?stat=1010&type=json&login=infosim,stablenet&id=5555&start=%d&end=%d&value=1&value=2&value=3", start.UnixNano()/int64(time.Millisecond), end.UnixNano()/int64(time.Millisecond))
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
