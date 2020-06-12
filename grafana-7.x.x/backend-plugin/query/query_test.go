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
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
			options := stablenet.DataQueryOptions{MeasurementObid: 1024, Metrics: []string{"123"}, Start: time.Time{}, End: time.Time{}}
			rawHandler.SnClient.(*mockSnClient).On("FetchDataForMetrics", options).Return(statisticResult, nil)
			metricsReq := []stablenet.Metric{{Name: "System Uptime", Key: "123"}}
			actual, err := rawHandler.fetchMetrics(query, 1024, metricsRequest(metricsReq))
			require.NoError(t, err, "no error expected")
			compareTimeSeries(t, tt.want, actual)
		})
	}
}
