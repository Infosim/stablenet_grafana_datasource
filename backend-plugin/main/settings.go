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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Parses the JSON object contained in settings and extracts the StableNet options.
// It panics if the the StableNet options are not complete or not parsable, or nil is passed as argument.
func stableNetOptions(settings *backend.DataSourceInstanceSettings) *stablenet.ConnectOptions {
	if settings == nil {
		panic("datasource settings are nil, are you in a datasource environment?")
	}

	options := make(map[string]string)

	// error checking of unmarshal is done by checking the specific fields
	_ = json.Unmarshal(settings.JSONData, &options)

	stableNetIp, stableNetIpOk := options["snip"]
	stableNetPort, stableNetPortOk := options["snport"]

	if !stableNetIpOk {
		panic("field \"snip\" is missing in the JSONData of the datasource")
	}
	if !stableNetPortOk {
		panic("the field \"snport\" is missing in the JSONData of the datasource")
	}
	if _, ok := options["snusername"]; !ok {
		panic("the field \"snusername\" is missing in the JSONData of the datasource")
	}
	if _, ok := settings.DecryptedSecureJSONData["snpassword"]; !ok {
		panic("the field \"snpassword\" is missing in the encryptedJSONData of the datasource")
	}

	port, portErr := strconv.Atoi(stableNetPort)
	if portErr != nil {
		panic(fmt.Sprintf("the field \"snport\" could not be parsed into a number: %v", portErr))
	}

	return &stablenet.ConnectOptions{
		Address:  fmt.Sprintf("https://%s:%d", stableNetIp, port),
		Username: options["snusername"],
		Password: settings.DecryptedSecureJSONData["snpassword"],
	}
}
