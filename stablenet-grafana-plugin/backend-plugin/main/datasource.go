package main

import (
	"backend-plugin/query"
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
	request := request{tsdbReq}
	connectOptions, err := request.stableNetOptions()
	if err != nil {
		j.logger.Error(fmt.Sprintf("could not extract StableNet(R) connect options: %v", err))
	}
	j.snClient = stablenet.NewClient(connectOptions)
	startTime, endTime := request.timeRange()
	handler := query.NewHandler(j.logger, j.snClient, startTime, endTime)
	results := make([]*datasource.QueryResult, 0, len(tsdbReq.Queries))
	for _, tsdbReq := range tsdbReq.Queries {
		query := query.Query{Query: *tsdbReq}
		result := handler.Handle(query)
		results = append(results, result)
	}
	response := &datasource.DatasourceResponse{
		Results: results,
	}
	return response, nil
}

type request struct {
	*datasource.DatasourceRequest
}

func (r *request) stableNetOptions() (*stablenet.ConnectOptions, error) {
	info := r.Datasource
	options := make(map[string]string)
	err := json.Unmarshal([]byte(info.JsonData), &options)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal jsonData of the datasource: %v", err)
	}
	if _, ok := options["snip"]; !ok {
		return nil, fmt.Errorf("the snip is missing in the jsonData of the datasource")
	}
	if _, ok := options["snport"]; !ok {
		return nil, fmt.Errorf("the snport is missing in the jsonData of the datasource")
	}
	if _, ok := options["snusername"]; !ok {
		return nil, fmt.Errorf("the snusername is missing in the jsonData of the datasource")
	}
	if _, ok := info.DecryptedSecureJsonData["snpassword"]; !ok {
		return nil, fmt.Errorf("the snpassword is missing in the encryptedJsonData of the datasource")
	}
	port, portErr := strconv.Atoi(options["snport"])
	if portErr != nil {
		return nil, fmt.Errorf("could not parse snport into number: %v", portErr)
	}
	return &stablenet.ConnectOptions{
		Host:     options["snip"],
		Port:     port,
		Username: options["snusername"],
		Password: info.DecryptedSecureJsonData["snpassword"],
	}, nil
}

func (r *request) timeRange() (startTime time.Time, endTime time.Time) {
	startTime = time.Unix(0, r.TimeRange.FromEpochMs*int64(time.Millisecond))
	endTime = time.Unix(0, r.TimeRange.ToEpochMs*int64(time.Millisecond))
	return
}
