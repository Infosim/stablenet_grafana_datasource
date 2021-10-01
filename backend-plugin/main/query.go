/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
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
	"strconv"
	"time"
)

type Mode int

const (
	Measurement   Mode = 0
	StatisticLink Mode = 10
)

/** Target describes the query coming directly from the frontend. It contains a lot of information which isn't needed
at all for querying data, but must be included in the Target in order to persist the whole config panel.
The unneeded elements aren't listed in the Go struct. The format of the Target is also not ideal for the data request,
but it's suited for the config panel itself. We convert the Target into a more suitable struct in the Target#toQuery method.
*/
type Target struct {
	Mode                Mode
	SelectedMeasurement struct {
		Value int
	} `json:"selectedMeasurement"`
	Interval         int64    `json:"customInterval"`
	ChosenMetrics    []string `json:"chosenMetrics"`
	MetricPrefix     string   `json:"metricPrefix"`
	IncludeMinStats  bool     `json:"includeMinStats"`
	IncludeAvgStats  bool     `json:"includeAvgStats"`
	IncludeMaxStats  bool     `json:"includeMaxStats"`
	StatisticLink    string   `json:"StatisticLink"`
	AveragePeriod    string   `json:"averagePeriod"`
	AverageUnit      int      `json:"averageUnit"`
	UseCustomAverage bool     `json:"useCustomAverage"`
	Metrics          []struct {
		Text string
		Key  string
	}
}

func (t *Target) toQuery(timeRange backend.TimeRange, refId string) MetricQuery {
	result := MetricQuery{
		Start: timeRange.From,
		End:   timeRange.To,
	}
	result.RefId = refId
	result.IncludeMinStats = t.IncludeMinStats
	result.IncludeAvgStats = t.IncludeAvgStats
	result.IncludeMaxStats = t.IncludeMaxStats
	period, err := strconv.Atoi(t.AveragePeriod)
	if t.UseCustomAverage && err == nil {
		result.Interval = int64(period * t.AverageUnit)
	} else {
		result.Interval = t.Interval
	}
	if t.Mode == StatisticLink && t.StatisticLink != "" {
		result.StatisticLink = &t.StatisticLink
	} else {
		result.MeasurementObid = t.SelectedMeasurement.Value
		metrics := make([]StringPair, 0, 0)
		for _, metric := range t.ChosenMetrics {
			for _, s := range t.Metrics {
				if s.Key == metric {
					metrics = append(metrics, StringPair{Key: metric, Name: fmt.Sprintf("%s %s", t.MetricPrefix, s.Text)})
				}
			}
		}
		result.Metrics = metrics
	}
	return result
}

type MetricQuery struct {
	Start           time.Time
	End             time.Time
	Interval        int64 `json:"customInterval"`
	IncludeAvgStats bool
	IncludeMaxStats bool
	IncludeMinStats bool
	StatisticLink   *string
	MeasurementObid int
	Metrics         []StringPair
	RefId           string
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
		RefId:           m.RefId,
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
