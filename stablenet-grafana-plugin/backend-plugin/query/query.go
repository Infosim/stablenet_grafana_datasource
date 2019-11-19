package query

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"regexp"
	"strconv"
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
		if query.includeAvgStats() {
			result = append(result, avgTimeSeries)
		}
		if query.includeMaxStats() {
			result = append(result, maxTimeSeries)
		}
	}
	return result, nil
}

type Handler interface {
	Process(Query) (*datasource.QueryResult, error)
}

type DeviceHandler struct {
	*StableNetHandler
}

func (d DeviceHandler) Process(q Query) (*datasource.QueryResult, error) {
	deviceQuery, err := q.GetCustomField("deviceQuery")
	if err != nil {
		return BuildErrorResult("could not extract the deviceQuery from the query", q.RefId), nil
	}
	devices, err := d.SnClient.QueryDevices(deviceQuery)
	if err != nil {
		e := fmt.Errorf("could not retrieve devices from StableNet(R): %v", err)
		d.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(devices, q.RefId)
}

type MeasurementHandler struct {
	*StableNetHandler
}

func (m MeasurementHandler) Process(query Query) (*datasource.QueryResult, error) {
	deviceObid, err := query.GetCustomIntField("deviceObid")
	if err != nil {
		return BuildErrorResult("could not extract deviceObid from the query", query.RefId), nil
	}
	measurements, err := m.SnClient.FetchMeasurementsForDevice(deviceObid)
	if err != nil {
		e := fmt.Errorf("could not fetch measurements from StableNet(R): %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(measurements, query.RefId)
}

type MetricNameHandler struct {
	*StableNetHandler
}

func (m MetricNameHandler) Process(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	metrics, err := m.SnClient.FetchMetricsForMeasurement(measurementObid)
	if err != nil {
		e := fmt.Errorf("could not retrieve metric names from StableNet(R): %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(metrics, query.RefId)
}

type MetricDataHandler struct {
	*StableNetHandler
}

func (m MetricDataHandler) Process(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	metricId, err := query.GetCustomIntField("metricId")
	if err != nil {
		return BuildErrorResult("could not extract metricName from query", query.RefId), nil
	}

	series, err := m.fetchMetrics(query, measurementObid, []int{metricId})
	if err != nil {
		e := fmt.Errorf("could not fetch metric data from server: %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}

	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: series,
	}
	return &result, nil
}

type DatasourceTestHandler struct {
	*StableNetHandler
}

func (d DatasourceTestHandler) Process(query Query) (*datasource.QueryResult, error) {
	_, err := d.SnClient.FetchMeasurementsForDevice(-1)
	if err != nil {
		return BuildErrorResult("Cannot login into StableNet(R) with the provided credentials", query.RefId), nil
	}
	return &datasource.QueryResult{
		Series: []*datasource.TimeSeries{},
	}, nil
}

type StatisticLinkHandler struct {
	*StableNetHandler
}

func (s StatisticLinkHandler) Process(query Query) (*datasource.QueryResult, error) {
	link, err := query.GetCustomField("statisticLink")
	if err != nil {
		return BuildErrorResult("could not extract statisticLink parameter from query", query.RefId), nil
	}
	measurementRegex := regexp.MustCompile("[?&]id=(\\d+)")
	idMatches := measurementRegex.FindAllStringSubmatch(link, 1)
	if len(idMatches) == 0 {
		return BuildErrorResult(fmt.Sprintf("the link \"%s\" does not carry a measurement id.", link), query.RefId), nil
	}
	measurementId, _ := strconv.Atoi(idMatches[0][1])
	valueRegex := regexp.MustCompile("[?&]value\\d*=(\\d+)")
	valueMatches := valueRegex.FindAllStringSubmatch(link, -1)
	valueIds := make([]int, 0, len(valueMatches))
	for _, valueMatch := range valueMatches {
		id, _ := strconv.Atoi(valueMatch[1])
		valueIds = append(valueIds, id)
	}

	series, err := s.fetchMetrics(query, measurementId, valueIds)
	if err != nil {
		e := fmt.Errorf("could not fetch data for statistic link from server: %v", err)
		s.Logger.Error(e.Error())
		return nil, e
	}

	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: series,
	}
	return &result, nil
}

func createResponseWithCustomData(data interface{}, refId string) (*datasource.QueryResult, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal json: %v", err)
	}
	result := datasource.QueryResult{
		RefId:    refId,
		MetaJson: string(payload),
		Series:   []*datasource.TimeSeries{},
	}
	return &result, nil
}
