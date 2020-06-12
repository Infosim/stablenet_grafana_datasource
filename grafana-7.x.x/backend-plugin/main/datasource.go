/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
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
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"regexp"
	"time"
)

type testDataSourceInstanceSettings struct {
	httpClient *http.Client
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &testDataSourceInstanceSettings{
		httpClient: &http.Client{},
	}, nil
}

func (s *testDataSourceInstanceSettings) Dispose() {
	// Cleanup
}

type testDataSource struct {
	im     instancemgmt.InstanceManager
	logger hclog.Logger
}

func newDataSource(logger hclog.Logger) datasource.ServeOpts {
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &testDataSource{
		im:     im,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/test", ds.handleTest)
	mux.HandleFunc("/devices", handleDeviceQuery)
	mux.HandleFunc("/measurements", handleMeasurementQuery)
	mux.HandleFunc("/metrics", handleMetricQuery)

	return datasource.ServeOpts{
		CallResourceHandler: httpadapter.New(mux),
		QueryDataHandler:    ds,
	}
}

func (ds *testDataSource) getSettings(pluginContext backend.PluginContext) (*testDataSourceInstanceSettings, error) {
	iface, err := ds.im.Get(pluginContext)
	if err != nil {
		return nil, err
	}

	return iface.(*testDataSourceInstanceSettings), nil
}

type MetricQuery struct {
	IncludeAvgStats bool
	IncludeMaxStats bool
	IncludeMinStats bool
	RequestData     []MetricRequest
}

type MetricRequest struct {
	MeasurementObid int
	Metrics         []struct {
		Key  string
		Name string
	}
}

func (ds *testDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	_, err := stableNetOptions(req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		return nil, fmt.Errorf("could not extract data source settings: %v", err)
	}

	queries := make([]MetricQuery, 0, len(req.Queries))
	for index, singleRequest := range req.Queries {
		query := MetricQuery{}
		err := json.Unmarshal(singleRequest.JSON, &query)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize query %d: %v", index, err)
		}
		queries = append(queries, query)
	}
	frames := make([]*data.Frame, 0, len(req.Queries))
	for index, _ := range queries {
		frame := data.NewFrame(fmt.Sprintf("Frame %d", index), data.NewField("timeValues", nil, []time.Time{}), data.NewField("data", nil, []float64{}))
		frame.AppendRow(time.Now().Add(-1*time.Hour), 5.3)
		frame.AppendRow(time.Now().Add(-30*time.Minute), 6.7)
		frame.AppendRow(time.Now(), 7.8)
		frames = append(frames, frame)
	}
	backend.Logger.Warn(fmt.Sprintf("%v", frames))
	response := backend.NewQueryDataResponse()
	response.Responses = backend.Responses{"whatever": backend.DataResponse{Frames: frames}}
	return response, nil
}

func (ds *testDataSource) handleTest(rw http.ResponseWriter, req *http.Request) {
	pluginContext := httpadapter.PluginConfigFromContext(req.Context())
	options, err := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract data source settings: %v", err), http.StatusInternalServerError)
		return
	}
	snClient := stablenet.NewClient(options)
	version, errStr := snClient.QueryStableNetVersion()
	if errStr != nil {
		http.Error(rw, *errStr, http.StatusBadRequest)
		return
	}
	versionRegex := regexp.MustCompile("^(?:9|[1-9]\\d)\\.")
	if !versionRegex.MatchString(version.Version) {
		http.Error(rw, fmt.Sprintf("The StableNet® version %s does not support Grafana®", version.Version), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func handleDeviceQuery(rw http.ResponseWriter, req *http.Request) {
	pluginContext := httpadapter.PluginConfigFromContext(req.Context())
	options, err := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract data source settings: %v", err), http.StatusInternalServerError)
		return
	}
	snClient := stablenet.NewClient(options)
	filterWrapper := struct {
		Filter string
	}{}
	err = json.NewDecoder(req.Body).Decode(&filterWrapper)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract filter from request body: %v", err), http.StatusUnprocessableEntity)
		return
	}
	devices, err := snClient.QueryDevices(filterWrapper.Filter)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query devices: %v", err), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(devices)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not serialize devices: %v", err), http.StatusInternalServerError)
	}
}

func handleMeasurementQuery(rw http.ResponseWriter, req *http.Request) {
	pluginContext := httpadapter.PluginConfigFromContext(req.Context())
	options, err := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract data source settings: %v", err), http.StatusInternalServerError)
		return
	}
	snClient := stablenet.NewClient(options)
	filterWrapper := struct {
		Filter     string
		DeviceObid int
	}{}
	err = json.NewDecoder(req.Body).Decode(&filterWrapper)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract filter and deviceObid from request body: %v", err), http.StatusUnprocessableEntity)
		return
	}
	measurements, err := snClient.FetchMeasurementsForDevice(&filterWrapper.DeviceObid, filterWrapper.Filter)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query measurements: %v", err), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(measurements)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not serialize measurements: %v", err), http.StatusInternalServerError)
	}
}

func handleMetricQuery(rw http.ResponseWriter, req *http.Request) {
	pluginContext := httpadapter.PluginConfigFromContext(req.Context())
	options, err := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract source settings: %v", err), http.StatusInternalServerError)
		return
	}
	snClient := stablenet.NewClient(options)
	filterWrapper := struct {
		MeasurementObid int
	}{}
	err = json.NewDecoder(req.Body).Decode(&filterWrapper)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract measurementObid from request body: %v"), http.StatusInternalServerError)
		return
	}
	metrics, err := snClient.FetchMetricsForMeasurement(filterWrapper.MeasurementObid, "")
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query metrics: %v", err), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(metrics)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not serialize measurements: %v", err), http.StatusInternalServerError)
	}
}
