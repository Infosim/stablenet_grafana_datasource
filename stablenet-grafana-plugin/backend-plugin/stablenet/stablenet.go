package stablenet

import (
	bytes2 "bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
)

type Client interface {
	FetchAllDevices() ([]Device, error)
	FetchMeasurementsForDevice(int) ([]Measurement, error)
}

type ConnectOptions struct {
	Host     string
	Port     int
	Username string
	Password string
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

func (c *ClientImpl) FetchAllDevices() ([]Device, error) {
	url := fmt.Sprintf("https://%s:%d/rest/devices/list", c.Host, c.Port)
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return c.unmarshalDevices(bytes2.NewReader(resp.Body()))
}

func (c *ClientImpl) unmarshalDevices(reader io.Reader) ([]Device, error) {
	buff := new(bytes2.Buffer)
	_, err := buff.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read from reader: %v", err)
	}
	type deviceCollection struct {
		XMLName xml.Name
		Devices []Device `xml:"device"`
	}
	bytes := buff.Bytes()
	collections := deviceCollection{}
	err = xml.Unmarshal(bytes, &collections)
	return collections.Devices, err
}

func (c *ClientImpl) FetchMeasurementsForDevice(deviceObid int) ([]Measurement, error) {
	url := fmt.Sprintf("https://%s:%d/rest/measurements/list", c.Host, c.Port)
	tagfilter := fmt.Sprintf("<valuetagfilter filtervalue=\"%d\"><tagcategory key=\"Device ID\" id=\"61\"/></valuetagfilter>", deviceObid)
	resp, err := c.client.R().SetBody([]byte(tagfilter)).SetHeader("Content-Type", "application/xml").Post(url)
	if err != nil {
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
