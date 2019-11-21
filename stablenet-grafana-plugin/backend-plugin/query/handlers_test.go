package query

import (
	"backend-plugin/stablenet"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetHandlersForRequest(t *testing.T) {
	json := "{\"snip\":\"127.0.0.1\", \"snport\": \"443\", \"snusername\":\"infosim\"}"
	decryptedData := map[string]string{"snpassword": "stablenet"}
	request := Request{
		DatasourceRequest: &datasource.DatasourceRequest{
			Datasource: &datasource.DatasourceInfo{
				JsonData:                json,
				DecryptedSecureJsonData: decryptedData,
			},
		},
	}
	handlers, err := GetHandlersForRequest(request)
	require.NoError(t, err, "no error expected")
	assert.Equal(t, 6, len(handlers), "expected six handlers")
	assert.IsType(t, deviceHandler{}, handlers["devices"], "deviceQuery handler hast not correct type")
	assert.IsType(t, measurementHandler{}, handlers["measurements"], "measurement handler hast not correct type")
	assert.IsType(t, metricNameHandler{}, handlers["metricNames"], "metric name handler hast not correct type")
	assert.IsType(t, metricDataHandler{}, handlers["metricData"], "metric data handler hast not correct type")
	assert.IsType(t, statisticLinkHandler{}, handlers["statisticLink"], "statistic link handler hast not correct type")
	assert.IsType(t, datasourceTestHandler{}, handlers["testDatasource"], "datasource test handler hast not correct type")
}

func TestGetHandlersForRequest_Error(t *testing.T) {
	request := Request{
		DatasourceRequest: &datasource.DatasourceRequest{},
	}
	handlers, err := GetHandlersForRequest(request)
	require.EqualError(t, err, "could not extract StableNet(R) connect options: datasource info is nil", "error message not correct")
	require.Nil(t, handlers, "handlers should be nil")
}

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

func TestHandlersSuccessfulResponse(t *testing.T) {
	var tests = []struct {
		*handlerServerTestCase
		name string
	}{
		{name: "device query", handlerServerTestCase: deviceHandlerTest()},
		{name: "measurement query", handlerServerTestCase: measurementHandlerTest()},
		{name: "metric query", handlerServerTestCase: metricNameHandlerTest()},
		{name: "metric data", handlerServerTestCase: metricDataHandlerTest()},
		{name: "statistic link", handlerServerTestCase: statisticLinkHandlerTest()},
		{name: "datasource test", handlerServerTestCase: datasourceTestHandlerTest()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clientArgs == nil {
				tt.clientArgs = tt.queryArgs
			}
			rawHandler, _ := setUpHandlerAndLogReceiver()
			clientArgs := make([]interface{}, 0, len(tt.clientArgs))
			for _, argument := range tt.clientArgs {
				clientArgs = append(clientArgs, argument.value)
			}
			rawHandler.SnClient.(*mockSnClient).On(tt.clientMethod, clientArgs...).Return(tt.clientReturn, nil)
			handler := tt.handler(rawHandler)
			jsonStructure := make(map[string]interface{})
			for _, argument := range tt.queryArgs {
				jsonStructure[argument.name] = argument.value
			}
			queryData, _ := json.Marshal(jsonStructure)
			actual, err := handler.Process(Query{
				Query: datasource.Query{ModelJson: string(queryData), RefId: "the cake is a lie"},
			})
			assert.NoError(t, err, "no error is expected")
			require.NotNil(t, actual, "actual must not be nil")
			if tt.successResult.MetaJson != "" {
				require.NotEmpty(t, actual.MetaJson, "metaJsonExpected, but got none")
				var actualMetaData map[string]interface{}
				_ = json.Unmarshal([]byte(actual.MetaJson), &actualMetaData)
				var wantMetaData map[string]interface{}
				_ = json.Unmarshal([]byte(tt.successResult.MetaJson), &wantMetaData)
				assert.Equal(t, wantMetaData, actualMetaData, "metaJson does differ")
			}
			assert.Equal(t, "the cake is a lie", actual.RefId, "the refId is wrong")
			if tt.successResult.Series != nil {
				compareTimeSeries(t, tt.successResult.Series, actual.Series)
			}
			assert.Equal(t, tt.successResult.Tables, actual.Tables, "tables differ")
			assert.Empty(t, actual.Error, "no error expected, this is checked in another testcase")
		})
	}
}

func compareTimeSeries(t *testing.T, want []*datasource.TimeSeries, actual []*datasource.TimeSeries) {
	require.NotNil(t, actual, "time series were expected, but not delivered")
	require.Equal(t, len(want), len(actual), "length of time series differes")
	for index, series := range want {
		actualSeries := actual[index]
		assert.Equal(t, series.Name, actualSeries.Name, "name of %dth timeseries differ", index+1)
		sameLength := assert.Equal(t, len(series.Points), len(actualSeries.Points), "number of points in %dth timeseries differ", index+1)
		if sameLength {
			for pIndex, point := range series.Points {
				assert.Equal(t, point, actualSeries.Points[pIndex], "point %d of %dth time series differs", pIndex+1, index+1)
			}
		}
	}
}

func TestHandlersServerErrors(t *testing.T) {
	var tests = []struct {
		*handlerServerTestCase
		name    string
		wantErr string
	}{
		{name: "device query", handlerServerTestCase: deviceHandlerTest(), wantErr: "could not retrieve devices from StableNet(R)"},
		{name: "measurement query", handlerServerTestCase: measurementHandlerTest(), wantErr: "could not fetch measurements from StableNet(R)"},
		{name: "metric query", handlerServerTestCase: metricNameHandlerTest(), wantErr: "could not retrieve metric names from StableNet(R)"},
		{name: "metric data", handlerServerTestCase: metricDataHandlerTest(), wantErr: "could not fetch metric data from server: could not retrieve metrics from StableNet(R)"},
		{name: "statistic link", handlerServerTestCase: statisticLinkHandlerTest(), wantErr: "could not fetch data for statistic link from server: could not retrieve metrics from StableNet(R)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawHandler, loggedBytes := setUpHandlerAndLogReceiver()
			clientArgs := make([]interface{}, 0, len(tt.clientArgs))
			for _, argument := range tt.clientArgs {
				clientArgs = append(clientArgs, argument.value)
			}
			rawHandler.SnClient.(*mockSnClient).On(tt.clientMethod, clientArgs...).Return(nil, errors.New("internal server error"))
			handler := tt.handler(rawHandler)
			jsonStructure := make(map[string]interface{})
			for _, argument := range tt.queryArgs {
				jsonStructure[argument.name] = argument.value
			}
			queryData, _ := json.Marshal(jsonStructure)
			actual, err := handler.Process(Query{
				Query: datasource.Query{ModelJson: string(queryData), RefId: "The cake is a lie"},
			})
			assert.Nil(t, actual, "result should be nil")
			assert.EqualError(t, err, tt.wantErr+": internal server error")
			assert.Equal(t, "no time [ERROR] "+tt.wantErr+": internal server error\n", loggedBytes.String())
		})
	}
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
		{name: "metric data no metricId", handler: metricDataHandler{}, json: "{\"measurementObid\": 1626}", wantErr: "could not extract metricIds from query"},
		{name: "statisticLinkHandler", handler: statisticLinkHandler{}, json: "{}", wantErr: "could not extract statisticLink parameter from query"},
		{name: "statisticLinkHandler no measurement id", handler: statisticLinkHandler{}, json: "{\"statisticLink\":\"hello\"}", wantErr: "the link \"hello\" does not carry a measurement id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := datasource.Query{ModelJson: tt.json, RefId: "the cake is a lie"}
			result, err := tt.handler.Process(Query{Query: query})
			assert.NoError(t, err, "on client fails, no error should be returned")
			require.NotNil(t, result, "result should not be nil")
			assert.Equal(t, "the cake is a lie", result.RefId, "the refId is wrong")
			assert.Equal(t, tt.wantErr, result.Error, "error message contained in result is wrong")
		})
	}
}

type arg struct {
	name  string
	value interface{}
}

type handlerServerTestCase struct {
	handler       func(*StableNetHandler) Handler
	queryArgs     []arg
	clientMethod  string
	clientArgs    []arg
	clientReturn  interface{}
	successResult *datasource.QueryResult
}

func datasourceTestHandlerTest() *handlerServerTestCase {
	clientReturn := []stablenet.Measurement{{}}
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return datasourceTestHandler{StableNetHandler: h} },
		queryArgs:     []arg{},
		clientMethod:  "FetchMeasurementsForDevice",
		clientArgs:    []arg{{value: -1}},
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{},
	}
}

func deviceHandlerTest() *handlerServerTestCase {
	args := []arg{{name: "deviceQuery", value: "lab"}}
	clientReturn := []stablenet.Device{{Name: "london.routerlab", Obid: 1024}, {Name: "berlin.routerlab", Obid: 5055}}
	metaJson, _ := json.Marshal(clientReturn)
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return deviceHandler{StableNetHandler: h} },
		queryArgs:     args,
		clientMethod:  "QueryDevices",
		clientArgs:    args,
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{MetaJson: string(metaJson), Series: []*datasource.TimeSeries{}},
	}
}

func measurementHandlerTest() *handlerServerTestCase {
	args := []arg{{name: "deviceObid", value: 1024}}
	clientReturn := []stablenet.Measurement{{Name: "london.routerlab Host", Obid: 4362}, {Name: "londen.routerlab Processor", Obid: 2623}}
	metaJson, _ := json.Marshal(clientReturn)
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return measurementHandler{StableNetHandler: h} },
		queryArgs:     args,
		clientMethod:  "FetchMeasurementsForDevice",
		clientArgs:    args,
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{MetaJson: string(metaJson), Series: []*datasource.TimeSeries{}},
	}
}

func metricNameHandlerTest() *handlerServerTestCase {
	args := []arg{{name: "measurementObid", value: 111}}
	clientReturn := []stablenet.Metric{{Name: "Uptime", Id: 4002}, {Name: "Processes", Id: 2003}}
	metaJson, _ := json.Marshal(clientReturn)
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return metricNameHandler{StableNetHandler: h} },
		queryArgs:     args,
		clientMethod:  "FetchMetricsForMeasurement",
		clientArgs:    args,
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{MetaJson: string(metaJson), Series: []*datasource.TimeSeries{}},
	}
}

func metricDataHandlerTest() *handlerServerTestCase {
	queryArgs := []arg{{name: "measurementObid", value: 1111}, {name: "metricIds", value: []int{123}}, {name: "includeMinStats", value: true}, {name: "includeMaxStats", value: true}}
	clientArgs := []arg{{value: 1111}, {value: []int{123}}, {value: time.Time{}}, {value: time.Time{}}}
	clientReturn, timeSeries := sampleStatisticData()
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return metricDataHandler{StableNetHandler: h} },
		queryArgs:     queryArgs,
		clientMethod:  "FetchDataForMetrics",
		clientArgs:    clientArgs,
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{Series: timeSeries[0:2]},
	}
}

func sampleStatisticData() (map[string]stablenet.MetricDataSeries, []*datasource.TimeSeries) {
	now := time.Now()
	then := time.Now().Add(5 * time.Minute)
	md1 := stablenet.MetricData{
		Interval: 5 * time.Minute,
		Time:     now,
		Min:      5,
		Max:      10,
		Avg:      7.5,
	}
	md2 := stablenet.MetricData{
		Interval: 6 * time.Minute,
		Time:     then,
		Min:      20,
		Max:      30,
		Avg:      25,
	}
	dataSeries := map[string]stablenet.MetricDataSeries{"System Uptime": {md1, md2}}
	maxSeries := &datasource.TimeSeries{
		Name:   "Max System Uptime",
		Points: []*datasource.Point{{Timestamp: now.UnixNano() / int64(time.Millisecond), Value: 10}, {Timestamp: then.UnixNano() / int64(time.Millisecond), Value: 30}},
	}
	minSeries := &datasource.TimeSeries{
		Name:   "Min System Uptime",
		Points: []*datasource.Point{{Timestamp: now.UnixNano() / int64(time.Millisecond), Value: 5}, {Timestamp: then.UnixNano() / int64(time.Millisecond), Value: 20}},
	}
	avgSeries := &datasource.TimeSeries{
		Name:   "Avg System Uptime",
		Points: []*datasource.Point{{Timestamp: now.UnixNano() / int64(time.Millisecond), Value: 7.5}, {Timestamp: then.UnixNano() / int64(time.Millisecond), Value: 25}},
	}
	return dataSeries, []*datasource.TimeSeries{minSeries, maxSeries, avgSeries}
}

func statisticLinkHandlerTest() *handlerServerTestCase {
	queryArgs := []arg{{name: "statisticLink", value: "stable.net/rest?id=1234&value0=1&value1=2"}, {name: "includeMinStats", value: true}, {name: "includeMaxStats", value: true}}
	clientArgs := []arg{{value: 1234}, {value: []int{1, 2}}, {value: time.Time{}}, {value: time.Time{}}}
	clientReturn, timeSeries := sampleStatisticData()
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return statisticLinkHandler{StableNetHandler: h} },
		queryArgs:     queryArgs,
		clientMethod:  "FetchDataForMetrics",
		clientArgs:    clientArgs,
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{Series: timeSeries[0:2]},
	}
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
	if args.Get(0) != nil {
		return args.Get(0).([]stablenet.Metric), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchDataForMetrics(measurementObid int, metrics []int, start time.Time, end time.Time) (map[string]stablenet.MetricDataSeries, error) {
	args := m.Called(measurementObid, metrics, start, end)
	if args.Get(0) != nil {
		return args.Get(0).(map[string]stablenet.MetricDataSeries), args.Error(1)
	}
	return nil, args.Error(1)
}
