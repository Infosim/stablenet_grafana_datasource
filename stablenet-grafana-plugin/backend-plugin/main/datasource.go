package main

import (
	"backend-plugin/query"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

type BackendPlugin interface {
	Query(context.Context, *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error)
}

func NewBackendPlugin(logger hclog.Logger) BackendPlugin {
	plugin := jsonDatasource{logger: logger}
	plugin.handlerFactory = query.GetHandlersForRequest
	return &plugin
}

type jsonDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	logger         hclog.Logger
	handlerFactory func(query.Request) (map[string]query.Handler, error)
}

func (j *jsonDatasource) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	request := query.Request{DatasourceRequest: tsdbReq, Logger: j.logger}
	handlers, err := j.handlerFactory(request)
	if err != nil {
		return &datasource.DatasourceResponse{Results: []*datasource.QueryResult{query.BuildErrorResult(err.Error(), "")}}, nil
	}
	results := make([]*datasource.QueryResult, 0, len(tsdbReq.Queries))
	for _, tsdbReq := range tsdbReq.Queries {
		startTime, endTime := request.ToTimeRange()
		q := query.Query{Query: *tsdbReq, StartTime: startTime, EndTime: endTime}
		var result *datasource.QueryResult
		queryType, queryTypeError := q.GetCustomField("queryType")
		if queryTypeError != nil {
			msg := fmt.Sprintf("could not retrieve q type: %v", queryTypeError)
			result = query.BuildErrorResult(msg, q.RefId)
			continue
		} else if _, ok := handlers[queryType]; !ok {
			msg := fmt.Sprintf("queryType \"%s\" is unknown", queryType)
			result = query.BuildErrorResult(msg, q.RefId)
			continue
		}
		handler := handlers[queryType]
		result, err := handler.Process(q)
		if err != nil {
			result = query.BuildErrorResult("Internal Plugin Error. Please consult the Grafana log files.", q.RefId)
		}
		results = append(results, result)
	}
	response := &datasource.DatasourceResponse{
		Results: results,
	}
	return response, nil
}
