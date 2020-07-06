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
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
)

type dataSource struct {
}

func newDataSource() datasource.ServeOpts {
	ds := &dataSource{}

	addClientThen := func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					backend.Logger.Error(fmt.Sprintf("An error occured in resource query %s: %v\n%s", req.URL.RawPath, err, debug.Stack()))
				}
			}()
			pluginContext := httpadapter.PluginConfigFromContext(req.Context())
			options := stableNetOptions(pluginContext.DataSourceInstanceSettings)
			client := stablenet.NewClient(options)
			ctx := context.WithValue(req.Context(), "SnClient", client)
			next.ServeHTTP(rw, req.WithContext(ctx))
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/devices", addClientThen(handleDeviceQuery))
	mux.HandleFunc("/measurements", addClientThen(handleMeasurementQuery))
	mux.HandleFunc("/metrics", addClientThen(handleMetricQuery))

	return datasource.ServeOpts{
		CallResourceHandler: httpadapter.New(mux),
		CheckHealthHandler:  ds,
		QueryDataHandler:    ds,
	}
}

func (ds *dataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			backend.Logger.Error(fmt.Sprintf("An error occured: %v\n%s", err, debug.Stack()))
		}
	}()
	options := stableNetOptions(req.PluginContext.DataSourceInstanceSettings)

	// We need this for testing purposes. Go's httptest package only allows to mock http, not https, and
	// it is not meant to separate ip and port. Thus, for testing purposes, we inject the test url here.
	if ctx.Value("sn_address") != nil {
		options.Address = ctx.Value("sn_address").(string)
	}

	queries := make([]MetricQuery, 0, len(req.Queries))
	for index, singleRequest := range req.Queries {
		query := NewQuery(singleRequest)
		err := json.Unmarshal(singleRequest.JSON, &query)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize query %d: %v", index, err)
		}
		queries = append(queries, query)
	}
	client := stablenet.NewClient(options)
	queries, err := ExpandStatisticLinks(queries, client.FetchMetricsForMeasurement)
	if err != nil {
		return nil, err
	}
	allFrames := make([]*data.Frame, 0, 0)
	for _, query := range queries {
		frames, err := query.FetchData(client.FetchDataForMetrics)
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

func (ds *dataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	defer func() {
		if err := recover(); err != nil {
			backend.Logger.Error(fmt.Sprintf("An error occured: %v\n%s", err, debug.Stack()))
		}
	}()
	options := stableNetOptions(req.PluginContext.DataSourceInstanceSettings)

	// We need this for testing purposes. Go's httptest package only allows to mock http, not https, and
	// it is not meant to separate ip and port. Thus, for testing purposes, we inject the test url here.
	if ctx.Value("sn_address") != nil {
		options.Address = ctx.Value("sn_address").(string)
	}

	client := stablenet.NewClient(options)
	version, errStr := client.QueryStableNetVersion()
	if errStr != nil {
		return &backend.CheckHealthResult{Status: backend.HealthStatusError, Message: *errStr}, nil
	}
	versionRegex := regexp.MustCompile("^(?:9|[1-9]\\d)\\.")
	if !versionRegex.MatchString(version.Version) {
		return &backend.CheckHealthResult{Status: backend.HealthStatusError, Message: fmt.Sprintf("The StableNet® version %s does not support Grafana®", version.Version)}, nil
	}
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}

func handleDeviceQuery(rw http.ResponseWriter, req *http.Request) {
	snClient := req.Context().Value("SnClient").(*stablenet.Client)
	filter := req.URL.Query().Get("filter")
	devices, err := snClient.QueryDevices(filter)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query devices: %v", err), http.StatusInternalServerError)
		return
	}
	encodeJson(rw, devices)
}

func handleMeasurementQuery(rw http.ResponseWriter, req *http.Request) {
	snClient := req.Context().Value("SnClient").(*stablenet.Client)
	filter := req.URL.Query().Get("filter")
	deviceObid, err := strconv.Atoi(req.URL.Query().Get("deviceObid"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not parse deviceObid query param: %v", err), http.StatusBadRequest)
		return
	}
	measurements, err := snClient.FetchMeasurementsForDevice(deviceObid, filter)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query measurements: %v", err), http.StatusInternalServerError)
		return
	}
	encodeJson(rw, measurements)
}

func handleMetricQuery(rw http.ResponseWriter, req *http.Request) {
	snClient := req.Context().Value("SnClient").(*stablenet.Client)
	measurementObid, err := strconv.Atoi(req.URL.Query().Get("measurementObid"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract measurementObid from request body: %v", err), http.StatusBadRequest)
		return
	}
	metrics, err := snClient.FetchMetricsForMeasurement(measurementObid)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query metrics: %v", err), http.StatusInternalServerError)
		return
	}
	encodeJson(rw, metrics)
}

// Encoding a json only results in an error if the data to be serialized does
// contain unserializable types, e.g. functions, channels, etc.
// Since we have absolute control over our types, we panic in case the json cannot be created.
func encodeJson(rw http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		panic(fmt.Errorf("unable to marshal data: %v", err))
	}
}
