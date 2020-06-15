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
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
	"time"
)

func TestExpandStatisticLinks(t *testing.T) {
	queries := []MetricQuery{
		{
			Start:         time.Now(),
			End:           time.Now().Add(5 * time.Hour),
			Interval:      4000,
			StatisticLink: ptr("http://example.com/measurements/?0id=4000&0value1=4&0value0=2"),
		},
		{
			Start:           time.Now(),
			End:             time.Now().Add(5 * time.Hour),
			Interval:        4000,
			MeasurementObid: 2121,
			Metrics:         []StringPair{{Name: "Host Uptime", Key: "SNMP_10"}},
		},
	}
	metricProvider := func(i int) ([]stablenet.Metric, error) {
		if i == 4000 {
			return []stablenet.Metric{{Key: "SNMP_1", Name: "In"}, {Key: "SNMP_2", Name: "Out"}, {Key: "SNMP_3", Name: "Up"}, {Key: "SNMP_4", Name: "Down"}}, nil
		}
		return nil, fmt.Errorf("measurement %d not found", i)
	}
	t.Run("success", func(t *testing.T) {
		got, err := ExpandStatisticLinks(queries, metricProvider)
		require.Nil(t, err, "no error expected")
		require.Equal(t, 2, len(got), "expanded queries wrong")
		assert.Equal(t, 4000, got[0].MeasurementObid, "measurement obid of first query not correct")
		assert.Equal(t, 2121, got[1].MeasurementObid, "measurement obid of second query not correct")
	})
	t.Run("expand error", func(t *testing.T) {
		q := []MetricQuery{{StatisticLink: ptr("not a link")}}
		got, err := ExpandStatisticLinks(q, metricProvider)
		require.Nil(t, got, "should be nil in case of error")
		assert.EqualError(t, err, "could not parse statistic link of query 0: the link \"not a link\" does not carry at least a measurement id", "error message wrong")
	})
}

func TestParseStatisticLink(t *testing.T) {
	query := MetricQuery{
		Start:           time.Now(),
		End:             time.Now().Add(5 * time.Hour),
		Interval:        4000,
		IncludeAvgStats: true,
		IncludeMaxStats: true,
		IncludeMinStats: true,
		StatisticLink:   ptr("http://example.com/measurements/?0id=4000&1id=5000&0value1=4&0value0=2&1value0=23&2id=6000"),
	}
	metricProvider := func(i int) ([]stablenet.Metric, error) {
		if i == 4000 {
			return []stablenet.Metric{{Key: "SNMP_1", Name: "In"}, {Key: "SNMP_2", Name: "Out"}, {Key: "SNMP_3", Name: "Up"}, {Key: "SNMP_4", Name: "Down"}}, nil
		} else if i == 5000 {
			return []stablenet.Metric{{Key: "SNMP_2", Name: "Disc Space"}, {Key: "SCRIPT_23", Name: "VMs"}}, nil
		} else if i == 6000 {
			return []stablenet.Metric{}, nil
		}
		return nil, fmt.Errorf("measurement %d not found", i)
	}
	t.Run("success", func(t *testing.T) {
		got, err := parseStatisticLink(query, metricProvider)
		require.NoError(t, err, "no error expected")
		require.Equal(t, 2, len(got), "number of expanded queries")
		sort.Slice(got, func(i, j int) bool {
			return got[i].MeasurementObid < got[j].MeasurementObid
		})
		one := got[0]
		assert.Equal(t, query.Start, one.Start, "start of first query not correct")
		assert.Equal(t, query.End, one.End, "end of first query not correct")
		assert.Equal(t, query.Interval, one.Interval, "Interval of first query not correct")
		assert.Equal(t, query.IncludeMinStats, one.IncludeMinStats, "min of first query not correct")
		assert.Equal(t, query.IncludeMaxStats, one.IncludeMaxStats, "max of first query not correct")
		assert.Equal(t, query.IncludeAvgStats, one.IncludeAvgStats, "avg of first query not correct")
		assert.Equal(t, 4000, one.MeasurementObid, "measurementObid of first query not correct")
		assert.Equal(t, []StringPair{{Key: "SNMP_2", Name: "Out"}, {Key: "SNMP_4", Name: "Down"}}, one.Metrics, "metrics of first query not correct")
	})
	t.Run("carries no link", func(t *testing.T) {
		q := MetricQuery{StatisticLink: ptr("not a link")}
		got, err := parseStatisticLink(q, metricProvider)
		assert.Nil(t, got, "should be nil in case of an error")
		assert.EqualError(t, err, "the link \"not a link\" does not carry at least a measurement id")
	})
	t.Run("carries no link", func(t *testing.T) {
		q := MetricQuery{StatisticLink: ptr("&id=10000")}
		got, err := parseStatisticLink(q, metricProvider)
		assert.Nil(t, got, "should be nil in case of an error")
		assert.EqualError(t, err, "could not fetch metrics for measurement 10000: measurement 10000 not found")
	})
}

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
			assert.Equal(t, len(tt.wanted), len(got), "length is different")
			for key, value := range tt.wanted {
				gotValue, ok := got[key]
				if !assert.True(t, ok, "key %d not available in got map", key) {
					continue
				}
				assert.Equal(t, value, gotValue, "for key %d the expected value %d differs from the got one %d", key, value, gotValue)
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
			assert.Equal(t, len(tt.wanted), len(got), "lenght is different")
			for key, value := range tt.wanted {
				gotValue, ok := got[key]
				if !assert.True(t, ok, "key %d not available in got map", key) {
					continue
				}
				assert.Equal(t, value, gotValue, "for key %d the expected value %d differs from the got one %d", key, value, gotValue)
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
			assert.Equal(t, len(tt.wanted), len(got))
			for _, wantedIndex := range tt.wanted {
				wanted := tt.metrics[wantedIndex]
				contained := false
				for _, metric := range got {
					if metric.Key == wanted.Key && metric.Name == wanted.Name {
						contained = true
					}
				}
				assert.True(t, contained, "%v was not present in to got metrics", wanted)
			}
		})
	}
}

func ptr(value string) *string {
	result := value
	return &result
}
