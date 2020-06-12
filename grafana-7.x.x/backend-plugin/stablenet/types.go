/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package stablenet

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"time"
)

type Device struct {
	Name string `json:"name"`
	Obid int    `json:"obid"`
}

type Measurement struct {
	Name string `json:"name"`
	Obid int    `json:"obid"`
}

type Metric struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type MetricData struct {
	Interval time.Duration
	Time     time.Time
	Min      float64
	Max      float64
	Avg      float64
}

type MetricDataSeries []MetricData

func (s MetricDataSeries) MinValues() []*datasource.Point {
	return s.toValues(func(data MetricData) float64 {
		return data.Min
	})
}

func (s MetricDataSeries) toValues(selector func(MetricData) float64) []*datasource.Point {
	result := make([]*datasource.Point, 0, len(s))
	for _, metricData := range s {
		point := datasource.Point{
			Timestamp: metricData.Time.UnixNano() / int64(1000*time.Microsecond),
			Value:     selector(metricData),
		}
		result = append(result, &point)
	}
	return result
}

func (s MetricDataSeries) MaxValues() []*datasource.Point {
	return s.toValues(func(data MetricData) float64 {
		return data.Max
	})
}

func (s MetricDataSeries) AvgValues() []*datasource.Point {
	return s.toValues(func(data MetricData) float64 {
		return data.Avg
	})
}

func (s MetricDataSeries) AsTable(max, min, avg bool) [][]interface{} {
	table := make([][]interface{}, 0, len(s))
	for _, data := range s {
		row := make([]interface{}, 0, 4)
		row = append(row, data.Time)
		if max {
			row = append(row, data.Max)
		}
		if min {
			row = append(row, data.Min)
		}
		if avg {
			row = append(row, data.Avg)
		}
		table = append(table, row)
	}
	return table
}

type ServerInfo struct {
	ServerVersion ServerVersion `xml:"serverversion"`
}

type ServerVersion struct {
	Version string `xml:"version,attr"`
}

type DataQuery struct {
	Start   int64    `json:"start"`
	End     int64    `json:"end"`
	Metrics []string `json:"metrics"`
	Raw     bool     `json:"raw"`
	Average int64    `json:"average"`
}

type DataQueryOptions struct {
	MeasurementObid int
	Metrics         []string
	Start           time.Time
	End             time.Time
	Average         int64
}
