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

func GetHandlersForRequest(request Request) map[string]Handler {
	connectOptions, err := request.stableNetOptions()
	if err != nil {
		request.Logger.Error(fmt.Sprintf("could not extract StableNet(R) connect options: %v", err))
	}
	snClient := stablenet.NewClient(connectOptions)
	baseHandler := StableNetHandler{
		Logger:   request.Logger,
		SnClient: snClient,
	}
	handlers := make(map[string]Handler)
	handlers["devices"] = deviceHandler{StableNetHandler: &baseHandler}
	handlers["measurements"] = MeasurementHandler{StableNetHandler: &baseHandler}
	handlers["metricNames"] = MetricNameHandler{StableNetHandler: &baseHandler}
	handlers["testDatasource"] = datasourceTestHandler{StableNetHandler: &baseHandler}
	handlers["metricData"] = MetricDataHandler{StableNetHandler: &baseHandler}
	handlers["statisticLink"] = StatisticLinkHandler{StableNetHandler: &baseHandler}
	return handlers
}

type Request struct {
	*datasource.DatasourceRequest
	Logger hclog.Logger
}

func (r *Request) stableNetOptions() (*stablenet.ConnectOptions, error) {
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
	deviceQuery, err := q.GetCustomField("deviceQuery")
	if err != nil {
		return BuildErrorResult("could not extract the deviceQuery from the query", q.RefId), nil
	}
	devices, err := d.SnClient.QueryDevices(deviceQuery)
	if err != nil {
		e := fmt.Errorf("could not retrieve devices from StableNet(R): %v", err)
		d.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(devices, q.RefId)
}

type MeasurementHandler struct {
	*StableNetHandler
}

func (m MeasurementHandler) Process(query Query) (*datasource.QueryResult, error) {
	deviceObid, err := query.GetCustomIntField("deviceObid")
	if err != nil {
		return BuildErrorResult("could not extract deviceObid from the query", query.RefId), nil
	}
	measurements, err := m.SnClient.FetchMeasurementsForDevice(deviceObid)
	if err != nil {
		e := fmt.Errorf("could not fetch measurements from StableNet(R): %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(measurements, query.RefId)
}

type MetricNameHandler struct {
	*StableNetHandler
}

func (m MetricNameHandler) Process(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	metrics, err := m.SnClient.FetchMetricsForMeasurement(measurementObid)
	if err != nil {
		e := fmt.Errorf("could not retrieve metric names from StableNet(R): %v", err)
		m.Logger.Error(e.Error())
		return nil, e
	}
	return createResponseWithCustomData(metrics, query.RefId)
}

type MetricDataHandler struct {
	*StableNetHandler
}

func (m MetricDataHandler) Process(query Query) (*datasource.QueryResult, error) {
	measurementObid, err := query.GetCustomIntField("measurementObid")
	if err != nil {
		return BuildErrorResult("could not extract measurementObid from query", query.RefId), nil
	}
	metricId, err := query.GetCustomIntField("metricId")
	if err != nil {
		return BuildErrorResult("could not extract metricName from query", query.RefId), nil
	}

	series, err := m.fetchMetrics(query, measurementObid, []int{metricId})
	if err != nil {
		e := fmt.Errorf("could not fetch metric data from server: %v", err)
		m.Logger.Error(e.Error())
		return nil, e
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
	_, err := d.SnClient.FetchMeasurementsForDevice(-1)
	if err != nil {
		return BuildErrorResult("Cannot login into StableNet(R) with the provided credentials", query.RefId), nil
	}
	return &datasource.QueryResult{
		Series: []*datasource.TimeSeries{},
	}, nil
}

type StatisticLinkHandler struct {
	*StableNetHandler
}

func (s StatisticLinkHandler) Process(query Query) (*datasource.QueryResult, error) {
	link, err := query.GetCustomField("statisticLink")
	if err != nil {
		return BuildErrorResult("could not extract statisticLink parameter from query", query.RefId), nil
	}
	measurementRegex := regexp.MustCompile("[?&]id=(\\d+)")
	idMatches := measurementRegex.FindAllStringSubmatch(link, 1)
	if len(idMatches) == 0 {
		return BuildErrorResult(fmt.Sprintf("the link \"%s\" does not carry a measurement id.", link), query.RefId), nil
	}
	measurementId, _ := strconv.Atoi(idMatches[0][1])
	valueRegex := regexp.MustCompile("[?&]value\\d*=(\\d+)")
	valueMatches := valueRegex.FindAllStringSubmatch(link, -1)
	valueIds := make([]int, 0, len(valueMatches))
	for _, valueMatch := range valueMatches {
		id, _ := strconv.Atoi(valueMatch[1])
		valueIds = append(valueIds, id)
	}

	series, err := s.fetchMetrics(query, measurementId, valueIds)
	if err != nil {
		e := fmt.Errorf("could not fetch data for statistic link from server: %v", err)
		s.Logger.Error(e.Error())
		return nil, e
	}

	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: series,
	}
	return &result, nil
}

func createResponseWithCustomData(data interface{}, refId string) (*datasource.QueryResult, error) {
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
