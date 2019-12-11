/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package query

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildErrorResult(t *testing.T) {
	limerick := `Ein Limerickdichter aus Aachen,
nicht ahnend, was Limericks versprachen,
der trieb es zu bunt,
und das war der Grund,
dass Freunde zuletzt mit ihm brachen.`
	result := BuildErrorResult(limerick, "XYZ")
	assert := testify.New(t)
	assert.Equal(result.Error, limerick, "result error message wrong")
	assert.Equal(result.RefId, "XYZ", "result refId wrong")
	assert.Nil(result.Series, "series should be nil")
	assert.Empty(result.MetaJson, "meta json should be empty")
	assert.Nil(result.Tables, "tables should be nil")
}

func TestQuery_GetCustomField(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"favouriteDish\": \"all that is tasty\"}",
	}
	query := Query{Query: rawquery}
	t.Run("test successful", func(t *testing.T) {
		actual, err := query.GetCustomField("favouriteDish")
		require.NoError(t, err)
		testify.Equal(t, "all that is tasty", actual)
	})
	t.Run("test error", func(t *testing.T) {
		_, err := query.GetCustomField("favouriteMeal")
		testify.EqualError(t, err, "type assertion to string failed")
	})
}

func TestQuery_GetCustomFieldNoJson(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"favouriteDish\": \"all that is tasty\"",
	}
	query := Query{Query: rawquery}
	t.Run("test successful", func(t *testing.T) {
		_, err := query.GetCustomField("favouriteDish")
		require.EqualError(t, err, "unexpected EOF")
	})
}

func TestQuery_GetCustomIntField(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"age\": 5}",
	}
	query := Query{Query: rawquery}
	t.Run("test successful", func(t *testing.T) {
		actual, err := query.GetCustomIntField("age")
		require.NoError(t, err)
		testify.Equal(t, 5, *actual)
	})
	t.Run("test error", func(t *testing.T) {
		_, err := query.GetCustomIntField("birthYear")
		testify.EqualError(t, err, "value 'birthYear' not present in the modelJson")
	})
}

func TestQuery_GetCustomIntFieldNoJson(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"favouriteDish\": \"all that is tasty\"",
	}
	query := Query{Query: rawquery}
	t.Run("test successful", func(t *testing.T) {
		_, err := query.GetCustomIntField("favouriteDish")
		require.EqualError(t, err, "unexpected end of JSON input")
	})
}

func TestQuery_GetMeasurementDataRequest(t *testing.T) {
	metricsRequest1 := []stablenet.Metric{{Name: "Storage", Key: "5"}, {Name: "Free Storage", Key: "4"}, {Name: "Free Storage (%)", Key: "7"}}
	metricsRequest2 := []stablenet.Metric{{Name: "Uptime", Key: "26"}, {Name: "Users", Key: "24"}}
	data := []measurementDataRequest{
		{MeasurementObid: 1234, Metrics:metricsRequest(metricsRequest1)},
		{MeasurementObid: 6747, Metrics:metricsRequest(metricsRequest2)},
	}
	jsonBytes, _ := json.Marshal(data)
	rawquery := datasource.Query{
		ModelJson: "{\"requestData\": " + string(jsonBytes) + "}",
	}
	query := Query{Query: rawquery}
	result, err := query.GetMeasurementDataRequest()
	require.Nil(t, err, "no error expected to be thrown")
	require.Equal(t, data, result, "extracted data should be equal")
}

func TestQuery_GetMeasurementDataRequest_Error(t *testing.T) {
	tests := []struct {
		name      string
		modelJson string
		wantErr   string
	}{
		{name: "invalid model json", modelJson: "this is not a model json", wantErr: "error while creating json from modelJson: invalid character 'h' in literal true (expecting 'r')"},
		{name: "unknown format", modelJson: "{\"requestData\" : \"still no a json\"}", wantErr: "requestData field of modelJson has not the expected format: json: cannot unmarshal string into Go value of type []query.measurementDataRequest"},
		{name: "requestData not preset", modelJson: "{}", wantErr: "dataRequest not present in the modelJson"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := datasource.Query{
				ModelJson: tt.modelJson,
			}
			query := Query{Query: raw}
			result, err := query.GetMeasurementDataRequest()
			require.EqualError(t, err, tt.wantErr, "error message not equal")
			require.Nil(t, result, "no result expected because error should be returned")
		})
	}
}

func TestStableNetHandler_fetchMetrics(t *testing.T) {
	statisticResult, series := sampleStatisticData()
	tests := []struct {
		name            string
		includeMinStats bool
		includeMaxStats bool
		includeAvgStats bool
		want            []*datasource.TimeSeries
	}{
		{name: "no stats", includeMinStats: false, includeMaxStats: false, includeAvgStats: false, want: []*datasource.TimeSeries{}},
		{name: "all stats", includeMinStats: true, includeMaxStats: true, includeAvgStats: true, want: series},
		{name: "some", includeMinStats: true, includeMaxStats: false, includeAvgStats: true, want: []*datasource.TimeSeries{series[0], series[2]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawHandler, _ := setUpHandlerAndLogReceiver()
			requiredStats := map[string]bool{"includeMinStats": tt.includeMinStats, "includeMaxStats": tt.includeMaxStats, "includeAvgStats": tt.includeAvgStats}
			jsonQuery, _ := json.Marshal(&requiredStats)
			query := Query{
				Query: datasource.Query{ModelJson: string(jsonQuery)},
			}
			rawHandler.SnClient.(*mockSnClient).On("FetchDataForMetrics", 1024, []string{"123"}, time.Time{}, time.Time{}).Return(statisticResult, nil)
			metricsReq := []stablenet.Metric{{Name: "System Uptime", Key: "123"}}
			actual, err := rawHandler.fetchMetrics(query, 1024, metricsRequest(metricsReq))
			require.NoError(t, err, "no error expected")
			compareTimeSeries(t, tt.want, actual)
		})
	}
}
