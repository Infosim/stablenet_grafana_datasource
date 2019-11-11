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
	client.client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	return &client
}

type ClientImpl struct {
	ConnectOptions
	client *resty.Client
}

func (c *ClientImpl) FetchAllDevices() ([]Device, error) {
	url := fmt.Sprintf("https://%s:%d/rest/devices/list", c.Host, c.Port)
	resp, err := c.client.R().Get(url)
	if err != nil{
		return nil, err
	}
	return c.unmarshalDevices(bytes2.NewReader(resp.Body()))
}

func (c *ClientImpl) unmarshalDevices(reader io.Reader) ([]Device, error) {
	buff := new(bytes2.Buffer)
	_, err := buff.ReadFrom(reader)
	if err != nil{
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
