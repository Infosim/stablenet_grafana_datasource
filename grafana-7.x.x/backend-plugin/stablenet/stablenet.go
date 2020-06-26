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
	"encoding/xml"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type QueryDataProvider interface {
	FetchDataForMetrics(DataQueryOptions) (map[string]MetricDataSeries, error)
	FetchMetricsForMeasurement(int) ([]Metric, error)
}

type VersionProvider interface {
	QueryStableNetVersion() (*ServerVersion, *string)
}

type DeviceProvider interface {
	QueryDevices(string) (*DeviceQueryResult, error)
}

type MeasurementProvider interface {
	FetchMeasurementsForDevice(int, string) (*MeasurementQueryResult, error)
}

type MetricProvider interface {
	FetchMetricsForMeasurement(int) ([]Metric, error)
}

type ConnectOptions struct {
	Host     string `json:"snip"`
	Port     int    `json:"snport"`
	Username string `json:"snusername"`
	Password string `json:"snpassword"`
}

func NewClient(options *ConnectOptions) *Client {
	client := Client{ConnectOptions: *options, client: resty.New()}
	client.client.SetBasicAuth(options.Username, options.Password)
	client.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return &client
}

type Client struct {
	ConnectOptions
	client *resty.Client
}

// Queries StableNet® for its version. Attention: Unlike Go-conventions state,
// this function returns a string point instead of an error in case the version cannot be fetched.
// The reason is that the returned string is meant to be presented to the end user, while an error type string
// should generally not be presented to the end user.
func (c *Client) QueryStableNetVersion() (*ServerVersion, *string) {
	var errorStr string
	// use old XML API here because all server versions should have this endpoint, opposed to the JSON API version info endpoint.
	url := fmt.Sprintf("https://%s:%d/rest/info", c.Host, c.Port)
	resp, err := c.client.R().Get(url)
	if err != nil {
		errorStr = fmt.Sprintf("Connecting to StableNet® failed: %v", err.Error())
		return nil, &errorStr
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		errorStr = fmt.Sprintf("The StableNet® server could be reached, but the credentials were invalid.")
		return nil, &errorStr
	}
	if resp.StatusCode() != http.StatusOK {
		errorStr = fmt.Sprintf("Log in to StableNet® successful, but the StableNet® version could not be queried. Status Code: %d", resp.StatusCode())
		return nil, &errorStr
	}
	var result ServerInfo
	err = xml.Unmarshal(resp.Body(), &result)
	if err != nil {
		errorStr = fmt.Sprintf("Log in to StableNet® successful, but the StableNet® answer \"%s\" could not be parsed: %v", resp.String(), err)
		return nil, &errorStr
	}
	return &result.ServerVersion, nil
}

func (c *Client) buildStatusError(msg string, resp *resty.Response) error {
	return fmt.Errorf("%s: status code: %d, response: %s", msg, resp.StatusCode(), string(resp.Body()))
}

type DeviceQueryResult struct {
	Devices []Device `json:"data"`
	HasMore bool     `json:"hasMore"`
}

func (c *Client) QueryDevices(filter string) (*DeviceQueryResult, error) {
	var url string
	if len(filter) != 0 {
		filterParam := fmt.Sprintf("name ct '%s'", filter)
		url = c.buildJsonApiUrl("devices", "name", filterParam)
	} else {
		url = c.buildJsonApiUrl("devices", "name")
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

func (c *Client) buildJsonApiUrl(endpoint string, orderBy string, filters ...string) string {
	url := fmt.Sprintf("https://%s:%d/api/1/%s?$top=100", c.Host, c.Port, endpoint)
	if len(orderBy) != 0 {
		url = fmt.Sprintf("%s&$orderBy=%s", url, orderBy)
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

func (c *Client) FetchMeasurementsForDevice(deviceObid int, filter string) (*MeasurementQueryResult, error) {
	var nameFilter string
	if len(filter) != 0 {
		nameFilter = fmt.Sprintf("name ct '%s'", filter)
	}
	deviceFilter := fmt.Sprintf("destDeviceId eq '%d'", deviceObid)

	url := c.buildJsonApiUrl("measurements", "name", deviceFilter, nameFilter)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving measurements for device filter \"%s\" failed: %v", deviceFilter, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving measurements for device filter \"%s\" failed", deviceFilter), resp)
	}
	var result MeasurementQueryResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return &result, nil
}

func (c *Client) FetchMeasurementName(id int) (*string, error) {
	url := c.buildJsonApiUrl("measurements", "name", fmt.Sprintf("obid eq '%d'", id))
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving name for measurement %d failed: %v", id, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving name for measurement %d failed", id), resp)
	}
	var responseData MeasurementQueryResult
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	if len(responseData.Measurements) == 0 {
		return nil, fmt.Errorf("measurement with id %d does not exist", id)
	}
	return &responseData.Measurements[0].Name, nil
}

func (c *Client) FetchMetricsForMeasurement(measurementObid int) ([]Metric, error) {
	endpoint := fmt.Sprintf("measurements/%d/metrics", measurementObid)
	//orderby is empty because it' currently not supported by the endpoint
	url := c.buildJsonApiUrl(endpoint, "")
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

func (c *Client) FetchDataForMetrics(options DataQueryOptions) (map[string]MetricDataSeries, error) {
	startMillis := options.Start.UnixNano() / int64(time.Millisecond)
	endMillis := options.End.UnixNano() / int64(time.Millisecond)
	query := DataQuery{Start: startMillis, End: endMillis, Metrics: options.Metrics, Raw: false, Average: options.Average}
	endpoint := fmt.Sprintf("measurements/%d/data", options.MeasurementObid)
	url := c.buildJsonApiUrl(endpoint, "")
	resp, err := c.client.R().SetHeader("Content-Type", "application/json").SetBody(query).Post(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving metric data for measurement %d failed: %v", options.MeasurementObid, err)
	}
	if resp.StatusCode() != 200 {
		return nil, c.buildStatusError(fmt.Sprintf("retrieving metric data for measurement %d failed", options.MeasurementObid), resp)
	}
	return parseStatisticByteSlice(resp.Body(), options.Metrics)
}

func parseStatisticByteSlice(bytes []byte, metricKeys []string) (map[string]MetricDataSeries, error) {
	var data []timestampResponse
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	resultMap := make(map[string]MetricDataSeries)
	for _, record := range data {
		converted := parseSingleTimestamp(record, metricKeys)
		for key, measurementData := range converted {
			if _, ok := resultMap[key]; !ok {
				resultMap[key] = make([]MetricData, 0, 0)
			}
			resultMap[key] = append(resultMap[key], measurementData)
		}
	}
	return resultMap, nil
}

func (c *Client) formatMetricIds(valueIds []int) string {
	if len(valueIds) == 1 {
		return fmt.Sprintf("value=%d", valueIds[0])
	}
	query := make([]string, 0, len(valueIds))
	for index, valueId := range valueIds {
		query = append(query, fmt.Sprintf("value%d=%d", index, valueId))
	}
	return strings.Join(query, "&")
}
