package stablenet

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	url2 "net/url"
	"strings"
	"time"
)

type Client interface {
	QueryDevices(string) ([]Device, error)
	FetchMeasurementsForDevice(int) ([]Measurement, error)
	FetchMetricsForMeasurement(int) ([]Metric, error)
	FetchDataForMetrics(int, []int, time.Time, time.Time) (map[string]MetricDataSeries, error)
}

type ConnectOptions struct {
	Host     string `json:"snip"`
	Port     int    `json:"snport"`
	Username string `json:"snusername"`
	Password string `json:"snpassword"`
}

func NewClient(options *ConnectOptions) Client {
	client := ClientImpl{ConnectOptions: *options, client: resty.New()}
	client.client.SetBasicAuth(options.Username, options.Password)
	client.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return &client
}

type ClientImpl struct {
	ConnectOptions
	client *resty.Client
}

func (c *ClientImpl) QueryDevices(deviceQuery string) ([]Device, error) {
	filter := fmt.Sprintf("name ct '%s'", deviceQuery)
	url := fmt.Sprintf("https://%s:%d/api/1/devices?$filter=%s", c.Host, c.Port, url2.QueryEscape(filter))
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("the statuscode was \"%d\" and the message was \"%s\"", resp.StatusCode(), resp.Status())
	}
	type serverResponse struct {
		Devices []Device `json:"data"`
	}
	var collection serverResponse
	err = json.Unmarshal(resp.Body(), &collection)
	return collection.Devices, err
}

func (c *ClientImpl) FetchMeasurementsForDevice(deviceObid int) ([]Measurement, error) {
	filter := fmt.Sprintf("destDeviceId eq '%d'", deviceObid)
	url := fmt.Sprintf("https://%s:%d/api/1/measurements?$filter=%s", c.Host, c.Port, url2.QueryEscape(filter))
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve measurements for device %d from StableNet: %v", deviceObid, err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("could not retrieve measurements for device %d, the error code was %d with message \"%s\"; %s", deviceObid, resp.StatusCode(), resp.Status(), url)
	}
	collection := struct {
		Measurements []Measurement `json:"data"`
	}{}
	err = json.Unmarshal(resp.Body(), &collection)
	return collection.Measurements, err
}

func (c *ClientImpl) FetchMetricsForMeasurement(measurementObid int) ([]Metric, error) {
	url := fmt.Sprintf("https://%s:%d/api/1/measurements/%d/metrics", c.Host, c.Port, measurementObid)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics for measurement %d from StableNet(R)", measurementObid)
	}
	responseData := struct {
		ValueOutputs []Metric `json:"valueOutputs"`
	}{}
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return responseData.ValueOutputs, nil
}

func (c *ClientImpl) FetchDataForMetrics(measurementObid int, metricIds []int, startTime time.Time, endTime time.Time) (map[string]MetricDataSeries, error) {
	startMillis := startTime.UnixNano() / int64(time.Millisecond)
	endMillis := endTime.UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("https://%s:%d/StatisticServlet?stat=1010&type=json&login=%s,%s&id=%d&start=%d&end=%d&%s", c.Host, c.Port, c.Username, c.Password, measurementObid, startMillis, endMillis, c.formatMetricIds(metricIds))
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics for measurement %d from StableNet: %v", measurementObid, err)
	}
	data := make([]map[string]string, 0, 0)
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	resultMap := make(map[string]MetricDataSeries)
	for _, record := range data {
		converted, err := parseSingleTimestamp(record)
		if err != nil {
			return nil, fmt.Errorf("parsing an entry from RawStatisticServlet failed: %v", err)
		}
		for key, measurementData := range converted {
			if _, ok := resultMap[key]; !ok {
				resultMap[key] = make([]MetricData, 0, 0)
			}
			resultMap[key] = append(resultMap[key], measurementData)
		}
	}
	return resultMap, nil
}

func (c *ClientImpl) formatMetricIds(valueIds []int) string {
	query := make([]string, 0, len(valueIds))
	for _, valueId := range valueIds {
		query = append(query, fmt.Sprintf("value=%d", valueId))
	}
	return strings.Join(query, "&")
}
