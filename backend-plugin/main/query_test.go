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
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShallowClone(t *testing.T) {
	url := "http://example.org"
	original := MetricQuery{
		Start:           time.Now(),
		End:             time.Now().Add(4 * time.Hour),
		Interval:        90000,
		IncludeAvgStats: true,
		IncludeMaxStats: true,
		IncludeMinStats: true,
		StatisticLink:   &url,
		MeasurementObid: 232,
		Metrics:         []StringPair{{Key: "SNMP_10", Name: "Host"}},
	}
	got := original.shallowClone()
	assert.Equal(t, original, got, "clone should be equal to original")
	assert.Equal(t, &original.Metrics, &got.Metrics, "clone should be shallow, metric slice should be the same")
}

func TestMetricQuery_metricKeys(t *testing.T) {
	pairs := []StringPair{{Name: "Berlin", Key: "3232"}, {Name: "Dallas", Key: "343"}, {Name: "Moscow", Key: "4545"}}
	query := MetricQuery{Metrics: pairs}
	keys := query.metricKeys()
	assert.Equal(t, []string{"3232", "343", "4545"}, keys, "keys are different")
}

func TestMetricQuery_keyNameMap(t *testing.T) {
	pairs := []StringPair{{Name: "Berlin", Key: "3232"}, {Name: "Dallas", Key: "343"}, {Name: "Moscow", Key: "4545"}}
	query := MetricQuery{Metrics: pairs}
	got := query.keyNameMap()
	assert.Equal(t, map[string]string{"3232": "Berlin", "4545": "Moscow", "343": "Dallas"}, got, "generated map not correct")
}

func TestMetricQuery_FetchData_Error(t *testing.T) {
	query := MetricQuery{}
	got, err := query.FetchData(func(options stablenet.DataQueryOptions) (map[string]stablenet.MetricDataSeries, error) {
		return nil, errors.New("internal error for testing")
	})
	assert.EqualError(t, err, "could not retrieve metrics from StableNet(R): internal error for testing", "error message wrong")
	assert.Nil(t, got, "got should be nil in case of an error")
}

func TestMetricQuery_FetchData(t *testing.T) {
	now := time.Now()
	five := time.Now().Add(5 * time.Minute)
	tests := []struct {
		min                bool
		max                bool
		avg                bool
		wantReadsFirstLine []interface{}
		wantHeader         []string
	}{
		{min: false, max: false, avg: false, wantReadsFirstLine: []interface{}{now}, wantHeader: []string{"Time"}},
		{min: false, max: false, avg: true, wantReadsFirstLine: []interface{}{now, 8.0}, wantHeader: []string{"Time", "Avg"}},
		{min: false, max: true, avg: false, wantReadsFirstLine: []interface{}{now, 10.0}, wantHeader: []string{"Time", "Max"}},
		{min: false, max: true, avg: true, wantReadsFirstLine: []interface{}{now, 10.0, 8.0}, wantHeader: []string{"Time", "Max", "Avg"}},
		{min: true, max: false, avg: false, wantReadsFirstLine: []interface{}{now, 6.0}, wantHeader: []string{"Time", "Min"}},
		{min: true, max: false, avg: true, wantReadsFirstLine: []interface{}{now, 6.0, 8.0}, wantHeader: []string{"Time", "Min", "Avg"}},
		{min: true, max: true, avg: false, wantReadsFirstLine: []interface{}{now, 6.0, 10.0}, wantHeader: []string{"Time", "Min", "Max"}},
		{min: true, max: true, avg: true, wantReadsFirstLine: []interface{}{now, 6.0, 10.0, 8.0}, wantHeader: []string{"Time", "Min", "Max", "Avg"}},
	}
	writes := stablenet.MetricDataSeries{{
		Time: now.Add(time.Minute),
		Min:  5,
		Max:  9,
		Avg:  7,
	}, {
		Time: five,
		Min:  7,
		Max:  13,
		Avg:  11,
	}}
	reads := stablenet.MetricDataSeries{{
		Time: now,
		Min:  6,
		Max:  10,
		Avg:  8,
	}}
	frameHeader := func(frame *data.Frame) []string {
		result := make([]string, 0, len(frame.Fields))
		for _, field := range frame.Fields {
			result = append(result, field.Name)
		}
		return result
	}
	for index, tt := range tests {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			query := MetricQuery{
				Start:           time.Now(),
				End:             time.Now().Add(5 * time.Minute),
				Interval:        25000,
				IncludeAvgStats: tt.avg,
				IncludeMaxStats: tt.max,
				IncludeMinStats: tt.min,
				StatisticLink:   nil,
				MeasurementObid: 2342,
				Metrics:         []StringPair{{Key: "SNMP_10", Name: "Writes"}, {Key: "SNMP_20", Name: "Reads"}},
			}
			got, err := query.FetchData(func(options stablenet.DataQueryOptions) (map[string]stablenet.MetricDataSeries, error) {
				assert.Equal(t, query.Start, options.Start, "start option must be set correctly")
				assert.Equal(t, query.End, options.End, "end option must be set correctly")
				assert.Equal(t, query.MeasurementObid, options.MeasurementObid, "measurement obid option must be set correctly")
				assert.Equal(t, []string{"SNMP_10", "SNMP_20"}, options.Metrics, "metrics option must be set correctly")
				assert.Equal(t, query.Interval, options.Average, "average option must be set correctly")
				return map[string]stablenet.MetricDataSeries{"SNMP_10": writes, "SNMP_20": reads}, nil
			})
			assert.Nil(t, err, "no error expected")
			assert.Equal(t, 2, len(got), "number of frames should be equal to number of requested metrics")
			assert.Equal(t, tt.wantHeader, frameHeader(got[0]), "frame header of first frame not correct")
			assert.Equal(t, "Writes", got[0].Name, "name of first frame")
			assert.Equal(t, 2, got[0].Rows(), "number of rows in first frame")
			assert.Equal(t, "Reads", got[1].Name, "name of second frame")
			assert.Equal(t, 1, got[1].Rows(), "number of rows in second frame")
			assert.Equal(t, tt.wantReadsFirstLine, got[1].RowCopy(0), "first line of second frame wrong")
			assert.Equal(t, tt.wantHeader, frameHeader(got[0]), "frame header of second frame not correct")
		})
	}
}
