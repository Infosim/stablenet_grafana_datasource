package main

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_request_stableNetOptionsErrors(t *testing.T) {
	tests := []struct {
		name              string
		jsonData          string
		decryptedJsonData map[string]string
		wantErr           string
	}{
		{name: "invalid json", jsonData: "{", decryptedJsonData: map[string]string{}, wantErr: "could not unmarshal jsonData of the datasource: unexpected end of JSON input"},
		{name: "missing snip", jsonData: "{}", decryptedJsonData: map[string]string{}, wantErr: "the snip is missing in the jsonData of the datasource"},
		{name: "missing snport", jsonData: "{\"snip\":\"127.0.0.1\"}", decryptedJsonData: map[string]string{}, wantErr: "the snport is missing in the jsonData of the datasource"},
		{name: "missing snusername", jsonData: "{\"snip\":\"127.0.0.1\", \"snport\":\"4444\"}", decryptedJsonData: map[string]string{}, wantErr: "the snusername is missing in the jsonData of the datasource"},
		{name: "missing snpassword", jsonData: "{\"snip\":\"127.0.0.1\", \"snport\":\"4444\", \"snusername\":\"infosim\"}", decryptedJsonData: map[string]string{}, wantErr: "the snpassword is missing in the encryptedJsonData of the datasource"},
		{name: "invalid snport", jsonData: "{\"snip\":\"127.0.0.1\", \"snport\": \"hello\", \"snusername\":\"infosim\"}", decryptedJsonData: map[string]string{"snpassword": "stablenet"}, wantErr: "could not parse snport into number: strconv.Atoi: parsing \"hello\": invalid syntax"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				DatasourceRequest: &datasource.DatasourceRequest{
					Datasource: &datasource.DatasourceInfo{
						JsonData:                tt.jsonData,
						DecryptedSecureJsonData: tt.decryptedJsonData,
					},
				},
			}
			got, err := r.stableNetOptions()
			require.Error(t, err, "must return a non-nil error")
			assert.Nil(t, got, "the options must be nil")
			assert.EqualError(t, err, tt.wantErr, "errors do not match")
		})
	}
}

func Test_request_stableNetOptions(t *testing.T) {
	json := "{\"snip\":\"127.0.0.1\", \"snport\": \"443\", \"snusername\":\"infosim\"}"
	decryptedData := map[string]string{"snpassword": "stablenet"}
	request := &request{
		DatasourceRequest: &datasource.DatasourceRequest{
			Datasource: &datasource.DatasourceInfo{
				JsonData:                json,
				DecryptedSecureJsonData: decryptedData,
			},
		},
	}
	actual, err := request.stableNetOptions()
	require.NoError(t, err, "no error is expected")
	require.NotNil(t, actual, "StableNet Options must not be nil")
	test := assert.New(t)
	test.Equal("127.0.0.1", actual.Host, "host differs")
	test.Equal(443, actual.Port, "port differs")
	test.Equal("infosim", actual.Username, "username differs")
	test.Equal("stablenet", actual.Password, "password differs")
}

func Test_request_timeRange(t *testing.T) {
	now := time.Now()
	then := now.Add(3 * time.Hour)
	nowRaw := now.UnixNano() / int64(time.Millisecond)
	thenRaw := then.UnixNano() / int64(time.Millisecond)
	request := &request{DatasourceRequest: &datasource.DatasourceRequest{TimeRange: &datasource.TimeRange{FromEpochMs: nowRaw, ToEpochMs: thenRaw}}}
	actualNow, actualThen := request.timeRange()
	assert.Equal(t, now.Second(), actualNow.Second(), "now differs")
	assert.Equal(t, then.Second(), actualThen.Second(), "then differs")
}
