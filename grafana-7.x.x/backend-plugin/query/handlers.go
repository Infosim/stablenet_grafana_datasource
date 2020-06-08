/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package query

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
	"regexp"
	"strconv"
	"time"
)

func GetHandlersForRequest(request Request) (map[string]Handler, error) {
	connectOptions, err := request.stableNetOptions()
	if err != nil {
		return nil, fmt.Errorf("could not extract StableNet(R) connect options: %v", err)
	}
	snClient := stablenet.NewClient(connectOptions)
	baseHandler := StableNetHandler{
		Logger:   request.Logger,
		SnClient: snClient,
	}
	handlers := make(map[string]Handler)
	handlers["devices"] = deviceHandler{StableNetHandler: &baseHandler}
	handlers["measurements"] = measurementHandler{StableNetHandler: &baseHandler}
	handlers["metricNames"] = metricNameHandler{StableNetHandler: &baseHandler}
	handlers["testDatasource"] = datasourceTestHandler{StableNetHandler: &baseHandler}
	handlers["metricData"] = metricDataHandler{StableNetHandler: &baseHandler}
	handlers["statisticLink"] = statisticLinkHandler{StableNetHandler: &baseHandler}
	return handlers, nil
}

type Request struct {
	*datasource.DatasourceRequest
	Logger hclog.Logger
}

func (r *Request) stableNetOptions() (*stablenet.ConnectOptions, error) {
	info := r.Datasource
	options := make(map[string]string)
	if info == nil {
		return nil, fmt.Errorf("datasource info is nil")
	}
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

func (r *Request) ToTimeRange() (startTime time.Time, endTime time.Time) {
	startTime = time.Unix(0, r.TimeRange.FromEpochMs*int64(time.Millisecond))
	endTime = time.Unix(0, r.TimeRange.ToEpochMs*int64(time.Millisecond))
	return
}

type Handler interface {
	Process(Query) (*datasource.QueryResult, error)
}

type deviceHandler struct {
	*StableNetHandler
}

func (d deviceHandler) Process(q Query) (*datasource.QueryResult, error) {
	filter, _ := q.GetCustomField("filter")
	queryResult, err := d.SnClient.QueryDevices(filter)
	if err != nil {
		e := fmt.Errorf("could not retrieve devices from StableNet(R): %v", err)
		d.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(queryResult, q.RefId), nil
}

type measurementHandler struct {
	*StableNetHandler
}

func (m measurementHandler) Process(query Query) (*datasource.QueryResult, error) {
	deviceObid, _ := query.GetCustomIntField("deviceObid")
	measurementFilter, _ := query.GetCustomField("filter")
	measurements, err := m.SnClient.FetchMeasurementsForDevice(deviceObid, measurementFilter)
	if err != nil {
		e := fmt.Errorf("could not fetch measurements from StableNet(R): %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(measurements, query.RefId), nil
}

type metricNameHandler struct {
	*StableNetHandler
}

func (m metricNameHandler) Process(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	filter, _ := query.GetCustomField("filter")
	metrics, err := m.SnClient.FetchMetricsForMeasurement(*measurementObid, filter)
	if err != nil {
		e := fmt.Errorf("could not retrieve metric names from StableNet(R): %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(metrics, query.RefId), nil
}

type metricDataHandler struct {
	*StableNetHandler
}

func (m metricDataHandler) Process(query Query) (*datasource.QueryResult, error) {
	requests, err := query.GetMeasurementDataRequest()
	if err != nil {
		return BuildErrorResult(fmt.Sprintf("could not extract measurement requests from query: %v", err), query.RefId), nil
	}

	series := make([]*datasource.TimeSeries, 0, 0)
	for _, request := range requests {
		requestSeries, err := m.fetchMetrics(query, request.MeasurementObid, request.Metrics)
		if err != nil {
			e := fmt.Errorf("could not fetch metric data from server: %v", err)
			m.Logger.Error(e.Error())
			return nil, e
		}
		series = append(series, requestSeries...)
	}

	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: series,
	}
	return &result, nil
}

type datasourceTestHandler struct {
	*StableNetHandler
}

func (d datasourceTestHandler) Process(query Query) (*datasource.QueryResult, error) {
	version, errStr := d.SnClient.QueryStableNetVersion()
	if errStr != nil {
		return BuildErrorResult(*errStr, query.RefId), nil
	}
	versionRegex := regexp.MustCompile("^(?:9|[1-9]\\d)\\.")
	if !versionRegex.MatchString(version.Version) {
		return BuildErrorResult(fmt.Sprintf("The StableNet® version %s does not support Grafana®.", version.Version), query.RefId), nil
	}
	return &datasource.QueryResult{
		Series: []*datasource.TimeSeries{},
		RefId:  query.RefId,
	}, nil
}

func createResponseWithCustomData(data interface{}, refId string) *datasource.QueryResult {
	payload, err := json.Marshal(data)
	if err != nil {
		//json.Marshal returns a non-nil error if the data contains an invalid type such as channels or math.Inf(1)
		//since these types are programming errors, the program should panic in that case.
		panic(fmt.Sprintf("marshalling failed: %v", err))
	}
	result := datasource.QueryResult{
		RefId:    refId,
		MetaJson: string(payload),
		Series:   []*datasource.TimeSeries{},
	}
	return &result
}
