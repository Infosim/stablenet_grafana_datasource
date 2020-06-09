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
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"net/http"
	"regexp"
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

	return datasource.ServeOpts{
		CheckHealthHandler:  ds,
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

func (ds *testDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	settings, err := ds.getSettings(req.PluginContext)
	if err != nil {
		return nil, err
	}

	// Handle request
	resp, err := settings.httpClient.Get("http://")
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return nil, nil
}

func (ds *testDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	var resp *backend.QueryDataResponse
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	err = ds.im.Do(req.PluginContext, func(settings *testDataSourceInstanceSettings) error {
		// Handle request
		resp, err := settings.httpClient.Get("http://")
		if err != nil {
			return err
		}
		resp.Body.Close()
		return nil
	})

	return resp, err
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