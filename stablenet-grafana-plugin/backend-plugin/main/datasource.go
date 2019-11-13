package main

import (
	"backend-plugin/request"
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

type JsonDatasource struct {
	plugin.NetRPCUnsupportedPlugin
	logger   hclog.Logger
	snClient stablenet.Client
}

func (j *JsonDatasource) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest) (*datasource.DatasourceResponse, error) {
	dsOptions := make(map[string]string)
	_ = json.Unmarshal([]byte(tsdbReq.Datasource.JsonData), &dsOptions)
	j.logger.Error(fmt.Sprintf("%v", tsdbReq.Datasource.DecryptedSecureJsonData))
	j.logger.Error(fmt.Sprintf("%v", dsOptions))
	port, portErr := strconv.Atoi(dsOptions["snport"])
	if portErr != nil {
		err := fmt.Errorf("could not parse port \"%s\"", dsOptions["snport"])
		j.logger.Error(portErr.Error())
		return &datasource.DatasourceResponse{
			Results: []*datasource.QueryResult{&datasource.QueryResult{
				Error: err.Error(),
			},
			},
		}, nil
	}
	j.snClient = stablenet.NewClient(stablenet.ConnectOptions{
		Host:     dsOptions["snip"],
		Port:     port,
		Username: dsOptions["snusername"],
		Password: tsdbReq.Datasource.DecryptedSecureJsonData["snpassword"],
	})
	startTime := time.Unix(0, tsdbReq.TimeRange.FromEpochMs*int64(time.Millisecond))
	endTime := time.Unix(0, tsdbReq.TimeRange.ToEpochMs*int64(time.Millisecond))
	handler := request.NewHandler(j.logger, j.snClient, startTime, endTime)
	results := make([]*datasource.QueryResult, 0, len(tsdbReq.Queries))
	for _, tsdbReq := range tsdbReq.Queries {
		query := request.Query{Query: *tsdbReq}
		result := handler.Handle(query)
		results = append(results, result)
	}
	response := &datasource.DatasourceResponse{
		Results: results,
	}
	return response, nil
}
