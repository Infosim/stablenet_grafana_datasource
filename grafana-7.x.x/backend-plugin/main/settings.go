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
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"strconv"
)

func stableNetOptions(settings *backend.DataSourceInstanceSettings) (*stablenet.ConnectOptions, error) {
	if settings == nil {
		return nil, fmt.Errorf("settings is nil")
	}
	options := make(map[string]string)
	err := json.Unmarshal(settings.JSONData, &options)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSONData of the datasource: %v", err)
	}
	if _, ok := options["snip"]; !ok {
		return nil, fmt.Errorf("the snip is missing in the JSONData of the datasource")
	}
	if _, ok := options["snport"]; !ok {
		return nil, fmt.Errorf("the snport is missing in the jsonData of the datasource")
	}
	if _, ok := options["snusername"]; !ok {
		return nil, fmt.Errorf("the snusername is missing in the JSONData of the datasource")
	}
	if _, ok := settings.DecryptedSecureJSONData["snpassword"]; !ok {
		return nil, fmt.Errorf("the snpassword is missing in the encryptedJSONData of the datasource")
	}
	port, portErr := strconv.Atoi(options["snport"])
	if portErr != nil {
		return nil, fmt.Errorf("could not parse snport into number: %v", portErr)
	}
	return &stablenet.ConnectOptions{
		Host:     options["snip"],
		Port:     port,
		Username: options["snusername"],
		Password: settings.DecryptedSecureJSONData["snpassword"],
	}, nil
}
