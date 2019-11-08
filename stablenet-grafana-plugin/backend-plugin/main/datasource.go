package main

import (
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

func (ds *JsonDatasource) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	ds.logger.Debug("Query", "datasource", tsdbReq.Datasource.Name, "TimeRange", tsdbReq.TimeRange)
	then := time.Now().AddDate(-1,0,0)
	points := make([]*datasource.Point, 0, 0)
	for i:=0; i < 10; i++{
		point := datasource.Point{
			Timestamp: then.UnixNano()/int64(time.Millisecond),
			Value:     float64(i*1000),
		}
		points = append(points, &point)
		then = then.Add(-time.Hour)
	}
	timeSeries := datasource.TimeSeries{
		Name:                 "Test Series",
		Tags:                 nil,
		Points:               points,
	}
	queryResult := datasource.QueryResult{
		Error:                "",
		RefId:                "A",
		MetaJson:             "",
		Series:               []*datasource.TimeSeries {&timeSeries},
		Tables:               nil,
	}
	response := &datasource.DatasourceResponse{
		Results:              []*datasource.QueryResult{&queryResult},
	}
	return response, nil
}

func (ds *JsonDatasource) MetricQuery(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	return ds.Query(ctx, tsdbReq)
}
