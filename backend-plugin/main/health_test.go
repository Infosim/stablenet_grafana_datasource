/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"backend-plugin/mock"
	"backend-plugin/stablenet"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var stableNetUsername = "infosim"
var stableNetPassword = "stablenet"

func TestDataSource_CheckHealth(t *testing.T) {
	tests := []struct {
		name             string
		snVersion        string
		wantStatus       backend.HealthStatus
		wantLicenseError bool
		wantBody         string
	}{
		{name: "too old", snVersion: "8.5.0", wantStatus: backend.HealthStatusError, wantBody: "The StableNet® version 8.5.0 does not support Grafana®."},
		{name: "recent", snVersion: "9.0.0", wantStatus: backend.HealthStatusOk, wantBody: "Connection to StableNet® successful"},
		{name: "recent with productname should fail", snVersion: "StableNet 9.0.0", wantStatus: backend.HealthStatusError, wantBody: "The StableNet® version StableNet 9.0.0 does not support Grafana®."},
		{name: "future", snVersion: "10.1.0", wantStatus: backend.HealthStatusOk, wantBody: "Connection to StableNet® successful"},
		{name: "rest-reporting not licensed", snVersion: "9.0.2", wantLicenseError: true, wantStatus: backend.HealthStatusError, wantBody: "The StableNet® server does not have the required license \"rest-reporting\"."},
	}

	snServer := mock.CreateMockServer(stableNetUsername, stableNetPassword)
	handler := mock.CreateHandler(snServer)
	server := httptest.NewServer(handler)
	defer server.Close()

	byteData, _ := json.Marshal(map[string]string{
		"snusername": stableNetUsername,
		"snip":       "to be changed by context",
		"snport":     "5443",
	})
	secureJsonData := map[string]string{
		"snpassword": stableNetPassword,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snServer.Info.ServerVersion = stablenet.ServerVersion{Version: tt.snVersion}

			if tt.wantLicenseError {
				snServer.Info.License.Modules.Modules = []stablenet.Module{{Name: "report"}, {Name: "cloud"}}
			} else {
				snServer.Info.License.Modules.Modules = []stablenet.Module{{Name: "rest-reporting"}, {Name: "ha"}}
			}

			instanceSettings := backend.DataSourceInstanceSettings{
				ID:                      5,
				JSONData:                byteData,
				DecryptedSecureJSONData: secureJsonData,
			}

			healthReq := &backend.CheckHealthRequest{
				PluginContext: backend.PluginContext{
					DataSourceInstanceSettings: &instanceSettings,
				},
			}

			ds := dataSource{validationStore: make(map[int64]bool)}
			ctx := context.WithValue(context.Background(), "sn_address", server.URL)

			got, err := ds.CheckHealth(ctx, healthReq)

			require.Nil(t, err, "no error expected")
			assert.Equal(t, tt.wantStatus, got.Status, "status is wrong")

			if tt.wantStatus == backend.HealthStatusError {
				assert.False(t, ds.validationStore[5], "validationStore should be set to false")
			} else {
				assert.True(t, ds.validationStore[5], "validationStore should be set to true")
			}

			assert.Equal(t, tt.wantBody, got.Message, "response message not correct")
		})
	}
	t.Run("server error", func(t *testing.T) {
		secureJsonData["snpassword"] = "wrong"

		healthReq := &backend.CheckHealthRequest{
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
					ID:                      5,
					JSONData:                byteData,
					DecryptedSecureJSONData: secureJsonData,
				},
			},
		}
		ds := dataSource{validationStore: make(map[int64]bool)}
		ctx := context.WithValue(context.Background(), "sn_address", server.URL)
		got, err := ds.CheckHealth(ctx, healthReq)
		require.Nil(t, err, "the error should be nil")
		assert.Equal(t, backend.HealthStatusError, got.Status, "the health status is wrong")
		assert.Equal(t, "The StableNet® server could be reached, but the credentials were invalid.", got.Message, "the message is wrong")
	})
}
