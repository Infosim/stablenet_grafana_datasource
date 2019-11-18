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

type QueryHandler interface {
	Handle(Query) *datasource.QueryResult
}

func NewHandler(logger hclog.Logger, snClient stablenet.Client, startTime time.Time, endTime time.Time) QueryHandler {
	return &queryHandlerImpl{
		logger:    logger,
		snClient:  snClient,
		startTime: startTime,
		endTime:   endTime,
	}
}

type queryHandlerImpl struct {
	logger    hclog.Logger
	snClient  stablenet.Client
	startTime time.Time
	endTime   time.Time
}

func (q *queryHandlerImpl) Handle(query Query) *datasource.QueryResult {
	queryType, queryTypeError := query.GetCustomField("queryType")
	if queryTypeError != nil {
		msg := fmt.Sprintf("could not retrieve query type: %v", queryTypeError)
		return BuildErrorResult(msg, query.RefId)
	}
	var result *datasource.QueryResult
	var queryError error
	switch queryType {
	case "devices":
		result, queryError = q.handleDeviceQuery(query)
		break
	case "measurements":
		result, queryError = q.handleMeasurementQuery(query)
		break
	case "metricNames":
		result, queryError = q.handleMetricNameQuery(query)
		break
	case "metricData":
		result, queryError = q.handleDataQuery(query)
		break
	case "testDatasource":
		result, queryError = q.handleDatasourceTest(query)
		break
	case "statisticLink":
		result, queryError = q.handleStatisticLink(query)
	default:
		msg := fmt.Sprintf("query type \"%s\" not supported", queryType)
		q.logger.Info(msg)
		return BuildErrorResult(msg, query.RefId)
	}

	if queryError != nil {
		q.logger.Error(queryError.Error())
		return BuildErrorResult("Internal Backend Plugin error. Please consult the Grafana logs.", query.RefId)
	}
	return result
}

func (q *queryHandlerImpl) handleDeviceQuery(query Query) (*datasource.QueryResult, error) {
	deviceQuery, err := query.GetCustomField("deviceQuery")
	if err != nil {
		return BuildErrorResult("could not extraxt the deviceQuery from the query", query.RefId), nil
	}
	devices, err := q.snClient.QueryDevices(deviceQuery)
	if err != nil {
		e := fmt.Errorf("could not retrieve devices from StableNet(R): %v", err)
		q.logger.Error("could not retrieve devices from StableNet(R)", e)
		return nil, e
	}
	return q.createResponseWithCustomData(devices, query.RefId)
}

func (q *queryHandlerImpl) createResponseWithCustomData(data interface{}, refId string) (*datasource.QueryResult, error) {
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

func (q *queryHandlerImpl) handleMeasurementQuery(query Query) (*datasource.QueryResult, error) {
	deviceObid, err := query.GetCustomIntField("deviceObid")
	if err != nil {
		return BuildErrorResult("could not extract deviceObid from the query", query.RefId), nil
	}
	measurements, err := q.snClient.FetchMeasurementsForDevice(deviceObid)
	if err != nil {
		return nil, fmt.Errorf("could not fetch measurements: %v", err)
	}
	return q.createResponseWithCustomData(measurements, query.RefId)
}

func (q *queryHandlerImpl) handleMetricNameQuery(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	metrics, err := q.snClient.FetchMetricsForMeasurement(measurementObid)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metric names from StableNet(R): %v", err)
	}
	return q.createResponseWithCustomData(metrics, query.RefId)
}

func (q *queryHandlerImpl) handleDataQuery(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	metricId, err := query.GetCustomIntField("metricId")
	if err != nil {
		return BuildErrorResult("could not extract metricName from query", query.RefId), nil
	}

	series, err := q.fetchMetrics(query, measurementObid, []int{metricId})
	if err != nil {
		return nil, fmt.Errorf("could not fetch data from server: %v", err)
	}

	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: series,
	}
	return &result, nil
}

func (q *queryHandlerImpl) fetchMetrics(query Query, measurementObid int, valueIds []int) ([]*datasource.TimeSeries, error) {
	data, err := q.snClient.FetchDataForMetrics(measurementObid, valueIds, q.startTime, q.endTime)
	q.logger.Error(fmt.Sprintf("%v", data))
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

func (q *queryHandlerImpl) handleDatasourceTest(query Query) (*datasource.QueryResult, error) {
	_, err := q.snClient.FetchMeasurementsForDevice(-1)
	if err != nil {
		return BuildErrorResult("Cannot login into StableNet(R) with the provided credentials", query.RefId), nil
	}
	return &datasource.QueryResult{
		Series: []*datasource.TimeSeries{},
	}, nil
}

func (q *queryHandlerImpl) handleStatisticLink(query Query) (*datasource.QueryResult, error) {
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

	series, err := q.fetchMetrics(query, measurementId, valueIds)
	if err != nil {
		return nil, fmt.Errorf("could not fetch data from server: %v", err)
	}

	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: series,
	}
	return &result, nil
}
