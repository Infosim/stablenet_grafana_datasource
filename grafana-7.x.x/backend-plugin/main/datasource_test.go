/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"backend-plugin/stablenet"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerTests struct {
	name       string
	urlParams  string
	filter     string
	wantStatus int
	wantErrMsg string
}

type mockVersionProvider struct {
	version   *stablenet.ServerVersion
	errString *string
}

func (m *mockVersionProvider) QueryStableNetVersion() (*stablenet.ServerVersion, *string) {
	return m.version, m.errString
}

func TestQueryTest(t *testing.T) {
	tests := []struct {
		name       string
		snVersion  *stablenet.ServerVersion
		wantStatus int
		wantBody   string
	}{
		{name: "too old", snVersion: &stablenet.ServerVersion{Version: "8.5.0"}, wantStatus: http.StatusInternalServerError, wantBody: "The StableNet® version 8.5.0 does not support Grafana®\n"},
		{name: "recent", snVersion: &stablenet.ServerVersion{Version: "9.0.0"}, wantStatus: http.StatusNoContent},
		{name: "recent with productname should fail", snVersion: &stablenet.ServerVersion{Version: "StableNet 9.0.0"}, wantStatus: http.StatusInternalServerError, wantBody: "The StableNet® version StableNet 9.0.0 does not support Grafana®\n"},
		{name: "future", snVersion: &stablenet.ServerVersion{Version: "10.1.0"}, wantStatus: http.StatusNoContent},
		{name: "internal error", snVersion: nil, wantStatus: http.StatusBadRequest, wantBody: "internal version error\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://example.org/", strings.NewReader(""))
			provider := &mockVersionProvider{}
			if tt.snVersion != nil {
				provider.version = tt.snVersion
			} else {
				errStr := "internal version error"
				provider.errString = &errStr
			}
			ctx := context.WithValue(request.Context(), "SnClient", provider)
			request = request.WithContext(ctx)
			recorder := httptest.NewRecorder()
			handleTest(recorder, request)
			assert.Equal(t, tt.wantStatus, recorder.Result().StatusCode, "status is wrong")
			assert.Equal(t, tt.wantBody, recorder.Body.String(), "response body is wrong")
		})
	}
}

type mockDeviceProvider struct {
	devices *stablenet.DeviceQueryResult
	err     error
	filter  string
}

func (m *mockDeviceProvider) QueryDevices(filter string) (*stablenet.DeviceQueryResult, error) {
	m.filter = filter
	return m.devices, m.err
}

func TestHandleDeviceQuery(t *testing.T) {
	tests := []handlerTests{
		{name: "empty device filter", urlParams: "", wantStatus: 200},
		{name: "device filter", urlParams: "?filter=bach", wantStatus: 200},
		{name: "internal device query error", urlParams: "?filter=bach", wantStatus: 500, wantErrMsg: "could not query devices: internal device error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://example.org/"+tt.urlParams, strings.NewReader(""))
			provider := &mockDeviceProvider{}
			devices := &stablenet.DeviceQueryResult{Devices: []stablenet.Device{{Name: "Berlin", Obid: 2323}, {Name: "Dallas", Obid: 923}}}
			if len(tt.wantErrMsg) > 0 {
				provider.err = errors.New("internal device error")
			} else {
				provider.devices = devices
			}
			ctx := context.WithValue(request.Context(), "SnClient", provider)
			request = request.WithContext(ctx)
			recorder := httptest.NewRecorder()
			handleDeviceQuery(recorder, request)
			assert.Equal(t, tt.wantStatus, recorder.Result().StatusCode, "status is wrong")
			if len(tt.wantErrMsg) == 0 {
				assert.Equal(t, request.URL.Query().Get("filter"), provider.filter, "provider not called with correct filter")
				var got stablenet.DeviceQueryResult
				err := json.Unmarshal(recorder.Body.Bytes(), &got)
				require.NoError(t, err, "no error expected")
				assert.Equal(t, *devices, got, "measurements differ")
			} else {
				assert.Equal(t, tt.wantErrMsg+"\n", recorder.Body.String(), "error message is wrong")
			}
		})
	}
}

type mockMeasurementProvider struct {
	measurements *stablenet.MeasurementQueryResult
	filter       string
	err          error
	obid         int
}

func (m *mockMeasurementProvider) FetchMeasurementsForDevice(obid int, filter string) (*stablenet.MeasurementQueryResult, error) {
	m.obid = obid
	m.filter = filter
	return m.measurements, m.err
}

func TestHandleMeasurementQuery(t *testing.T) {
	tests := []handlerTests{
		{name: "no device obid", urlParams: "?obid=4500&filter=host", wantStatus: 400, wantErrMsg: "could not parse deviceObid query param: strconv.Atoi: parsing \"\": invalid syntax"},
		{name: "unparsable device obid", urlParams: "?deviceObid=a_string", wantStatus: 400, wantErrMsg: "could not parse deviceObid query param: strconv.Atoi: parsing \"a_string\": invalid syntax"},
		{name: "measurement provider error", urlParams: "?deviceObid=4500", wantStatus: 500, wantErrMsg: "could not query measurements: internal measurement error"},
		{name: "success", urlParams: "?deviceObid=4500", wantStatus: 200},
		{name: "success", urlParams: "?deviceObid=-1&filter=processor", wantStatus: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://example.org/"+tt.urlParams, strings.NewReader(""))
			provider := &mockMeasurementProvider{}
			measurements := &stablenet.MeasurementQueryResult{Measurements: []stablenet.Measurement{{Name: "Host", Obid: 2323}, {Name: "Uptime", Obid: 923}}}
			if len(tt.wantErrMsg) > 0 {
				provider.err = errors.New("internal measurement error")
			} else {
				provider.measurements = measurements
			}
			ctx := context.WithValue(request.Context(), "SnClient", provider)
			request = request.WithContext(ctx)
			recorder := httptest.NewRecorder()
			handleMeasurementQuery(recorder, request)
			assert.Equal(t, tt.wantStatus, recorder.Result().StatusCode, "status is wrong")
			if len(tt.wantErrMsg) == 0 {
				assert.Equal(t, request.URL.Query().Get("filter"), provider.filter, "filter should be set")
				assert.Equal(t, request.URL.Query().Get("deviceObid"), fmt.Sprintf("%d", provider.obid), "provider not called with correct obid")
				var got stablenet.MeasurementQueryResult
				err := json.Unmarshal(recorder.Body.Bytes(), &got)
				require.NoError(t, err, "no error expected")
				assert.Equal(t, *measurements, got, "measurements differ")
			} else {
				assert.Equal(t, tt.wantErrMsg+"\n", recorder.Body.String(), "error message is wrong")
			}
		})
	}
}

type mockMetricProvider struct {
	metrics []stablenet.Metric
	err     error
	obid    int
}

func (m *mockMetricProvider) FetchMetricsForMeasurement(obid int) ([]stablenet.Metric, error) {
	m.obid = obid
	return m.metrics, m.err
}

func TestHandleMetricQuery(t *testing.T) {
	tests := []handlerTests{
		{name: "no measurement obid", urlParams: "?obid=4500&filter=host", wantStatus: 400, wantErrMsg: "could not extract measurementObid from request body: strconv.Atoi: parsing \"\": invalid syntax"},
		{name: "unparsable measurement obid", urlParams: "?measurementObid=string", wantStatus: 400, wantErrMsg: "could not extract measurementObid from request body: strconv.Atoi: parsing \"string\": invalid syntax"},
		{name: "metric provider error", urlParams: "?measurementObid=4500", wantStatus: 500, wantErrMsg: "could not query metrics: internal metric error"},
		{name: "success", urlParams: "?measurementObid=4500", wantStatus: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://example.org/"+tt.urlParams, strings.NewReader(""))
			provider := &mockMetricProvider{}
			metrics := []stablenet.Metric{{Name: "Host", Key: "SNMP_1"}, {Name: "Uptime", Key: "SNMP_2"}}
			if len(tt.wantErrMsg) > 0 {
				provider.err = errors.New("internal metric error")
			} else {
				provider.metrics = metrics
			}
			ctx := context.WithValue(request.Context(), "SnClient", provider)
			request = request.WithContext(ctx)
			recorder := httptest.NewRecorder()
			handleMetricQuery(recorder, request)
			assert.Equal(t, tt.wantStatus, recorder.Result().StatusCode, "status is wrong")
			if len(tt.wantErrMsg) == 0 {
				assert.Equal(t, 4500, provider.obid, "provider not called with correct obid")
				var got []stablenet.Metric
				err := json.Unmarshal(recorder.Body.Bytes(), &got)
				require.NoError(t, err, "no error expected")
				assert.Equal(t, metrics, got, "metrics differ")
			} else {
				assert.Equal(t, tt.wantErrMsg+"\n", recorder.Body.String(), "error message is wrong")
			}
		})
	}
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
