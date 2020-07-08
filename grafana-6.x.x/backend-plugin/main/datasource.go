/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package main

import (
	"backend-plugin/query"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

type BackendPlugin interface {
	Query(context.Context, *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error)
}

func NewBackendPlugin(logger hclog.Logger) BackendPlugin {
	snPlugin := jsonDatasource{logger: logger, validationStore: make(map[int64]bool)}
	snPlugin.handlerFactory = query.GetHandlersForRequest
	return &snPlugin
}

type jsonDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	logger          hclog.Logger
	validationStore map[int64]bool
	handlerFactory  func(query.Request, map[int64]bool, int64) (map[string]query.Handler, error)
}

func (j *jsonDatasource) Query(ctx context.Context, datasourceRequest *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	request := query.Request{DatasourceRequest: datasourceRequest, Logger: j.logger}
	handlers, err := j.handlerFactory(request, j.validationStore, datasourceRequest.Datasource.Id)
	if err != nil {
		return &datasource.DatasourceResponse{Results: []*datasource.QueryResult{query.BuildErrorResult(err.Error(), "")}}, nil
	}
	results := make([]*datasource.QueryResult, 0, len(datasourceRequest.Queries))
	for _, datasourceQuery := range datasourceRequest.Queries {
		startTime, endTime := request.ToTimeRange()
		q := query.Query{Query: *datasourceQuery, StartTime: startTime, EndTime: endTime}
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
