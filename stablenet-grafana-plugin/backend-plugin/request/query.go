package request

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"time"
)

type Query struct {
	datasource.Query
}

func BuildErrorResult(msg string, refId string) *datasource.QueryResult {
	return &datasource.QueryResult{
		Error: msg,
		RefId: refId,
	}
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
		q.logger.Error(msg)
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
	default:
		msg := fmt.Sprintf("query type \"%s\" not supported")
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
		return BuildErrorResult("could not extraxt the deviceQuery from the request", query.RefId), nil
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
		return BuildErrorResult("could not extract deviceObid from the request", query.RefId), nil
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
		return BuildErrorResult("could not extract measurementObid from request", query.RefId), nil
	}
	metrics, err := q.snClient.FetchMetricsForMeasurement(measurementObid, q.startTime, q.endTime)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metric names from StableNet(R): %v", err)
	}
	return q.createResponseWithCustomData(metrics, query.RefId)
}

func (q queryHandlerImpl) handleDataQuery(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from request", query.RefId), nil
	}
	metricName, err := query.GetCustomField("metricName")
	if err != nil {
		return BuildErrorResult("could not extract metricName from request", query.RefId), nil
	}
	data, err := q.snClient.FetchDataForMetric(measurementObid, metricName, q.startTime, q.endTime)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics from StableNet(R): %v", err)
	}

	points := make([]*datasource.Point, 0, len(data))
	for _, metricData := range data {
		p := datasource.Point{
			Timestamp: metricData.Time.UnixNano() / int64(1000*time.Microsecond),
			Value:     metricData.Value,
		}
		points = append(points, &p)
	}
	timeSeries := datasource.TimeSeries{
		Points: points,
	}
	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: []*datasource.TimeSeries{&timeSeries},
	}
	return &result, nil
}

func (q queryHandlerImpl) handleDatasourceTest(query Query) (*datasource.QueryResult, error) {
	_, err := q.snClient.FetchMeasurementsForDevice(-1)
	if err != nil {
		return BuildErrorResult("Cannot login into StableNet(R) with the provided credentials", query.RefId), nil
	}
	return &datasource.QueryResult{
		Series: []*datasource.TimeSeries{},
	}, nil
}
