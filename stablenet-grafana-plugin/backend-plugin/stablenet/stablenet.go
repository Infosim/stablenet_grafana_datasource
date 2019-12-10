/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
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
	QueryDevices(string) (*DeviceQueryResult, error)
	FetchMeasurementsForDevice(*int, string) (*MeasurementQueryResult, error)
	FetchMetricsForMeasurement(int, string) ([]Metric, error)
	FetchDataForMetrics(int, []string, time.Time, time.Time) (map[string]MetricDataSeries, error)
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

func (c *ClientImpl) buildStatusError(msg string, resp *resty.Response) error {
	return fmt.Errorf("%s: status code: %d, response: %s", msg, resp.StatusCode(), string(resp.Body()))
}

type DeviceQueryResult struct {
	Devices []Device `json:"data"`
	HasMore bool     `json:"hasMore"`
}

func (c *ClientImpl) QueryDevices(filter string) (*DeviceQueryResult, error) {
	var url string
	if len(filter) != 0 {
		filterParam := fmt.Sprintf("name ct '%s'", filter)
		url = c.buildJsonApiUrl("devices", filterParam)
	} else {
		url = c.buildJsonApiUrl("devices")
	}
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving devices matching query \"%s\" failed: %v", filter, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving devices matching query \"%s\" failed", filter), resp)
	}
	var result DeviceQueryResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return &result, nil
}

func (c *ClientImpl) buildJsonApiUrl(endpoint string, filters ...string) string {
	url := fmt.Sprintf("https://%s:%d/api/1/%s?$top=100", c.Host, c.Port, endpoint)
	nonEmpty := make([]string, 0, len(filters))
	for _, f := range filters {
		if len(f) > 0 {
			nonEmpty = append(nonEmpty, f)
		}
	}
	if len(nonEmpty) == 0 {
		return url
	}
	filter := "&$filter=" + url2.QueryEscape(strings.Join(nonEmpty, " and "))
	return url + filter
}

func (c *ClientImpl) buildJsonApiUrlWithLimit(endpoint string, limit bool, filters ...string) string {
	url := fmt.Sprintf("https://%s:%d/api/1/%s?$top=100", c.Host, c.Port, endpoint)
	if !limit {
		url = fmt.Sprintf("https://%s:%d/api/1/%s?top=-1", c.Host, c.Port, endpoint)
	}
	nonEmpty := make([]string, 0, len(filters))
	for _, f := range filters {
		if len(f) > 0 {
			nonEmpty = append(nonEmpty, f)
		}
	}
	if len(nonEmpty) == 0 {
		return url
	}
	filter := "&$filter=" + url2.QueryEscape(strings.Join(nonEmpty, " and "))
	return url + filter
}

type MeasurementQueryResult struct {
	Measurements []Measurement `json:"data"`
	HasMore      bool          `json:"hasMore"`
}

func (c *ClientImpl) FetchMeasurementsForDevice(deviceObid *int, filter string) (*MeasurementQueryResult, error) {
	var deviceFilter, nameFilter string
	if deviceObid != nil {
		deviceFilter = fmt.Sprintf("destDeviceId eq '%d'", *deviceObid)
	}
	if len(filter) != 0 {
		nameFilter = fmt.Sprintf("name ct '%s'", filter)
	}
	url := c.buildJsonApiUrl("measurements", deviceFilter, nameFilter)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving measurements for device filter \"%s\" and name filter \"%s\" failed: %v", deviceFilter, nameFilter, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving measurements for device filter \"%s\" and name filter \"%s\" failed", deviceFilter, nameFilter), resp)
	}
	var result MeasurementQueryResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return &result, nil
}

func (c *ClientImpl) FetchMetricsForMeasurement(measurementObid int, filter string) ([]Metric, error) {
	var nameFilter string
	if len(filter) != 0 {
		nameFilter = fmt.Sprintf("name ct '%s'", filter)
	}
	endpoint := fmt.Sprintf("measurements/%d/metrics", measurementObid)
	url := c.buildJsonApiUrl(endpoint, nameFilter)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving metrics for measurement %d failed: %v", measurementObid, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving metrics for measurement %d failed", measurementObid), resp)
	}
	responseData := make([]Metric, 0, 0)
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return responseData, nil
}

func (c *ClientImpl) FetchDataForMetrics(measurementObid int, metricKeys []string, startTime time.Time, endTime time.Time) (map[string]MetricDataSeries, error) {
	startMillis := startTime.UnixNano() / int64(time.Millisecond)
	endMillis := endTime.UnixNano() / int64(time.Millisecond)
	query := struct {
		Start   int64    `json:"start"`
		End     int64    `json:"end"`
		Metrics []string `json:"metrics"`
		Raw     bool     `json:"raw"`
	}{
		Start: startMillis, End: endMillis, Metrics: metricKeys, Raw: true,
	}
	endpoint := fmt.Sprintf("measurements/%d/data", measurementObid)
	url := c.buildJsonApiUrlWithLimit(endpoint, false)
	resp, err := c.client.R().SetHeader("Content-Type", "application/json").SetBody(query).Post(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving metric data for measurement %d failed: %v", measurementObid, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving metric data for measurement %d failed", measurementObid), resp)
	}
	return parseStatisticByteSlice(resp.Body(), metricKeys)
}

func parseStatisticByteSlice(bytes []byte, metricKeys []string) (map[string]MetricDataSeries, error) {
	var data []timestampResponse
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	resultMap := make(map[string]MetricDataSeries)
	for _, record := range data {
		converted, err := parseSingleTimestamp(record, metricKeys)
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
	if len(valueIds) == 1 {
		return fmt.Sprintf("value=%d", valueIds[0])
	}
	query := make([]string, 0, len(valueIds))
	for index, valueId := range valueIds {
		query = append(query, fmt.Sprintf("value%d=%d", index, valueId))
	}
	return strings.Join(query, "&")
}
