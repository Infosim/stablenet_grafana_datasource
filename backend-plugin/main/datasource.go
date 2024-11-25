/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
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
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

type dataSource struct {
	validationStore map[int64]bool
}

func newDataSource() datasource.ServeOpts {
	ds := &dataSource{validationStore: make(map[int64]bool)}

	addClientThen := func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					backend.Logger.Error(fmt.Sprintf("An error occured in resource query %s: %v\n%s", req.URL.RawPath, err, debug.Stack()))
				}
			}()

			pluginContext := backend.PluginConfigFromContext(req.Context())
			options, err := loadStableNetSettings(pluginContext.DataSourceInstanceSettings)

			if err != nil {
				http.Error(rw, "The datasource configuration is not valid, please check make sure that the health test is successful.", http.StatusInternalServerError)
			}

			valid, present := ds.validationStore[pluginContext.DataSourceInstanceSettings.ID]
			if !present {
				valid, _ = ds.checkAndUpdateHealth(options, pluginContext.DataSourceInstanceSettings.ID)
			}
			if !valid {
				http.Error(rw, "The datasource is not valid, please check the data source configuration and make sure that the test is successful.", http.StatusInternalServerError)
				return
			}

			client := stablenet.NewStableNetClient(options)
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

	options, err := loadStableNetSettings(req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		return nil, err
	}

	// We need this for testing purposes. Go's httptest package only allows to mock http, not https, and
	// it is not meant to separate ip and port. Thus, for testing purposes, we inject the test url here.
	if ctx.Value("sn_address") != nil {
		options.Address = ctx.Value("sn_address").(string)
	}

	valid, present := ds.validationStore[req.PluginContext.DataSourceInstanceSettings.ID]
	if !present {
		valid, _ = ds.checkAndUpdateHealth(options, req.PluginContext.DataSourceInstanceSettings.ID)
	}
	if !valid {
		responses := backend.Responses{"queryResponse": backend.DataResponse{Error: errors.New("the datasource is not valid, please check the data source configuration and make sure that the test is successful")}}
		return &backend.QueryDataResponse{Responses: responses}, nil
	}

	queries := make([]MetricQuery, 0, len(req.Queries))
	for index, singleRequest := range req.Queries {
		target := &Target{}
		err := json.Unmarshal(singleRequest.JSON, target)
		if err != nil {
			return nil, fmt.Errorf("could not deserialize query %d: %v", index, err)
		}
		query := target.toQuery(singleRequest.TimeRange, singleRequest.RefID)
		if (len(query.Metrics) == 0) && query.StatisticLink == nil {
			continue
		}
		queries = append(queries, query)
	}

	client := stablenet.NewStableNetClient(options)
	queries, err = ExpandStatisticLinks(queries, client.FetchMetricsForMeasurement)
	if err != nil {
		return nil, err
	}
	response := backend.NewQueryDataResponse()
	response.Responses = make(map[string]backend.DataResponse)
	for _, query := range queries {
		frames, err := query.FetchData(client.FetchDataForMetrics)
		if err != nil {
			return nil, fmt.Errorf("could not fetch data for query %v: %v", query, err)
		}
		response.Responses[query.RefId] = backend.DataResponse{Frames: frames}
	}
	return response, nil
}

func handleDeviceQuery(rw http.ResponseWriter, req *http.Request) {
	filter := req.URL.Query().Get("filter")

	snClient := req.Context().Value("SnClient").(*stablenet.StableNetClient)

	devices, err := snClient.QueryDevices(filter)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query devices: %v", err), http.StatusInternalServerError)
		return
	}

	encodeJson(rw, devices)
}

func handleMeasurementQuery(rw http.ResponseWriter, req *http.Request) {
	deviceObid, err := strconv.Atoi(req.URL.Query().Get("deviceObid"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not parse deviceObid query param: %v", err), http.StatusBadRequest)
		return
	}

	filter := req.URL.Query().Get("filter")

	snClient := req.Context().Value("SnClient").(*stablenet.StableNetClient)

	measurements, err := snClient.FetchMeasurementsForDevice(deviceObid, filter)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query measurements: %v", err), http.StatusInternalServerError)
		return
	}
	encodeJson(rw, measurements)
}

func handleMetricQuery(rw http.ResponseWriter, req *http.Request) {
	measurementObid, err := strconv.Atoi(req.URL.Query().Get("measurementObid"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not extract measurementObid from request body: %v", err), http.StatusBadRequest)
		return
	}

	snClient := req.Context().Value("SnClient").(*stablenet.StableNetClient)

	metrics, err := snClient.FetchMetricsForMeasurement(measurementObid)
	if err != nil {
		http.Error(rw, fmt.Sprintf("could not query metrics: %v", err), http.StatusInternalServerError)
		return
	}

	encodeJson(rw, metrics)
}

// Encoding a json only results in an error if the data to be serialized does contain unserializable types, e.g. functions, channels, etc.
// Since we have absolute control over our types, we panic in case the json cannot be created.
func encodeJson(rw http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		panic(fmt.Errorf("unable to marshal data: %v", err))
	}
}
