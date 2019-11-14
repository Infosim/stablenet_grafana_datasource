package stablenet

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
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

func TestClientImpl_FetchMeasurementsForDevice(t *testing.T) {
	client := NewClient(ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	measurements, err := client.FetchMeasurementsForDevice(1024)
	require.NoError(t, err)
	assert := testify.New(t)
	assertMeasurementsCorrect(measurements, assert)
}

func assertMeasurementsCorrect(measurements []Measurement, assert *testify.Assertions) {
	assert.Equal(14, len(measurements), "number of measurements incorrect")
	for index, measurement := range measurements {
		assert.NotEmpty(measurement.Name, "name of measurement %d is empty", index+1)
		assert.NotEmpty(measurement.Obid, "obid of measurement %d is empty", index+1)
	}
}

func TestClientImpl_FetchMetricsForMeasurement(t *testing.T) {
	client := NewClient(ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	metrics, err := client.FetchMetricsForMeasurement(6330, time.Now().Add(-48*time.Hour), time.Now())
	require.NoError(t, err)
	assert := testify.New(t)
	assert.Equal(9, len(metrics))
}

func TestClientImpl_unmarshalMeasurements(t *testing.T) {
	file, err := os.Open("./test-data/measurements.xml")
	require.NoError(t, err)
	client := ClientImpl{}
	measurements, err := client.unmarshalMeasurements(file)
	require.NoError(t, err)
	assert := testify.New(t)
	assertMeasurementsCorrect(measurements, assert)
}

func TestClientImpl_QueryDevices(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	devices, err := ioutil.ReadFile("./test-data/devices.json")
	require.NoError(t, err)
	httpmock.RegisterResponder("GET", "https://127.0.0.1:5443/api/1/devices?$filter=name+ct+%27lab%27", httpmock.NewBytesResponder(200, devices))
	client := NewClient(ConnectOptions{Port:5443, Host:"127.0.0.1"})
	clientImpl := client.(*ClientImpl)
	httpmock.ActivateNonDefault(clientImpl.client.GetClient())
	actual, err := client.QueryDevices("lab")
	require.NoError(t, err)

	assert := testify.New(t)
	assert.Equal(1, httpmock.GetTotalCallCount())
	assert.Equal(10, len(actual))
	assert.Equal("newyork.routerlab.infosim.net", actual[7].Name)
}
