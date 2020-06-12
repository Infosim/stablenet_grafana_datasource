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
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"sort"
	"time"
)

func BuildErrorResult(msg string, refId string) *datasource.QueryResult {
	return &datasource.QueryResult{
		Error: msg,
		RefId: refId,
	}
}

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

type Query struct {
	datasource.Query
	StartTime time.Time
	EndTime   time.Time
}

func (q *Query) GetCustomField(name string) (string, error) {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return "", err
	}
	return queryJson.Get(name).String()
}

func (q *Query) GetCustomIntField(name string) (*int, error) {
	var object map[string]interface{}
	err := json.Unmarshal([]byte(q.ModelJson), &object)
	if err != nil {
		return nil, err
	}
	if _, ok := object[name]; !ok {
		return nil, fmt.Errorf("value '%s' not present in the modelJson", name)
	}
	floatValue, ok := object[name].(float64)
	if !ok {
		return nil, fmt.Errorf("value '%s' is supposed to be an int, but was not", name)
	}
	intValue := int(floatValue)
	return &intValue, nil
}

func (q *Query) GetMeasurementDataRequest() ([]measurementDataRequest, error) {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return nil, fmt.Errorf("error while creating json from modelJson: %v", err)
	}
	if queryJson.Get("requestData").Interface() == nil {
		return nil, fmt.Errorf("dataRequest not present in the modelJson")
	}
	dataRequestBytes, err := queryJson.Get("requestData").Encode()
	var result []measurementDataRequest
	err = json.Unmarshal(dataRequestBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("requestData field of modelJson has not the expected format: %v", err)
	}
	return result, nil
}

func (q *Query) includeMinStats() bool {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return false
	}
	result, _ := queryJson.Get("includeMinStats").Bool()
	return result
}

func (q *Query) includeAvgStats() bool {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return false
	}
	result, _ := queryJson.Get("includeAvgStats").Bool()
	return result
}

func (q *Query) includeMaxStats() bool {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return false
	}
	result, _ := queryJson.Get("includeMaxStats").Bool()
	return result
}

type StableNetHandler struct {
	SnClient stablenet.Client
}

type MetricQuery struct {
	IncludeAvgStats bool
	IncludeMaxStats bool
	IncludeMinStats bool
	MeasurementObid int
	Metrics         []struct {
		Key  string
		Name string
	}
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

func (s *StableNetHandler) FetchMetrics(originalQuery backend.DataQuery, metricQuery MetricQuery) ([]*data.Frame, error) {
	options := stablenet.DataQueryOptions{
		MeasurementObid: metricQuery.MeasurementObid,
		Metrics:         metricQuery.metricKeys(),
		Start:           originalQuery.TimeRange.From,
		End:             originalQuery.TimeRange.To,
		Average:         int64(originalQuery.Interval / time.Millisecond),
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
