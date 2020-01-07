/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
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
	modelJson := "{\"snip\":\"127.0.0.1\", \"snport\": \"443\", \"snusername\":\"infosim\"}"
	decryptedData := map[string]string{"snpassword": "stablenet"}
	request := Request{
		DatasourceRequest: &datasource.DatasourceRequest{
			Datasource: &datasource.DatasourceInfo{
				JsonData:                modelJson,
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
	modelJson := "{\"snip\":\"127.0.0.1\", \"snport\": \"443\", \"snusername\":\"infosim\"}"
	decryptedData := map[string]string{"snpassword": "stablenet"}
	request := &Request{
		DatasourceRequest: &datasource.DatasourceRequest{
			Datasource: &datasource.DatasourceInfo{
				JsonData:                modelJson,
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
		/**{name: "statistic link", handlerServerTestCase: statisticLinkHandlerTest()},**/
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
		//statistic link handler is currently not tested here because it needs two calls to StableNet which cannot be handled by this test framework.
		//However, it is tested partially at the end of this file.
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
		{name: "metrics for measurement", handler: metricNameHandler{}, json: "{}", wantErr: "could not extract measurementObid from query"},
		{name: "metric data", handler: metricDataHandler{}, json: "{}", wantErr: "could not extract measurement requests from query: dataRequest not present in the modelJson"},
		{name: "statisticLinkHandler", handler: statisticLinkHandler{}, json: "{}", wantErr: "could not extract statisticLink parameter from query"},
		{name: "statisticLinkHandler no measurement id", handler: statisticLinkHandler{}, json: "{\"statisticLink\":\"hello\"}", wantErr: "the link \"hello\" does not carry a measurement id or value ids"},
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
	clientReturn := &stablenet.MeasurementQueryResult{HasMore: false, Measurements: []stablenet.Measurement{}}
	var nilInt *int
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return datasourceTestHandler{StableNetHandler: h} },
		queryArgs:     []arg{},
		clientMethod:  "FetchMeasurementsForDevice",
		clientArgs:    []arg{{value: nilInt}, {value: ""}},
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{},
	}
}

func deviceHandlerTest() *handlerServerTestCase {
	args := []arg{{name: "filter", value: "lab"}}
	clientReturn := &stablenet.DeviceQueryResult{Devices: []stablenet.Device{{Name: "london.routerlab", Obid: 1024}, {Name: "berlin.routerlab", Obid: 5055}}}
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
	var deviceObid = 1024
	clientReturn := &stablenet.MeasurementQueryResult{
		Measurements: []stablenet.Measurement{{Name: "london.routerlab Host", Obid: 4362}, {Name: "londen.routerlab Processor", Obid: 2623}},
		HasMore:      true,
	}
	metaJson, _ := json.Marshal(clientReturn)
	return &handlerServerTestCase{
		handler:       func(h *StableNetHandler) Handler { return measurementHandler{StableNetHandler: h} },
		queryArgs:     []arg{{name: "deviceObid", value: deviceObid}, {name: "filter", value: "o"}},
		clientMethod:  "FetchMeasurementsForDevice",
		clientArgs:    []arg{{value: &deviceObid}, {value: "o"}},
		clientReturn:  clientReturn,
		successResult: &datasource.QueryResult{MetaJson: string(metaJson), Series: []*datasource.TimeSeries{}},
	}
}

func metricNameHandlerTest() *handlerServerTestCase {
	args := []arg{{name: "measurementObid", value: 111}, {name: "filter", value: "Host"}}
	clientReturn := []stablenet.Metric{{Name: "Uptime", Key: "4002"}, {Name: "Processes", Key: "2003"}}
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
	metricReq := []stablenet.Metric{{Name: "System Uptime", Key: "123"}}
	requestData := []measurementDataRequest{
		{MeasurementObid: 1111, Metrics: metricsRequest(metricReq)},
	}
	queryArgs := []arg{{name: "requestData", value: requestData}, {name: "includeMinStats", value: true}, {name: "includeMaxStats", value: true}}
	clientArgs := []arg{{value: 1111}, {value: []string{"123"}}, {value: time.Time{}}, {value: time.Time{}}}
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
	dataSeries := map[string]stablenet.MetricDataSeries{"123": {md1, md2}}
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

func TestDatasourceTestHandler_Process_Error(t *testing.T) {
	rawHandler, _ := setUpHandlerAndLogReceiver()
	var nilInt *int
	rawHandler.SnClient.(*mockSnClient).On("FetchMeasurementsForDevice", nilInt, "").Return(nil, errors.New("login not possible"))
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

func (m *mockSnClient) QueryDevices(query string) (*stablenet.DeviceQueryResult, error) {
	args := m.Called(query)
	if args.Get(0) != nil {
		result := args.Get(0).(*stablenet.DeviceQueryResult)
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchMeasurementName(id int) (*string, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*string), nil
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchMeasurementsForDevice(deviceObid *int, nameFilter string) (*stablenet.MeasurementQueryResult, error) {
	args := m.Called(deviceObid, nameFilter)
	if args.Get(0) != nil {
		return args.Get(0).(*stablenet.MeasurementQueryResult), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchMetricsForMeasurement(measurementObid int, filter string) ([]stablenet.Metric, error) {
	args := m.Called(measurementObid, filter)
	if args.Get(0) != nil {
		return args.Get(0).([]stablenet.Metric), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSnClient) FetchDataForMetrics(measurementObid int, metricKeys []string, start time.Time, end time.Time) (map[string]stablenet.MetricDataSeries, error) {
	args := m.Called(measurementObid, metricKeys, start, end)
	if args.Get(0) != nil {
		return args.Get(0).(map[string]stablenet.MetricDataSeries), args.Error(1)
	}
	return nil, args.Error(1)
}
