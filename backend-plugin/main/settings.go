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
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func loadStableNetSettings(settings *backend.DataSourceInstanceSettings) (*stablenet.ConnectOptions, error) {
	if settings == nil {
		return nil, fmt.Errorf("datasource settings are nil, are you in a datasource environment?")
	}

	if _, ok := settings.DecryptedSecureJSONData["password"]; !ok {
		return nil, fmt.Errorf("no password was provided")
	}

	return &stablenet.ConnectOptions{
		Address:  settings.URL,
		Username: settings.User,
		Password: settings.DecryptedSecureJSONData["password"],
	}, nil
}
