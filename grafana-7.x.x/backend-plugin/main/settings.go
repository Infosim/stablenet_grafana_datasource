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

// Parses the JSON object contained in settings and extracts the StableNet options.
// It panics if the the StableNet options are not complete or not parsable, or nil is passed as argument.
func stableNetOptions(settings *backend.DataSourceInstanceSettings) *stablenet.ConnectOptions {
	if settings == nil {
		panic("datasource settings are nil, are you in a datasource environment?")
	}
	options := make(map[string]string)
	// error checking of unmarshal is done by checking the specific fields
	_ = json.Unmarshal(settings.JSONData, &options)
	if _, ok := options["snip"]; !ok {
		panic("field \"snip\" is missing in the JSONData of the datasource")
	}
	if _, ok := options["snport"]; !ok {
		panic("the field \"snport\" is missing in the JSONData of the datasource")
	}
	if _, ok := options["snusername"]; !ok {
		panic("the field \"snusername\" is missing in the JSONData of the datasource")
	}
	if _, ok := settings.DecryptedSecureJSONData["snpassword"]; !ok {
		panic("the field \"snpassword\" is missing in the encryptedJSONData of the datasource")
	}
	port, portErr := strconv.Atoi(options["snport"])
	if portErr != nil {
		panic(fmt.Sprintf("the field \"snport\" could not be parsed into a number: %v", portErr))
	}
	address := fmt.Sprintf("%s:%d", options["snip"], port)
	return &stablenet.ConnectOptions{
		Address:  address,
		Username: options["snusername"],
		Password: settings.DecryptedSecureJSONData["snpassword"],
	}
}
