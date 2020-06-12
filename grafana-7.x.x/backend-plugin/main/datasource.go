/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	query2 "backend-plugin/query"
	"backend-plugin/stablenet"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
)

type testDataSource struct {
}

func newDataSource() datasource.ServeOpts {
	ds := &testDataSource{}

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

func (ds *testDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			backend.Logger.Error(fmt.Sprintf("An error occured: %v\n%s", err, debug.Stack()))
		}
	}()
	options := stableNetOptions(req.PluginContext.DataSourceInstanceSettings)

	queries := make([]query2.MetricQuery, 0, len(req.Queries))
	for index, singleRequest := range req.Queries {
		query := query2.NewQuery(singleRequest)
		err := json.Unmarshal(singleRequest.JSON, &query)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize query %d: %v", index, err)
		}
		queries = append(queries, query)
	}
	client := stablenet.NewClient(options)
	handler := query2.StableNetHandler{SnClient: client}
	queries, err := handler.ExpandStatisticLinks(queries)
	if err != nil {
		return nil, err
	}
	allFrames := make([]*data.Frame, 0, 0)
	for _, query := range queries {
		frames, err := handler.FetchMetrics(query)
		if err != nil {
			return nil, fmt.Errorf("could not fetch data for query %v: %v", query, err)
		}
		for _, frame := range frames {
			allFrames = append(allFrames, frame)
		}
	}
	response := backend.NewQueryDataResponse()
	response.Responses = backend.Responses{"queryResponse": backend.DataResponse{Frames: allFrames}}
	return response, nil
}

func (ds *testDataSource) handleTest(rw http.ResponseWriter, req *http.Request) {
	pluginContext := httpadapter.PluginConfigFromContext(req.Context())
	options := stableNetOptions(pluginContext.DataSourceInstanceSettings)
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
	options := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	snClient := stablenet.NewClient(options)
	filter := req.URL.Query().Get("filter")
	devices, err := snClient.QueryDevices(filter)
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
	options := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	snClient := stablenet.NewClient(options)
	deviceObid, err := strconv.Atoi(req.URL.Query().Get("deviceObid"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not parse deviceObid from request body: %v", err), http.StatusBadRequest)
		return
	}
	measurements, err := snClient.FetchMeasurementsForDevice(deviceObid)
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
	options := stableNetOptions(pluginContext.DataSourceInstanceSettings)
	snClient := stablenet.NewClient(options)
	measurementObid, err := strconv.Atoi(req.URL.Query().Get("measurementObid"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract measurementObid from request body: %v"), http.StatusInternalServerError)
		return
	}
	metrics, err := snClient.FetchMetricsForMeasurement(measurementObid)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query metrics: %v", err), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(metrics)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not serialize measurements: %v", err), http.StatusInternalServerError)
	}
}
