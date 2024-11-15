/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"backend-plugin/mock"
	"backend-plugin/stablenet"
	"context"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type handlerTests struct {
	name       string
	urlParams  string
	filter     string
	wantStatus int
	wantErrMsg string
}

func TestDataSource_QueryData(t *testing.T) {
	stableNetUsername := "infosim"
	stableNetPassword := "stablenet"

	snServer := mock.CreateMockServer(stableNetUsername, stableNetPassword)
	handler := mock.CreateHandler(snServer)
	server := httptest.NewServer(handler)
	defer server.Close()

	byteData, _ := json.Marshal(map[string]string{
		"snusername": stableNetUsername,
		"snip":       "to be changed by context",
		"snport":     "5443",
	})

	instanceSettings := backend.DataSourceInstanceSettings{
		ID:       5,
		JSONData: byteData,
		DecryptedSecureJSONData: map[string]string{
			"snpassword": stableNetPassword,
		},
	}

	dataQueryByteData, _ := json.Marshal(map[string]interface{}{
		"StatisticLink":   "?id=1001",
		"includeMinStats": true,
		"mode":            StatisticLink,
	})

	request := backend.QueryDataRequest{
		PluginContext: backend.PluginContext{DataSourceInstanceSettings: &instanceSettings},
		Queries: []backend.DataQuery{
			{RefID: "A", JSON: dataQueryByteData},
		},
	}

	ctx := context.WithValue(context.Background(), "sn_address", server.URL)
	datasource := dataSource{validationStore: map[int64]bool{5: true}}
	got, err := datasource.QueryData(ctx, &request)
	require.NoError(t, err, "no error expected")
	require.Equal(t, 1, len(got.Responses), "number of responses wrong")
	frames := got.Responses["A"].Frames
	require.Equal(t, 1, len(frames), "number of frames wrong")
	name := frames[0].Name
	assert.Equal(t, snServer.Metrics[0].Name, name, "name of frame is wrong")
	assert.Equal(t, 5.0, frames[0].Fields[1].At(0), "value is wrong")
}

func TestHandleDeviceQuery(t *testing.T) {
	snServer := mock.CreateMockServer("infosim", "stablenet")
	handler := mock.CreateHandler(snServer)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := stablenet.NewStableNetClient(&stablenet.ConnectOptions{Username: snServer.Username, Password: snServer.Password, Address: server.URL})
	t.Run("empty device filter", func(t *testing.T) {
		request := httptest.NewRequest("GET", "http://example.org/", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleDeviceQuery(recorder, request)
		assert.Equal(t, 200, recorder.Result().StatusCode, "status is wrong")
		var got stablenet.DeviceQueryResult
		_ = json.Unmarshal(recorder.Body.Bytes(), &got)
		assert.Equal(t, snServer.Devices, got.Data, "devices wrong")
		wantQueryParam := map[string][]string{
			"$orderBy": {"name"},
			"$top":     {"100"},
		}
		assert.Equal(t, url.Values(wantQueryParam), snServer.LastQueries, "no queries wrong params")
	})
	t.Run("device filter", func(t *testing.T) {
		request := httptest.NewRequest("GET", "http://example.org?filter=bach", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleDeviceQuery(recorder, request)
		assert.Equal(t, 200, recorder.Result().StatusCode, "status is wrong")
		var got stablenet.DeviceQueryResult
		_ = json.Unmarshal(recorder.Body.Bytes(), &got)
		assert.Equal(t, snServer.Devices, got.Data, "devices wrong")
		wantQueryParam := map[string][]string{
			"$orderBy": {"name"},
			"$filter":  {"name ct 'bach'"},
			"$top":     {"100"},
		}
		assert.Equal(t, url.Values(wantQueryParam), snServer.LastQueries, "no queries wrong params")
	})
	t.Run("server error", func(t *testing.T) {
		client := stablenet.NewStableNetClient(&stablenet.ConnectOptions{Username: "", Password: "", Address: server.URL})
		request := httptest.NewRequest("GET", "http://example.org/", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleDeviceQuery(recorder, request)
		assert.Equal(t, 500, recorder.Result().StatusCode, "status is wrong")
		assert.Equal(t, "could not query devices: retrieving devices matching query \"\" failed: status code: 401, response: Authentication Error\n\n", recorder.Body.String(), "error message is wrong")
	})
}

func TestHandleMeasurementQuery(t *testing.T) {
	snServer := mock.CreateMockServer("infosim", "stablenet")
	handler := mock.CreateHandler(snServer)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := stablenet.NewStableNetClient(&stablenet.ConnectOptions{Username: snServer.Username, Password: snServer.Password, Address: server.URL})
	requestErrorTests := []handlerTests{
		{name: "no device obid", urlParams: "?obid=4500&filter=host", wantStatus: 400, wantErrMsg: "could not parse deviceObid query param: strconv.Atoi: parsing \"\": invalid syntax"},
		{name: "unparsable device obid", urlParams: "?deviceObid=a_string", wantStatus: 400, wantErrMsg: "could not parse deviceObid query param: strconv.Atoi: parsing \"a_string\": invalid syntax"},
	}
	for _, tt := range requestErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://example.org/"+tt.urlParams, strings.NewReader(""))
			request = request.WithContext(context.WithValue(request.Context(), "SnClient", client))
			recorder := httptest.NewRecorder()
			handleMeasurementQuery(recorder, request)
			assert.Equal(t, tt.wantStatus, recorder.Result().StatusCode, "status is wrong")
			assert.Equal(t, tt.wantErrMsg+"\n", recorder.Body.String(), "error message is wrong")
		})
	}

	t.Run("success without filter", func(t *testing.T) {
		request := httptest.NewRequest("GET", "http://example.org/?deviceObid=4500", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleMeasurementQuery(recorder, request)
		assert.Equal(t, 200, recorder.Result().StatusCode, "status is wrong")
		wantValues := map[string][]string{
			"$filter":  {"destDeviceId eq '4500'"},
			"$orderBy": {"name"},
			"$top":     {"100"},
		}
		assert.Equal(t, url.Values(wantValues), snServer.LastQueries, "query params are wrong")
		var got stablenet.MeasurementQueryResult
		err := json.Unmarshal(recorder.Body.Bytes(), &got)
		require.NoError(t, err, "no error expected")
		assert.Equal(t, snServer.Measurements, got.Data, "measurements differ")
	})
	t.Run("success with filter", func(t *testing.T) {
		request := httptest.NewRequest("GET", "http://example.org/?deviceObid=4500&filter=processor", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleMeasurementQuery(recorder, request)
		assert.Equal(t, 200, recorder.Result().StatusCode, "status is wrong")
		wantValues := map[string][]string{
			"$filter":  {"destDeviceId eq '4500' and name ct 'processor'"},
			"$orderBy": {"name"},
			"$top":     {"100"},
		}
		assert.Equal(t, url.Values(wantValues), snServer.LastQueries, "query params are wrong")
		var got stablenet.MeasurementQueryResult
		err := json.Unmarshal(recorder.Body.Bytes(), &got)
		require.NoError(t, err, "no error expected")
		assert.Equal(t, snServer.Measurements, got.Data, "measurements differ")
	})
	t.Run("server error", func(t *testing.T) {
		client := stablenet.NewStableNetClient(&stablenet.ConnectOptions{Username: "", Password: "", Address: server.URL})
		request := httptest.NewRequest("GET", "http://example.org/?deviceObid=1111", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleMeasurementQuery(recorder, request)
		assert.Equal(t, 500, recorder.Result().StatusCode, "status is wrong")
		assert.Equal(t, "could not query measurements: retrieving measurements for device filter \"destDeviceId eq '1111'\" failed: status code: 401, response: Authentication Error\n\n", recorder.Body.String(), "error message is wrong")
	})
}

func TestHandleMetricQuery(t *testing.T) {
	snServer := mock.CreateMockServer("infosim", "stablenet")
	handler := mock.CreateHandler(snServer)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := stablenet.NewStableNetClient(&stablenet.ConnectOptions{Username: snServer.Username, Password: snServer.Password, Address: server.URL})
	tests := []handlerTests{
		{name: "no measurement obid", urlParams: "?obid=4500&filter=host", wantStatus: 400, wantErrMsg: "could not extract measurementObid from request body: strconv.Atoi: parsing \"\": invalid syntax"},
		{name: "unparsable measurement obid", urlParams: "?measurementObid=string", wantStatus: 400, wantErrMsg: "could not extract measurementObid from request body: strconv.Atoi: parsing \"string\": invalid syntax"},
		{name: "success", urlParams: "?measurementObid=1001", wantStatus: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://example.org/"+tt.urlParams, strings.NewReader(""))
			ctx := context.WithValue(request.Context(), "SnClient", client)
			request = request.WithContext(ctx)
			recorder := httptest.NewRecorder()
			handleMetricQuery(recorder, request)
			assert.Equal(t, tt.wantStatus, recorder.Result().StatusCode, "status is wrong")
			if len(tt.wantErrMsg) == 0 {
				var got []stablenet.Metric
				err := json.Unmarshal(recorder.Body.Bytes(), &got)
				require.NoError(t, err, "no error expected")
				assert.Equal(t, snServer.Metrics, got, "metrics differ")
			} else {
				assert.Equal(t, tt.wantErrMsg+"\n", recorder.Body.String(), "error message is wrong")
			}
		})
	}
	t.Run("server error", func(t *testing.T) {
		client := stablenet.NewStableNetClient(&stablenet.ConnectOptions{Username: "", Password: "", Address: server.URL})
		request := httptest.NewRequest("GET", "http://example.org/?measurementObid=1001", strings.NewReader(""))
		ctx := context.WithValue(request.Context(), "SnClient", client)
		request = request.WithContext(ctx)
		recorder := httptest.NewRecorder()
		handleMetricQuery(recorder, request)
		assert.Equal(t, 500, recorder.Result().StatusCode, "status is wrong")
		assert.Equal(t, "could not query metrics: retrieving metrics for measurement 1001 failed: status code: 401, response: Authentication Error\n\n", recorder.Body.String(), "error message is wrong")
	})
}

func TestEncodeVersion(t *testing.T) {
	t.Run("panics", func(t *testing.T) {
		assert.PanicsWithError(t, "unable to marshal data: json: unsupported type: chan int", func() {
			encodeJson(httptest.NewRecorder(), make(chan int))
		}, "function should panic")
	})
	t.Run("no panic", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		encodeJson(recorder, 42)
		assert.Equal(t, "42\n", recorder.Body.String(), "json encoding wrong")
	})
}
