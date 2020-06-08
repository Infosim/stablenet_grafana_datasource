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
	"bufio"
	"bytes"
	"github.com/bmizerany/assert"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
	"time"
)

func TestFindMeasurementIdsInLink(t *testing.T) {
	cases := []struct {
		name   string
		link   string
		wanted map[int]int
	}{
		{name: "one without index 1", link: "stablenet.de/?id=33", wanted: map[int]int{0: 33}},
		{name: "one without index 2", link: "stablenet.de/?chart=555&id=34", wanted: map[int]int{0: 34}},
		{name: "several measurement ids 1", link: "stablenet.de/?chart=555&0id=34&0value1=1000&0value1=2000&1id=56&1value0=1001", wanted: map[int]int{0: 34, 1: 56}},
		{name: "several measurement ids 2", link: "stablenet.de/?chart=555&1id=34&1value1=1000&0value1=2000&0id=56&1value0=1001", wanted: map[int]int{1: 34, 0: 56}},
		{name: "several measurement ids mixed", link: "stablenet.de/?chart=555&id=34&0value1=1000&0value1=2000&1id=56&1value0=1001", wanted: map[int]int{0: 34, 1: 56}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := findMeasurementIdsInLink(tt.link)
			testify.Equal(t, len(tt.wanted), len(got), "length is different")
			for key, value := range tt.wanted {
				gotValue, ok := got[key]
				if !testify.True(t, ok, "key %d not available in got map", key) {
					continue
				}
				testify.Equal(t, value, gotValue, "for key %d the expected value %d differs from the got one %d", key, value, gotValue)
			}
		})
	}
}

func TestExtractMetricKeysForMeasurement(t *testing.T) {
	cases := []struct {
		name   string
		link   string
		wanted map[int][]string
	}{
		{name: "one measurement, one metric", link: "https://localhost:5443/PlotServlet?id=1643&chart=5504&last=0,1440&offset=0,0&tz=Europe%2FBerlin&value0=1002", wanted: map[int][]string{1643: {"1002"}}},
		{name: "one measurement, several metrics", link: "https://localhost:5443/PlotServlet?id=1643&chart=5504&last=0,1440&offset=0,0&tz=Europe%2FBerlin&value0=1002&value1=1000&value2=1001", wanted: map[int][]string{1643: {"1002", "1000", "1001"}}},
		{name: "several measurements", link: "https://localhost:5443/PlotServlet?multicharttype=0&dns=1&log=0&width=1252&height=1126&quality=-1.0&0last=0,1440&0offset=0,0&0interval=60000&0id=1643&0chart=5504&0value0=1000&0value1=1001&0value2=1002&1last=0,1440&1offset=0,0&1interval=60000&1id=3889&1chart=5504&1value0=1", wanted: map[int][]string{1643: {"1000", "1001", "1002"}, 3889: {"1"}}},
		{name: "ping measurements w/o valueIds", link: "https://localhost:5443/PlotServlet?multicharttype=0&dns=1&log=0&width=1252&height=1126&quality=-1.0&0last=0,1440&0offset=0,0&0interval=60000&0id=3088&0chart=100&1last=0,1440&1offset=0,0&1interval=60000&1id=7228&1chart=100&tz=Europe/Berlin", wanted: map[int][]string{3088: {}, 7228: {}}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMetricKeysForMeasurements(tt.link)
			testify.Equal(t, len(tt.wanted), len(got), "lenght is different")
			for key, value := range tt.wanted {
				gotValue, ok := got[key]
				if !testify.True(t, ok, "key %d not available in got map", key) {
					continue
				}
				testify.Equal(t, value, gotValue, "for key %d the expected value %d differs from the got one %d", key, value, gotValue)
			}
		})
	}
}

func TestFilterWantedMetrics(t *testing.T) {
	cases := []struct {
		name          string
		metrics       []stablenet.Metric
		wantedMetrics []string
		wanted        []int
	}{
		{name: "Wanted Metrics 1", metrics: []stablenet.Metric{{Name: "Uptime", Key: "SNMP1001"}, {Name: "Processes", Key: "SNMP1000"}, {Name: "Users", Key: "SNMP1002"}}, wantedMetrics: []string{"1001", "1002"}, wanted: []int{0, 2}},
		{name: "Wanted Metrics 2", metrics: []stablenet.Metric{{Name: "Uptime", Key: "SNMP1010"}, {Name: "Processes", Key: "SNMP1020"}}, wantedMetrics: []string{"1000"}, wanted: []int{}},
		{name: "Wanted Metrics is empty", metrics: []stablenet.Metric{{Name: "Uptime", Key: "SNMP1010"}, {Name: "Processes", Key: "SNMP1020"}}, wantedMetrics: []string{}, wanted: []int{0, 1}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := filterWantedMetrics(tt.wantedMetrics, tt.metrics)
			testify.Equal(t, len(tt.wanted), len(got))
			for _, wantedIndex := range tt.wanted {
				wanted := tt.metrics[wantedIndex]
				contained := false
				for _, metric := range got {
					if metric.Key == wanted.Key && metric.Name == wanted.Name {
						contained = true
					}
				}
				testify.True(t, contained, "%v was not present in to got metrics", wanted)
			}
		})
	}
}

func Test_statisticLinkHandler2_Successful(t *testing.T) {
	startTime := time.Now().Add(-5 * time.Hour)
	endTime := time.Now()

	metrics1643 := []stablenet.Metric{{Key: "SNMP1000", Name: "Uptime"}, {Key: "SNMP1001", Name: "Processes"}, {Key: "SNMP1002", Name: "Users"}}
	metrics3886 := []stablenet.Metric{{Key: "PING1", Name: "Ping"}}

	data1643 := map[string]stablenet.MetricDataSeries{"SNMP1000": {{Min: 5, Max: 7, Avg: 6, Interval: 300000, Time: endTime.Add(-1 * time.Hour)}}, "SNMP1002": {{Min: 1, Max: 1, Avg: 1, Interval: 300000, Time: endTime.Add(-1 * time.Hour)}}}
	data3886 := map[string]stablenet.MetricDataSeries{"PING1": {{Min: 300, Max: 400, Avg: 350, Interval: 300000, Time: endTime.Add(-1 * time.Hour)}}}

	name1643 := "ThinkStation Host"
	name3886 := "ThinkStation Ping"

	options1 := stablenet.DataQueryOptions{
		MeasurementObid: 1643,
		Metrics:         []string{"SNMP1000", "SNMP1002"},
		Start:           startTime,
		End:             endTime,
		Average:         250,
	}
	options2 := stablenet.DataQueryOptions{
		MeasurementObid: 3886,
		Metrics:         []string{"PING1"},
		Start:           startTime,
		End:             endTime,
		Average:         250,
	}

	client := new(mockSnClient)
	client.On("FetchMetricsForMeasurement", 1643, "").Return(metrics1643, nil)
	client.On("FetchMetricsForMeasurement", 3886, "").Return(metrics3886, nil)
	client.On("FetchDataForMetrics", options1).Return(data1643, nil)
	client.On("FetchDataForMetrics", options2).Return(data3886, nil)
	client.On("FetchMeasurementName", 1643).Return(&name1643, nil)
	client.On("FetchMeasurementName", 3886).Return(&name3886, nil)

	link := "https://localhost:5443/PlotServlet?multicharttype=0&dns=1&log=0&width=1252&height=1126&quality=-1.0&0last=0,1440&0offset=0,0&0interval=60000&0id=1643&0chart=5504&0value0=1000&0value1=1002&1last=0,1440&1offset=0,0&1interval=60000&1id=3886&1chart=5504&1value0=1"
	logData := bytes.Buffer{}
	logReceiver := bufio.NewWriter(&logData)
	snHandler := StableNetHandler{SnClient: client, Logger: hclog.New(&hclog.LoggerOptions{Output: logReceiver, TimeFormat: "no time"})}
	handler := statisticLinkHandler{StableNetHandler: &snHandler}
	query := Query{Query: datasource.Query{RefId: "A", ModelJson: "{\"includeAvgStats\": true, \"statisticLink\": \"" + link + "\"}", IntervalMs: 250}, StartTime: startTime, EndTime: endTime}
	got, err := handler.Process(query)
	require.NoError(t, err, "no error expected")
	assert.Equal(t, "A", got.RefId, "refId is wrong")
	series := got.Series
	sort.Slice(series, func(i, j int) bool {
		return len(series[i].Name) < len(series[j].Name)
	})
	require.Equal(t, 3, len(series), "number of series is wrong")

	assert.Equal(t, "ThinkStation Ping Avg Ping", series[0].Name, "name of first series wrong")
	assert.Equal(t, 1, len(series[0].Points), "number of data points of first series wrong")
	assert.Equal(t, 350.0, series[0].Points[0].Value, "value of data of first series wrong")
	assert.Equal(t, "ThinkStation Host Avg Users", series[1].Name, "name of second series wrong")
	assert.Equal(t, 1, len(series[1].Points), "number of data points of second series wrong")
	assert.Equal(t, 1.0, series[1].Points[0].Value, "value of data of second series wrong")
	assert.Equal(t, "ThinkStation Host Avg Uptime", series[2].Name, "name of third series wrong")
	assert.Equal(t, 1, len(series[2].Points), "number of data points of third series wrong")
	assert.Equal(t, 6.0, series[2].Points[0].Value, "value of data of third series wrong")
}
