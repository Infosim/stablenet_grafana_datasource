package query

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"time"
)

func BuildErrorResult(msg string, refId string) *datasource.QueryResult {
	return &datasource.QueryResult{
		Error: msg,
		RefId: refId,
	}
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

func (q *Query) GetCustomIntArray(name string) ([]int, error) {
	queryJson, err := simplejson.NewJson([]byte(q.ModelJson))
	if err != nil {
		return nil, err
	}
	array, err := queryJson.Get(name).Array()
	if err != nil {
		return nil, err
	}
	result := make([]int, 0, len(array))
	for _, value := range array {
		intVal, ok := value.(json.Number)
		if !ok {
			return nil, fmt.Errorf("the value %v is not an integer", value)
		}
		realInt, _ := intVal.Int64()
		result = append(result, int(realInt))
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
	result := make([]*datasource.TimeSeries, 0, len(data))
	for name, series := range data {
		series = series.ExpandWithMissingValues()
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
