package stablenet

import (
	"encoding/xml"
	"time"
)

type Device struct {
	XMLName xml.Name
	Name    string `xml:"name,attr" json:"name"`
	Obid    int    `xml:"obid,attr" json:"obid"`
}

type Measurement struct {
	XMLName xml.Name
	Name    string `xml:"name,attr" json:"name"`
	Obid    int    `xml:"obid,attr" json:"obid"`
}

type MetricData struct {
	Time time.Time
	Value float64
}
