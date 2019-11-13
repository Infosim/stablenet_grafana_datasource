package main

import (
	"backend-plugin/request"
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

type JsonDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
	snClient stablenet.Client
}

func (j *JsonDatasource) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	req := &request.Content{*tsdbReq}
	port, portErr := strconv.Atoi(req.Datasource.DecryptedSecureJsonData["snport"])
	if portErr != nil{
		err := fmt.Errorf("could not parse port: %v", portErr)
		j.logger.Error(portErr.Error())
		return nil, err
	}
	j.snClient = stablenet.NewClient(stablenet.ConnectOptions{
		Host:     req.Datasource.DecryptedSecureJsonData["snip"],
		Port:     port,
		Username: req.Datasource.DecryptedSecureJsonData["snusername"],
		Password: req.Datasource.DecryptedSecureJsonData["snpassword"],
	})
	j.logger.Error("Hello World")
	j.logger.Error(fmt.Sprintf("%v", tsdbReq.Datasource.JsonData))
	j.logger.Error(fmt.Sprintf("%v", tsdbReq.Datasource.DecryptedSecureJsonData))
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
	err = fmt.Errorf("queryType \"%s\" not known", queryType)
	j.logger.Error(err.Error())
	return nil, err
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
	devices, err := j.snClient.FetchAllDevices()
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
	deviceObid, err := req.GetCustomIntField("deviceObid")
	if err != nil {
		j.logger.Error(err.Error())
		return nil, err
	}
	measurements, err := j.snClient.FetchMeasurementsForDevice(deviceObid)
	if err != nil {
		j.logger.Error(err.Error())
		return nil, err
	}
	return j.createResponseWithCustomData(measurements)
}

func (j *JsonDatasource) handleMetricNameQuery(req *request.Content) (*datasource.DatasourceResponse, error) {
	measurementObid, err := req.GetCustomIntField("measurementObid")
	if err != nil {
		e := fmt.Errorf("could not extract measurementObid: %v", err)
		j.logger.Error(e.Error())
		return nil, e
	}
	startTime := time.Unix(0, req.TimeRange.FromEpochMs*int64(time.Millisecond))
	endTime := time.Unix(0, req.TimeRange.ToEpochMs*int64(time.Millisecond))
	metrics, err := j.snClient.FetchMetricsForMeasurement(measurementObid, startTime, endTime)
	if err != nil {
		e := fmt.Errorf("could not retrieve metrics from StableNet: %v", err)
		j.logger.Error(e.Error())
		return nil, e
	}
	return j.createResponseWithCustomData(metrics)
}

func (j *JsonDatasource) handleDataQuery(req *request.Content) (*datasource.DatasourceResponse, error) {
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
	data, err := j.snClient.FetchDataForMetric(measurementObid, metricName, startTime, endTime)
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
