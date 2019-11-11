package stablenet

import "encoding/xml"

type Device struct {
	XMLName xml.Name
	Name string `xml:"name,attr" json:"name"`
	Obid int `xml:"obid,attr" json:"obid"`
}
