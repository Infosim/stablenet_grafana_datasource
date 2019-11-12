package stablenet

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)
import testify "github.com/stretchr/testify/assert"

func TestClientImpl_FetchAllDevices(t *testing.T) {
	client := NewClient(ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	devices, err := client.FetchAllDevices()
	require.NoError(t, err)
	assert := testify.New(t)
	assertDevicesCorrect(devices, assert)
}

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

func TestClientImpl_unmarshalDevices(t *testing.T) {
	file, err := os.Open("./test-data/devices.xml")
	require.NoError(t, err)
	client := ClientImpl{}
	devices, err := client.unmarshalDevices(file)
	require.NoError(t, err)
	assert := testify.New(t)
	assertDevicesCorrect(devices, assert)
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
