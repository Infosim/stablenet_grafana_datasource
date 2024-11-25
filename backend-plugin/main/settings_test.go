/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadStableNetSettings_Nil(t *testing.T) {
	_, err := loadStableNetSettings(nil)
	assert.NotNil(t, err)
}

func TestLoadStableNetSettings_NoPassword(t *testing.T) {
	_, err := loadStableNetSettings(
		&backend.DataSourceInstanceSettings{
			URL:                     testStableNetUrl,
			User:                    testStableNetUsername,
			DecryptedSecureJSONData: nil,
		},
	)

	require.NotNil(t, err)
}

func TestLoadStableNetSettings(t *testing.T) {
	values := backend.DataSourceInstanceSettings{
		URL:                     testStableNetUrl,
		User:                    testStableNetUsername,
		DecryptedSecureJSONData: map[string]string{"password": testStableNetPassword},
	}

	options, err := loadStableNetSettings(&values)

	assert.NoError(t, err)

	assert.Equal(t, testStableNetUrl, options.Address, "host not correct")
	assert.Equal(t, testStableNetUsername, options.Username, "username not correct")
	assert.Equal(t, testStableNetPassword, options.Password, "password not correct")
}
