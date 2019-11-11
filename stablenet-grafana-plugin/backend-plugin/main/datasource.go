package main

import (
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
	queryType, err := j.getQueryType(tsdbReq)
	if err != nil {
		j.logger.Error("could not retrieve query type: %v", err)
		return nil, err
	}
	if queryType == "devices" {
		return j.handleDeviceQuery(tsdbReq)
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

func (j *JsonDatasource) getQueryType(tsdbReq *datasource.DatasourceRequest) (string, error) {
	queryType := "query"
	if len(tsdbReq.Queries) > 0 {
		firstQuery := tsdbReq.Queries[0]
		j.logger.Info(firstQuery.ModelJson)
		queryJson, err := simplejson.NewJson([]byte(firstQuery.ModelJson))
		if err != nil {
			return "", err
		}
		queryType = queryJson.Get("queryType").MustString("devices")

	}
	return queryType, nil
}

func (j *JsonDatasource) handleDeviceQuery(tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	if len(tsdbReq.Queries) != 1 {
		err := fmt.Errorf("client queried devices but provided more than one query")
		j.logger.Error(err.Error())
		return nil, err
	}
	snClient := stablenet.NewClient(stablenet.ConnectOptions{Host: "127.0.0.1", Port: 5443, Username: "infosim", Password: "stablenet"})
	devices, err := snClient.FetchAllDevices()
	if err != nil {
		j.logger.Error("could not retreive devices from StableNet(R)", err)
		return nil, err
	}
	payload, err := json.Marshal(devices)
	if err != nil {
		j.logger.Error("could not parse json", err)
		return nil, err
	}
	result := datasource.QueryResult{
		RefId:    tsdbReq.Queries[0].RefId,
		MetaJson: string(payload),
		Series: []*datasource.TimeSeries{},
	}
	response := datasource.DatasourceResponse{
		Results: []*datasource.QueryResult{&result},
	}
	return &response, nil
}
