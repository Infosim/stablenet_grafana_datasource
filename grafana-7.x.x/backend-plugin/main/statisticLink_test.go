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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStatisticLink(t *testing.T) {

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
