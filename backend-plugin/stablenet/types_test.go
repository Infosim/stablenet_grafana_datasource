/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package stablenet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExampleModules_IsRestReportingLicensed(t *testing.T) {
	modules := Modules{Modules: []Module{{Name: "nqa"}, {Name: "policy"}}}
	assert.False(t, modules.IsRestReportingLicensed())

	modules = Modules{Modules: []Module{{Name: "nqa"}, {Name: "policy"}, {Name: "rest-reporting"}}}
	assert.True(t, modules.IsRestReportingLicensed())
}

func TestMetricDataSeries_AsTable(t *testing.T) {
	now := time.Now()
	five := now.Add(5 * time.Minute)
	ten := now.Add(10 * time.Minute)
	series := MetricDataSeries{
		{
			Interval: 5000,
			Time:     now,
			Min:      1,
			Max:      101,
			Avg:      11},
		{
			Interval: 5000,
			Time:     five,
			Min:      2,
			Max:      102,
			Avg:      12},
		{
			Interval: 500,
			Time:     ten,
			Min:      0,
			Max:      100,
			Avg:      10},
	}
	tests := []struct {
		name string
		min  bool
		max  bool
		avg  bool
		want [][]interface{}
	}{
		{name: "all false", want: [][]interface{}{{now}, {five}, {ten}}},
		{name: "min", min: true, want: [][]interface{}{{now, 1.0}, {five, 2.0}, {ten, 0.0}}},
		{name: "min,max", min: true, max: true, want: [][]interface{}{{now, 1.0, 101.0}, {five, 2.0, 102.0}, {ten, 0.0, 100.0}}},
		{name: "min,max,avg", min: true, max: true, avg: true, want: [][]interface{}{{now, 1.0, 101.0, 11.0}, {five, 2.0, 102.0, 12.0}, {ten, 0.0, 100.0, 10.0}}},
		{name: "min,avg", min: true, avg: true, want: [][]interface{}{{now, 1.0, 11.0}, {five, 2.0, 12.0}, {ten, 0.0, 10.0}}},
		{name: "max,avg", max: true, avg: true, want: [][]interface{}{{now, 101.0, 11.0}, {five, 102.0, 12.0}, {ten, 100.0, 10.0}}},
		{name: "avg", avg: true, want: [][]interface{}{{now, 11.0}, {five, 12.0}, {ten, 10.0}}},
		{name: "max", max: true, want: [][]interface{}{{now, 101.0}, {five, 102.0}, {ten, 100.0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := series.AsTable(tt.min, tt.max, tt.avg)
			assert.Equal(t, tt.want, got, "computed table is wrong")
		})
	}
}
