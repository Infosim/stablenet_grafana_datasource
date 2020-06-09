/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/hashicorp/go-hclog"
	"os"
)

var pluginLogger = hclog.New(&hclog.LoggerOptions{
	Name:  "stablenet-datasource-logger",
	Level: hclog.Info,
})

func main() {
	err := datasource.Serve(newDataSource(pluginLogger))
	if err != nil {
		pluginLogger.Error("could not initialize datasource: %v", err)
		os.Exit(1)
	}
	//plugin.Serve(&plugin.ServeConfig{
	//	Logger: pluginLogger,
	//	HandshakeConfig: plugin.HandshakeConfig{
	//		ProtocolVersion:  1,
	//		MagicCookieKey:   "grafana_plugin_type",
	//		MagicCookieValue: "datasource",
	//	},
	//	Plugins: map[string]plugin.Plugin{
	//		"stablenet-datasource": &datasource.DatasourcePluginImpl{Plugin: NewBackendPlugin(pluginLogger)},
	//	},
	//
	//	// A non-nil value here enables gRPC serving for this plugin...
	//	GRPCServer: plugin.DefaultGRPCServer,
	//})
}
