package main

import (
	"backend-plugin/request"
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"time"

	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

type JsonDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

func (j *JsonDatasource) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	req := &request.Content{*tsdbReq}
	queryType, err := req.GetCustomField("queryType")
	if err != nil {
		j.logger.Error("could not retrieve query type: %v", err)
		return nil, err
	}
	if queryType == "devices" {
		return j.handleDeviceQuery(req)
	}
	if queryType == "measurements" {
		return j.handleMeasurementQuery(req)
	}
	if queryType == "metricNames" {
		return j.handleMetricNameQuery(req)
	}
	if queryType == "metricData" {
		return j.handleDataQuery(req)
	}
	j.logger.Error("Query", "datasource", tsdbReq.Datasource.Name, "TimeRange", tsdbReq.TimeRange)
	then := time.Now().AddDate(-1, 0, 0)
	points := make([]*datasource.Point, 0, 0)
	for i := 0; i < 10; i++ {
		point := datasource.Point{
			Timestamp: then.UnixNano() / int64(time.Millisecond),
			Value:     float64(i * 1000),
		}
		points = append(points, &point)
		then = then.Add(-time.Hour)
	}
	timeSeries := datasource.TimeSeries{
		Name:   "Test Series",
		Tags:   nil,
		Points: points,
	}
	queryResult := datasource.QueryResult{
		Error:    "",
		RefId:    "A",
		MetaJson: "",
		Series:   []*datasource.TimeSeries{&timeSeries},
		Tables:   nil,
	}
	response := &datasource.DatasourceResponse{
		Results: []*datasource.QueryResult{&queryResult},
	}
	j.logger.Error("Context", fmt.Sprintf("%v", ctx))
	return response, nil
}

func (j *JsonDatasource) getQueryType(req *request.Content) (string, error) {
	queryType := "query"
	if len(req.Queries) > 0 {
		firstQuery := req.Queries[0]
		queryJson, err := simplejson.NewJson([]byte(firstQuery.ModelJson))
		if err != nil {
			return "", err
		}
		queryType = queryJson.Get("queryType").MustString("devices", "measurements")
	}
	return queryType, nil
}

func (j *JsonDatasource) handleDeviceQuery(req *request.Content) (*datasource.DatasourceResponse, error) {
	snClient := stablenet.NewClient(stablenet.ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	devices, err := snClient.FetchAllDevices()
	if err != nil {
		j.logger.Error("could not retrieve devices from StableNet(R)", err)
		return nil, err
	}
	return j.createResponseWithCustomData(devices)
}

func (j *JsonDatasource) createResponseWithCustomData(data interface{}) (*datasource.DatasourceResponse, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		j.logger.Error("could not marshal json", err)
		return nil, err
	}
	result := datasource.QueryResult{
		RefId:    "A",
		MetaJson: string(payload),
		Series:   []*datasource.TimeSeries{},
	}
	response := datasource.DatasourceResponse{
		Results: []*datasource.QueryResult{&result},
	}
	return &response, nil
}

func (j *JsonDatasource) handleMeasurementQuery(req *request.Content) (*datasource.DatasourceResponse, error) {
	snClient := stablenet.NewClient(stablenet.ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	deviceObid, err := req.GetCustomIntField("deviceObid")
	if err != nil {
		j.logger.Error(err.Error())
		return nil, err
	}
	measurements, err := snClient.FetchMeasurementsForDevice(deviceObid)
	if err != nil {
		j.logger.Error(err.Error())
		return nil, err
	}
	return j.createResponseWithCustomData(measurements)
}

func (j *JsonDatasource) handleMetricNameQuery(req *request.Content) (*datasource.DatasourceResponse, error) {
	snClient := stablenet.NewClient(stablenet.ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	measurementObid, err := req.GetCustomIntField("measurementObid")
	if err != nil {
		e := fmt.Errorf("could not extract measurementObid: %v", err)
		j.logger.Error(e.Error())
		return nil, e
	}
	startTime := time.Unix(0, req.TimeRange.FromEpochMs*int64(time.Millisecond))
	endTime := time.Unix(0, req.TimeRange.ToEpochMs*int64(time.Millisecond))
	metrics, err := snClient.FetchMetricsForMeasurement(measurementObid, startTime, endTime)
	if err != nil {
		e := fmt.Errorf("could not retrieve metrics from StableNet: %v", err)
		j.logger.Error(e.Error())
		return nil, e
	}
	return j.createResponseWithCustomData(metrics)
}

func (j *JsonDatasource) handleDataQuery(req *request.Content) (*datasource.DatasourceResponse, error) {
	snClient := stablenet.NewClient(stablenet.ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	measurementObid, err := req.GetCustomIntField("measurementObid")
	if err != nil {
		e := fmt.Errorf("could not extract measurementObid: %v", err)
		j.logger.Error(e.Error())
		return nil, e
	}
	metricName, err := req.GetCustomField("metricName")
	if err != nil {
		e := fmt.Errorf("could not extract metricName: %v", err)
		j.logger.Error(e.Error())
		return nil, e
	}
	startTime := time.Unix(0, req.TimeRange.FromEpochMs*int64(time.Millisecond))
	endTime := time.Unix(0, req.TimeRange.ToEpochMs*int64(time.Millisecond))
	data, err := snClient.FetchDataForMetric(measurementObid, metricName, startTime, endTime)
	if err != nil {
		e := fmt.Errorf("could not retrieve metrics from StableNet: %v", err)
		j.logger.Error(e.Error())
		return nil, e
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
		Points:               points,
	}
	result := datasource.QueryResult{
		RefId:    "A",
		Series:   []*datasource.TimeSeries{&timeSeries},
	}
	response := datasource.DatasourceResponse{
		Results: []*datasource.QueryResult{&result},
	}
	return &response, nil
}
