package stablenet

import (
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)
import testify "github.com/stretchr/testify/assert"

func assertDevicesCorrect(devices []Device, assert *testify.Assertions) {
	assert.Equal(97, len(devices), "number of devices incorrect")
	for index, device := range devices {
		assert.NotEmpty(device.Name, "name of device %d is empty", index+1)
		assert.NotEmpty(device.Obid, "obid of device %d is empty", index+1)
	}
}

func assertMeasurementsCorrect(measurements []Measurement, assert *testify.Assertions) {
	assert.Equal(14, len(measurements), "number of measurements incorrect")
	for index, measurement := range measurements {
		assert.NotEmpty(measurement.Name, "name of measurement %d is empty", index+1)
		assert.NotEmpty(measurement.Obid, "obid of measurement %d is empty", index+1)
	}
}

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
