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
	"github.com/hashicorp/go-hclog"
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
	MeasurementObid int   `json:"measurementObid"`
	MetricIds       []int `json:"metricIds"`
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

func (q *Query) GetCustomIntField(name string) (int, error) {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return 0, err
	}
	return queryJson.Get(name).Int()
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
	Logger   hclog.Logger
}

func (s *StableNetHandler) fetchMetrics(query Query, measurementObid int, valueIds []int) ([]*datasource.TimeSeries, error) {
	data, err := s.SnClient.FetchDataForMetrics(measurementObid, valueIds, query.StartTime, query.EndTime)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics from StableNet(R): %v", err)
	}
	keys := make([]string, 0, len(data))
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]*datasource.TimeSeries, 0, len(data))
	for _, name := range keys {
		series := data[name].ExpandWithMissingValues()
		maxTimeSeries := &datasource.TimeSeries{
			Points: series.MaxValues(),
			Name:   "Max " + name,
		}
		minTimeSeries := &datasource.TimeSeries{
			Points: series.MinValues(),
			Name:   "Min " + name,
		}
		avgTimeSeries := &datasource.TimeSeries{
			Points: series.AvgValues(),
			Name:   "Avg " + name,
		}
		if query.includeMinStats() {
			result = append(result, minTimeSeries)
		}
		if query.includeMaxStats() {
			result = append(result, maxTimeSeries)
		}
		if query.includeAvgStats() {
			result = append(result, avgTimeSeries)
		}
	}
	return result, nil
}
