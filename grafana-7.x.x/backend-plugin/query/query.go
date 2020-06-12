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
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"sort"
	"time"
)

type measurementDataRequest struct {
	MeasurementObid int            `json:"measurementObid"`
	Metrics         metricsRequest `json:"metrics"`
}

type metricsRequest []stablenet.Metric

func (m metricsRequest) metricKeys() []string {
	result := make([]string, 0, len(m))
	for _, metric := range m {
		result = append(result, metric.Key)
	}
	return result
}

func (m metricsRequest) keyNameMap() map[string]string {
	result := make(map[string]string)
	for _, metric := range m {
		result[metric.Key] = metric.Name
	}
	return result
}

type StableNetHandler struct {
	SnClient stablenet.Client
}

type MetricQuery struct {
	Start           time.Time
	End             time.Time
	Interval        time.Duration
	IncludeAvgStats bool
	IncludeMaxStats bool
	IncludeMinStats bool
	StatisticLink   *string
	MeasurementObid int
	Metrics         []StringPair
}

func (m *MetricQuery) ShallowClone() MetricQuery {
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
		Start:    query.TimeRange.From,
		End:      query.TimeRange.To,
		Interval: query.Interval,
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

func (s *StableNetHandler) ExpandStatisticLinks(queries []MetricQuery) ([]MetricQuery, error) {
	result := make([]MetricQuery, 0, len(queries))
	for index, query := range queries {
		if query.StatisticLink == nil {
			result = append(result, query)
			continue
		}
		linkQueries, err := ParseStatisticLink(query, s.SnClient.FetchMetricsForMeasurement)
		if err != nil {
			return nil, fmt.Errorf("could not parse statistic link of query %d: %v", index, err)
		}
		for _, linkQuery := range linkQueries {
			result = append(result, linkQuery)
		}
	}
	return result, nil
}

func (s *StableNetHandler) FetchMetrics(metricQuery MetricQuery) ([]*data.Frame, error) {
	options := stablenet.DataQueryOptions{
		MeasurementObid: metricQuery.MeasurementObid,
		Metrics:         metricQuery.metricKeys(),
		Start:           metricQuery.Start,
		End:             metricQuery.End,
		Average:         int64(metricQuery.Interval / time.Millisecond),
	}
	snData, err := s.SnClient.FetchDataForMetrics(options)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics from StableNet(R): %v", err)
	}
	keys := make([]string, 0, len(snData))
	for key, _ := range snData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	frames := make([]*data.Frame, 0, len(snData))
	names := metricQuery.keyNameMap()
	for _, key := range keys {
		columns := make([]*data.Field, 0, 4)
		columns = append(columns, data.NewField("timeValues", nil, []time.Time{}))
		if metricQuery.IncludeMaxStats {
			columns = append(columns, data.NewField("Max", nil, []float64{}))
		}
		if metricQuery.IncludeMinStats {
			columns = append(columns, data.NewField("Min", nil, []float64{}))
		}
		if metricQuery.IncludeAvgStats {
			columns = append(columns, data.NewField("Avg", nil, []float64{}))
		}
		frame := data.NewFrame(names[key], columns...)
		for _, row := range snData[key].AsTable(metricQuery.IncludeMaxStats, metricQuery.IncludeMinStats, metricQuery.IncludeAvgStats) {
			frame.AppendRow(row...)
		}
		frames = append(frames, frame)
	}

	return frames, nil
}
