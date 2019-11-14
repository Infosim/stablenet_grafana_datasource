package stablenet

import (
	bytes2 "bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	url2 "net/url"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	QueryDevices(deviceQuey string) ([]Device, error)
	FetchMeasurementsForDevice(int) ([]Measurement, error)
	FetchMetricsForMeasurement(int, time.Time, time.Time) ([]string, error)
	FetchDataForMetric(int, string, time.Time, time.Time) ([]MetricData, error)
}

type ConnectOptions struct {
	Host     string `json:"snip"`
	Port     int    `json:"snport"`
	Username string `json:"snusername"`
	Password string `json:"snpassword"`
}

func NewClient(options ConnectOptions) Client {
	client := ClientImpl{ConnectOptions: options, client: resty.New()}
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
	if resp.StatusCode() != 200{
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
	url := fmt.Sprintf("https://%s:%d/rest/measurements/list", c.Host, c.Port)
	tagfilter := fmt.Sprintf("<valuetagfilter filtervalue=\"%d\"><tagcategory key=\"Device ID\" id=\"61\"/></valuetagfilter>", deviceObid)
	resp, err := c.client.R().SetBody([]byte(tagfilter)).SetHeader("Content-Type", "application/xml").Post(url)
	if err != nil || resp.StatusCode() != 200 {
		return nil, fmt.Errorf("could not retrieve Measurements for device %d from StableNet", deviceObid)
	}
	return c.unmarshalMeasurements(bytes2.NewReader(resp.Body()))
}

func (c *ClientImpl) unmarshalMeasurements(reader io.Reader) ([]Measurement, error) {
	buff := new(bytes2.Buffer)
	_, err := buff.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read from reader: %v", err)
	}
	type measurementCollection struct {
		XMLName      xml.Name
		Measurements []Measurement `xml:",any"`
	}
	bytes := buff.Bytes()
	collections := measurementCollection{}
	err = xml.Unmarshal(bytes, &collections)
	return collections.Measurements, err
}

func (c *ClientImpl) FetchMetricsForMeasurement(measurementObid int, startTime time.Time, endTime time.Time) ([]string, error) {
	startMillis := startTime.UnixNano() / int64(time.Millisecond)
	endMillis := endTime.UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("https://%s:%d/StatisticServlet?stat=1020&type=json&login=%s,%s&id=%d&start=%d&end=%d", c.Host, c.Port, c.Username, c.Password, measurementObid, startMillis, endMillis)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics for measurement %d from StableNet", measurementObid)
	}
	data := make([]map[string]interface{}, 0, 0)
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	metricNames := make(map[string]bool)
	for _, record := range data {
		for key, _ := range record {
			metricNames[key] = true
		}
	}
	delete(metricNames, "Time")
	delete(metricNames, "time")

	result := make([]string, 0, len(metricNames))
	for key, _ := range metricNames {
		result = append(result, key)
	}
	return result, nil
}

func (c *ClientImpl) FetchDataForMetric(measurementObid int, metricName string, startTime time.Time, endTime time.Time) ([]MetricData, error) {
	startMillis := startTime.UnixNano() / int64(time.Millisecond)
	endMillis := endTime.UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("https://%s:%d/StatisticServlet?stat=1020&type=json&login=%s,%s&id=%d&start=%d&end=%d", c.Host, c.Port, c.Username, c.Password, measurementObid, startMillis, endMillis)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve metrics for measurement %d from StableNet", measurementObid)
	}
	data := make([]map[string]interface{}, 0, 0)
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	result := make([]MetricData, 0, len(data))
	timeFormat := "2006-01-02 15:04:05 -0700"
	for _, record := range data {
		measurementTime, err := time.Parse(timeFormat, record["Time"].(string))
		if err != nil {
			return nil, err
		}
		for key, value := range record {
			if key != metricName {
				continue
			}
			floatString := value.(string)
			floatString = strings.Replace(floatString, ",", "", -1)
			value, err := strconv.ParseFloat(floatString, 64)
			if err != nil {
				return nil, fmt.Errorf("could not format value: %v", err)
			}
			result = append(result, MetricData{
				Time:  measurementTime,
				Value: value,
			})
		}
	}
	return result, nil
}
