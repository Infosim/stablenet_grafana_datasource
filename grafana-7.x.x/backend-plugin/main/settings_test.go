/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStableNetSettings_Nil(t *testing.T) {
	panicFunc := func() {
		stableNetOptions(nil)
	}
	assert.PanicsWithValue(t, "datasource settings are nil, are you in a datasource environment?", panicFunc, "should panic with correct string")
}

func TestStableNetSettings(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		secureData  map[string]string
		wantPanic   bool
		panicString string
	}{
		{name: "no stablenet ip", json: "{}", secureData: nil, panicString: "field \"snip\" is missing in the JSONData of the datasource"},
		{name: "no stablenet port", json: "{\"snip\":\"55.66.77.88\"}", secureData: nil, panicString: "the field \"snport\" is missing in the JSONData of the datasource"},
		{name: "no stablenet user", json: "{\"snip\":\"55.66.77.88\", \"snport\": 12345}", secureData: nil, panicString: "the field \"snusername\" is missing in the JSONData of the datasource"},
		{name: "no stablenet password", json: "{\"snip\":\"55.66.77.88\", \"snport\": 12345, \"snusername\":\"infosim\"}", secureData: nil, panicString: "the field \"snpassword\" is missing in the encryptedJSONData of the datasource"},
		{name: "stablenet invalid port", json: "{\"snip\":\"55.66.77.88\", \"snport\": \"not port\", \"snusername\":\"infosim\"}", secureData: map[string]string{"snpassword": "stablenet"}, panicString: "the field \"snport\" could not be parsed into a number: strconv.Atoi: parsing \"not port\": invalid syntax"},
		{name: "success", json: "{\"snip\":\"55.66.77.88\", \"snport\": \"12345\", \"snusername\":\"infosim\"}", secureData: map[string]string{"snpassword": "stablenet"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawMessage := json.RawMessage{}
			require.NoError(t, json.Unmarshal([]byte(tt.json), &rawMessage))
			settings := &backend.DataSourceInstanceSettings{JSONData: rawMessage, DecryptedSecureJSONData: tt.secureData}
			if len(tt.panicString) != 0 {
				panicFunc := func() {
					stableNetOptions(settings)
				}
				assert.PanicsWithValue(t, tt.panicString, panicFunc, "should panic with correct string")
			} else {
				options := stableNetOptions(settings)
				assert.Equal(t, "infosim", options.Username, "username not correct")
				assert.Equal(t, "stablenet", options.Password, "password not correct")
				assert.Equal(t, "https://55.66.77.88:12345", options.Address, "host not correct")
			}
		})
	}
}
