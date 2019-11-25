package query

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/hashicorp/go-hclog"
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
	handlers["devices"] = DeviceHandler{StableNetHandler: &baseHandler}
	handlers["measurements"] = MeasurementHandler{StableNetHandler: &baseHandler}
	handlers["metricNames"] = MetricNameHandler{StableNetHandler: &baseHandler}
	handlers["testDatasource"] = DatasourceTestHandler{StableNetHandler: &baseHandler}
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
