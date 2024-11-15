/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
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
	"net/http"
	url2 "net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type ConnectOptions struct {
	Address  string `json:"snip"`
	Username string `json:"snusername"`
	Password string `json:"snpassword"`
}

func NewStableNetClient(options *ConnectOptions) *StableNetClient {
	client := resty.New().
		SetBasicAuth(options.Username, options.Password).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return &StableNetClient{ConnectOptions: *options, client: client}
}

type StableNetClient struct {
	ConnectOptions
	client *resty.Client
}

func (stableNetClient *StableNetClient) get(path string) (*resty.Response, error) {
	url := stableNetClient.Address + path
	return stableNetClient.client.R().Get(url)
}

// Queries StableNet® for its version. Attention: Unlike Go-conventions state, this function returns a string point instead of an error in case the version cannot be fetched.
// The reason is that the returned string is meant to be presented to the end user, while an error type string should generally not be presented to the end user.
func (stableNetClient *StableNetClient) QueryStableNetInfo() (*ServerInfo, *string) {
	// use old XML API here because all server versions should have this endpoint, opposed to the JSON API version info endpoint.
	response, err := stableNetClient.get("/rest/info")

	if err != nil {
		errorStr := fmt.Sprintf("Connecting to StableNet® failed: %v", err.Error())
		return nil, &errorStr
	}

	if response.StatusCode() == http.StatusUnauthorized {
		errorStr := "The StableNet® server could be reached, but the credentials were invalid."
		return nil, &errorStr
	}

	if response.StatusCode() != http.StatusOK {
		errorStr := fmt.Sprintf("Log in to StableNet® successful, but the StableNet® version could not be queried. Status Code: %d", response.StatusCode())
		return nil, &errorStr
	}

	var result ServerInfo
	err = xml.Unmarshal(response.Body(), &result)
	if err != nil {
		errorStr := fmt.Sprintf("Log in to StableNet® successful, but the StableNet® answer \"%s\" could not be parsed: %v", response.String(), err)
		return nil, &errorStr
	}

	return &result, nil
}

func buildStatusError(msg string, resp *resty.Response) error {
	return fmt.Errorf("%s: status code: %d, response: %s", msg, resp.StatusCode(), string(resp.Body()))
}

func buildJsonApiUrl(endpoint string, orderBy string, filters ...string) string {
	url := fmt.Sprintf("/api/1/%s?$top=100", endpoint)

	if len(orderBy) != 0 {
		url = url + fmt.Sprintf("&$orderBy=%s", orderBy)
	}

	nonEmptyFilters := make([]string, 0, len(filters))
	for _, f := range filters {
		if len(f) > 0 {
			nonEmptyFilters = append(nonEmptyFilters, f)
		}
	}

	if len(nonEmptyFilters) == 0 {
		return url
	}

	return url + "&$filter=" + url2.QueryEscape(strings.Join(nonEmptyFilters, " and "))
}

// Queries devices from the StableNet server that contain the string "nameFilter" in their nae
func (stableNetClient *StableNetClient) QueryDevices(nameFilter string) (*DeviceQueryResult, error) {
	path := buildJsonApiUrl("devices", "name")

	if len(nameFilter) != 0 {
		path = path + "&$filter=" + url2.QueryEscape(fmt.Sprintf("name ct '%s'", nameFilter))
	}

	resp, err := stableNetClient.get(path)
	if err != nil {
		return nil, fmt.Errorf("retrieving devices matching query \"%s\" failed: %v", nameFilter, err)
	}
	if resp.StatusCode() != 200 {
		return nil, buildStatusError(fmt.Sprintf("retrieving devices matching query \"%s\" failed", nameFilter), resp)
	}

	var result DeviceQueryResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return &result, nil
}

func (stableNetClient *StableNetClient) FetchMeasurementsForDevice(deviceObid int, fieldFilter string) (*MeasurementQueryResult, error) {
	var nameFilter string
	if len(fieldFilter) != 0 {
		nameFilter = fmt.Sprintf("name ct '%s'", fieldFilter)
	}

	deviceFilter := fmt.Sprintf("destDeviceId eq '%d'", deviceObid)

	path := buildJsonApiUrl("measurements", "name", deviceFilter, nameFilter)

	resp, err := stableNetClient.get(path)
	if err != nil {
		return nil, fmt.Errorf("retrieving measurements for device filter \"%s\" failed: %v", deviceFilter, err)
	}
	if resp.StatusCode() != 200 {
		return nil, buildStatusError(fmt.Sprintf("retrieving measurements for device filter \"%s\" failed", deviceFilter), resp)
	}

	var result MeasurementQueryResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return &result, nil
}

func (stableNetCliet *StableNetClient) FetchMeasurementName(id int) (*string, error) {
	url := buildJsonApiUrl("measurements", "name", fmt.Sprintf("obid eq '%d'", id))

	resp, err := stableNetCliet.get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving name for measurement %d failed: %v", id, err)
	}
	if resp.StatusCode() != 200 {
		return nil, buildStatusError(fmt.Sprintf("retrieving name for measurement %d failed", id), resp)
	}

	var responseData MeasurementQueryResult
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}

	if len(responseData.Data) == 0 {
		return nil, fmt.Errorf("measurement with id %d does not exist", id)
	}

	return &responseData.Data[0].Name, nil
}

func (stableNetClient *StableNetClient) FetchMetricsForMeasurement(measurementObid int) ([]Metric, error) {
	url := fmt.Sprintf("/api/1/measurement-data/%d/metrics?$top=100", measurementObid)

	resp, err := stableNetClient.get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving metrics for measurement %d failed: %v", measurementObid, err)
	}
	if resp.StatusCode() != 200 {
		return nil, buildStatusError(fmt.Sprintf("retrieving metrics for measurement %d failed", measurementObid), resp)
	}

	responseData := make([]Metric, 0)
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return responseData, nil
}

func (stableNetClient *StableNetClient) FetchDataForMetrics(options DataQueryOptions) (map[string]MetricDataSeries, error) {
	query := DataQuery{
		Start:   options.Start.UnixNano() / int64(time.Millisecond),
		End:     options.End.UnixNano() / int64(time.Millisecond),
		Metrics: options.Metrics,
		Average: options.Average,
		Raw:     false,
	}

	url := stableNetClient.Address + fmt.Sprintf("/api/1/measurement-data/%d?$top=100", options.MeasurementObid)

	resp, err := stableNetClient.client.R().SetHeader("Content-Type", "application/json").SetBody(query).Post(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving metric data for measurement %d failed: %v", options.MeasurementObid, err)
	}
	if resp.StatusCode() != 200 {
		return nil, buildStatusError(fmt.Sprintf("retrieving metric data for measurement %d failed", options.MeasurementObid), resp)
	}

	return parseStatisticByteSlice(resp.Body())
}

type MeasurementDataEntryDTO struct {
	Timestamp       int64    `json:"timestamp"`
	Interval        int64    `json:"interval"`
	MissingInterval int64    `json:"missingInterval"`
	Min             *float64 `json:"min"`
	Max             *float64 `json:"max"`
	Avg             *float64 `json:"avg"`
}

type MeasurementMetricResultDataDTO struct {
	MetricName   string                    `json:"metricName"`
	MetricKey    string                    `json:"metricKey"`
	MetricDataId int                       `json:"metricDataId"`
	Data         []MeasurementDataEntryDTO `json:"data"`
}

type MeasurementMultiMetricResultDataDTO struct {
	MeasruementId int                              `json:"measurementId"`
	Values        []MeasurementMetricResultDataDTO `json:"values"`
}

func convertMeasurementData(data MeasurementDataEntryDTO) MetricData {
	return MetricData{
		Interval: time.Duration(data.Interval) * time.Millisecond,
		Time:     time.Unix(0, data.Timestamp*int64(time.Millisecond)),
		Min:      *data.Min,
		Avg:      *data.Avg,
		Max:      *data.Avg,
	}
}

func parseStatisticByteSlice(bytes []byte) (map[string]MetricDataSeries, error) {
	var data MeasurementMultiMetricResultDataDTO
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}

	resultMap := make(map[string]MetricDataSeries)
	for _, record := range data.Values {
		key := record.MetricKey
		data := record.Data

		for _, measurementData := range data {
			metricData := convertMeasurementData(measurementData)
			if _, ok := resultMap[key]; !ok {
				resultMap[key] = make([]MetricData, 0)
			}
			resultMap[key] = append(resultMap[key], metricData)
		}
	}

	return resultMap, nil
}
