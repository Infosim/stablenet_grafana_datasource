package query

import (
	"backend-plugin/stablenet"
	"bufio"
	"bytes"
	"errors"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_request_stableNetOptionsErrors(t *testing.T) {
	tests := []struct {
		name              string
		jsonData          string
		decryptedJsonData map[string]string
		wantErr           string
	}{
		{name: "invalid json", jsonData: "{", decryptedJsonData: map[string]string{}, wantErr: "could not unmarshal jsonData of the datasource: unexpected end of JSON input"},
		{name: "missing snip", jsonData: "{}", decryptedJsonData: map[string]string{}, wantErr: "the snip is missing in the jsonData of the datasource"},
		{name: "missing snport", jsonData: "{\"snip\":\"127.0.0.1\"}", decryptedJsonData: map[string]string{}, wantErr: "the snport is missing in the jsonData of the datasource"},
		{name: "missing snusername", jsonData: "{\"snip\":\"127.0.0.1\", \"snport\":\"4444\"}", decryptedJsonData: map[string]string{}, wantErr: "the snusername is missing in the jsonData of the datasource"},
		{name: "missing snpassword", jsonData: "{\"snip\":\"127.0.0.1\", \"snport\":\"4444\", \"snusername\":\"infosim\"}", decryptedJsonData: map[string]string{}, wantErr: "the snpassword is missing in the encryptedJsonData of the datasource"},
		{name: "invalid snport", jsonData: "{\"snip\":\"127.0.0.1\", \"snport\": \"hello\", \"snusername\":\"infosim\"}", decryptedJsonData: map[string]string{"snpassword": "stablenet"}, wantErr: "could not parse snport into number: strconv.Atoi: parsing \"hello\": invalid syntax"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				DatasourceRequest: &datasource.DatasourceRequest{
					Datasource: &datasource.DatasourceInfo{
						JsonData:                tt.jsonData,
						DecryptedSecureJsonData: tt.decryptedJsonData,
					},
				},
			}
			got, err := r.stableNetOptions()
			require.Error(t, err, "must return a non-nil error")
			assert.Nil(t, got, "the options must be nil")
			assert.EqualError(t, err, tt.wantErr, "errors do not match")
		})
	}
}

func Test_request_stableNetOptions(t *testing.T) {
	json := "{\"snip\":\"127.0.0.1\", \"snport\": \"443\", \"snusername\":\"infosim\"}"
	decryptedData := map[string]string{"snpassword": "stablenet"}
	request := &Request{
		DatasourceRequest: &datasource.DatasourceRequest{
			Datasource: &datasource.DatasourceInfo{
				JsonData:                json,
				DecryptedSecureJsonData: decryptedData,
			},
		},
	}
	actual, err := request.stableNetOptions()
	require.NoError(t, err, "no error is expected")
	require.NotNil(t, actual, "StableNet Options must not be nil")
	test := assert.New(t)
	test.Equal("127.0.0.1", actual.Host, "host differs")
	test.Equal(443, actual.Port, "port differs")
	test.Equal("infosim", actual.Username, "username differs")
	test.Equal("stablenet", actual.Password, "password differs")
}

func Test_request_timeRange(t *testing.T) {
	now := time.Now()
	then := now.Add(3 * time.Hour)
	nowRaw := now.UnixNano() / int64(time.Millisecond)
	thenRaw := then.UnixNano() / int64(time.Millisecond)
	request := &Request{DatasourceRequest: &datasource.DatasourceRequest{TimeRange: &datasource.TimeRange{FromEpochMs: nowRaw, ToEpochMs: thenRaw}}}
	actualNow, actualThen := request.ToTimeRange()
	assert.Equal(t, now.Second(), actualNow.Second(), "now differs")
	assert.Equal(t, then.Second(), actualThen.Second(), "then differs")
}

func TestDeviceHandler_Process(t *testing.T) {
	rawHandler, _ := setUpHandlerAndLogReceiver()
	devices := []stablenet.Device{
		{Name: "RoGat", Obid: 1024},
		{Name: "localhost", Obid: 1003},
	}
	client := rawHandler.SnClient.(*mockSnClient)
	client.On("QueryDevices", "lab").Return(devices, nil)
	handler := deviceHandler{StableNetHandler: rawHandler}
	result, err := handler.Process(Query{
		Query: datasource.Query{ModelJson: "{\"deviceQuery\":\"lab\"}"},
	})
	assert.NoError(t, err, "no error is expected to be thrown")
	require.NotNil(t, result, "the result must not be nil")
	assert.Equal(t, "[{\"name\":\"RoGat\",\"obid\":1024},{\"name\":\"localhost\",\"obid\":1003}]", result.MetaJson)
}

func TestDeviceHandler_Process_ServerError(t *testing.T) {
	rawHandler, loggedBytes := setUpHandlerAndLogReceiver()
	rawHandler.SnClient.(*mockSnClient).On("QueryDevices", "local").Return(nil, errors.New("internal server error"))
	handler := deviceHandler{StableNetHandler: rawHandler}
	actual, err := handler.Process(Query{
		Query: datasource.Query{ModelJson: "{\"deviceQuery\":\"local\"}"},
	})
	assert.Nil(t, actual, "result should be nil")
	assert.EqualError(t, err, "could not retrieve devices from StableNet(R): internal server error")
	assert.Equal(t, "no time [ERROR] could not retrieve devices from StableNet(R): internal server error\n", loggedBytes.String())
}

func TestHandlersClientErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		json    string
		wantErr string
	}{
		{name: "device query", handler: deviceHandler{}, json: "{}", wantErr: "could not extract the deviceQuery from the query"},
		{name: "measurements for device", handler: measurementHandler{}, json: "{}", wantErr: "could not extract deviceObid from the query"},
		{name: "metrics for measurement", handler: metricNameHandler{}, json: "{}", wantErr: "could not extract measurementObid from query"},
		{name: "metric data", handler: metricDataHandler{}, json: "{}", wantErr: "could not extract measurementObid from query"},
		{name: "metric data no metricId", handler: metricDataHandler{}, json: "{\"measurementObid\": 1626}", wantErr: "could not extract metricId from query"},
		{name: "statisticLinkHandler", handler: statisticLinkHandler{}, json: "{}", wantErr: "could not extract statisticLink parameter from query"},
		{name: "statisticLinkHandler no measurement id", handler: statisticLinkHandler{}, json: "{\"statisticLink\":\"hello\"}", wantErr: "the link \"hello\" does not carry a measurement id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := datasource.Query{ModelJson: tt.json}
			result, err := tt.handler.Process(Query{Query: query})
			assert.NoError(t, err, "on client fails, no error should be returned")
			require.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.wantErr, result.Error, "error message contained in result is wrong")
		})
	}
}

func TestMeasurementHandler_Process(t *testing.T) {
	rawHandler, _ := setUpHandlerAndLogReceiver()
	measurements := []stablenet.Measurement{
		{Name: "Host", Obid: 1124},
		{Name: "CPU", Obid: 1004},
	}
	client := rawHandler.SnClient.(*mockSnClient)
	client.On("FetchMeasurementsForDevice", 1024).Return(measurements, nil)
	handler := measurementHandler{StableNetHandler: rawHandler}
	result, err := handler.Process(Query{
		Query: datasource.Query{ModelJson: "{\"deviceObid\":1024}"},
	})
	assert.NoError(t, err, "no error is expected to be thrown")
	require.NotNil(t, result, "the result must not be nil")
	assert.Equal(t, "[{\"name\":\"Host\",\"obid\":1124},{\"name\":\"CPU\",\"obid\":1004}]", result.MetaJson)
}

func TestMeasurementHandler_Process_ServerError(t *testing.T) {
	rawHandler, loggedBytes := setUpHandlerAndLogReceiver()
	rawHandler.SnClient.(*mockSnClient).On("FetchMeasurementsForDevice", 1024).Return(nil, errors.New("internal server error"))
	handler := measurementHandler{StableNetHandler: rawHandler}
	actual, err := handler.Process(Query{
		Query: datasource.Query{ModelJson: "{\"deviceObid\":1024}"},
	})
	assert.Nil(t, actual, "result should be nil")
	assert.EqualError(t, err, "could not fetch measurements from StableNet(R): internal server error")
	assert.Equal(t, "no time [ERROR] could not fetch measurements from StableNet(R): internal server error\n", loggedBytes.String())
}

func TestDatasourceTestHandler_Process(t *testing.T) {
	rawHandler, _ := setUpHandlerAndLogReceiver()
	rawHandler.SnClient.(*mockSnClient).On("FetchMeasurementsForDevice", -1).Return([]stablenet.Measurement{}, nil)
	handler := datasourceTestHandler{StableNetHandler: rawHandler}
	result, err := handler.Process(Query{})
	assert.NoError(t, err, "no error is expected to be thrown")
	require.NotNil(t, result, "the result must not be nil")
	assert.NotNil(t, result.Series, "the result must contain series")
}

func TestDatasourceTestHandler_Process_Error(t *testing.T) {
	rawHandler, _ := setUpHandlerAndLogReceiver()
	rawHandler.SnClient.(*mockSnClient).On("FetchMeasurementsForDevice", -1).Return(nil, errors.New("login not possible"))
	handler := datasourceTestHandler{StableNetHandler: rawHandler}
	result, err := handler.Process(Query{})
	assert.NoError(t, err, "no error is expected to be thrown")
	require.NotNil(t, result, "the result must not be nil")
	assert.Equal(t, "Cannot login into StableNet(R) with the provided credentials", result.Error)
}

func setUpHandlerAndLogReceiver() (*StableNetHandler, *bytes.Buffer) {
	logData := bytes.Buffer{}
	logReceiver := bufio.NewWriter(&logData)
	logger := hclog.New(&hclog.LoggerOptions{Output: logReceiver, TimeFormat: "no time"})
	return &StableNetHandler{
		SnClient: &mockSnClient{},
		Logger:   logger,
	}, &logData
}

type mockSnClient struct {
	mock.Mock
}

func (m *mockSnClient) QueryDevices(query string) ([]stablenet.Device, error) {
	args := m.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]stablenet.Device), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchMeasurementsForDevice(deviceObid int) ([]stablenet.Measurement, error) {
	args := m.Called(deviceObid)
	if args.Get(0) != nil {
		return args.Get(0).([]stablenet.Measurement), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchMetricsForMeasurement(measurementObid int) ([]stablenet.Metric, error) {
	args := m.Called(measurementObid)
	return args.Get(0).([]stablenet.Metric), args.Error(1)
}

func (m *mockSnClient) FetchDataForMetrics(measurementObid int, metrics []int, start time.Time, end time.Time) (map[string]stablenet.MetricDataSeries, error) {
	args := m.Called(measurementObid, metrics, start, end)
	return args.Get(0).(map[string]stablenet.MetricDataSeries), args.Error(1)
}
