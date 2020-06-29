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
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"sort"
	"time"
)

type MetricQuery struct {
	Start time.Time
	End   time.Time
	// Do not use `json:"intervalMs" here because this property gets overridden by Grafana.
	// We want to use our own average period.
	Interval        int64 `json:"customInterval"`
	IncludeAvgStats bool
	IncludeMaxStats bool
	IncludeMinStats bool
	StatisticLink   *string
	MeasurementObid int
	Metrics         []StringPair
}

func (m *MetricQuery) shallowClone() MetricQuery {
	return MetricQuery{
		Start:           m.Start,
		End:             m.End,
		Interval:        m.Interval,
		IncludeAvgStats: m.IncludeAvgStats,
		IncludeMaxStats: m.IncludeMaxStats,
		IncludeMinStats: m.IncludeMinStats,
		StatisticLink:   m.StatisticLink,
		MeasurementObid: m.MeasurementObid,
		Metrics:         m.Metrics,
	}
}

func NewQuery(query backend.DataQuery) MetricQuery {
	return MetricQuery{
		Start: query.TimeRange.From,
		End:   query.TimeRange.To,
	}
}

type StringPair struct {
	Key  string
	Name string
}

func (m *MetricQuery) metricKeys() []string {
	result := make([]string, 0, len(m.Metrics))
	for _, metric := range m.Metrics {
		result = append(result, metric.Key)
	}
	return result
}

func (m *MetricQuery) keyNameMap() map[string]string {
	result := make(map[string]string)
	for _, metric := range m.Metrics {
		result[metric.Key] = metric.Name
	}
	return result
}

func (m *MetricQuery) FetchData(provider func(stablenet.DataQueryOptions) (map[string]stablenet.MetricDataSeries, error)) ([]*data.Frame, error) {
	options := stablenet.DataQueryOptions{
		MeasurementObid: m.MeasurementObid,
		Metrics:         m.metricKeys(),
		Start:           m.Start,
		End:             m.End,
		Average:         m.Interval,
	}
	snData, err := provider(options)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics from StableNet(R): %v", err)
	}
	keys := make([]string, 0, len(snData))
	for key, _ := range snData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	frames := make([]*data.Frame, 0, len(snData))
	names := m.keyNameMap()
	for _, key := range keys {
		columns := make([]*data.Field, 0, 4)
		columns = append(columns, data.NewField("Time", nil, []time.Time{}))
		if m.IncludeMinStats {
			columns = append(columns, data.NewField("Min", nil, []float64{}))
		}
		if m.IncludeMaxStats {
			columns = append(columns, data.NewField("Max", nil, []float64{}))
		}
		if m.IncludeAvgStats {
			columns = append(columns, data.NewField("Avg", nil, []float64{}))
		}
		frame := data.NewFrame(names[key], columns...)
		for _, row := range snData[key].AsTable(m.IncludeMinStats, m.IncludeMaxStats, m.IncludeAvgStats) {
			frame.AppendRow(row...)
		}
		frames = append(frames, frame)
	}
	return frames, nil
}
