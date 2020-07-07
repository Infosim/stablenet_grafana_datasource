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
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"regexp"
	"runtime/debug"
)

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
	valid, msg := ds.checkAndUpdateHealth(options, req.PluginContext.DataSourceInstanceSettings.ID)
	status := backend.HealthStatusError
	if valid {
		status = backend.HealthStatusOk
	}
	return &backend.CheckHealthResult{Status: status, Message: msg}, nil
}

func (ds *dataSource) checkAndUpdateHealth(options *stablenet.ConnectOptions, datasourceId int64) (bool, string) {
	client := stablenet.NewClient(options)
	info, errStr := client.QueryStableNetInfo()
	if errStr != nil {
		return false, *errStr
	}
	versionRegex := regexp.MustCompile("^(?:9|[1-9]\\d)\\.")
	if !versionRegex.MatchString(info.ServerVersion.Version) {
		ds.validationStore[datasourceId] = false
		return false, fmt.Sprintf("The StableNet® version %s does not support Grafana®.", info.ServerVersion.Version)
	}
	if !info.License.Modules.IsRestReportingLicensed() {
		ds.validationStore[datasourceId] = false
		return false, fmt.Sprintf("The StableNet® server does not have the required license \"REST_REPORTING\".")
	}
	ds.validationStore[datasourceId] = true
	return true, "Connection to StableNet® successful"
}
