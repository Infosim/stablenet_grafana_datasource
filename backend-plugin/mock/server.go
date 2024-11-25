/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package mock

import (
	"backend-plugin/stablenet"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"time"
)

type MeasurementData struct {
	Min *float64 `json:"min"`
	Avg *float64 `json:"avg"`
	Max *float64 `json:"max"`
}

type TimestampResponse struct {
	Interval  int               `json:"interval"`
	TimeStamp int64             `json:"timestamp"`
	Row       []MeasurementData `json:"row"`
}

type SnServer struct {
	Username     string
	Password     string
	Devices      []stablenet.Device
	Measurements []stablenet.Measurement
	Metrics      []stablenet.Metric
	Data         stablenet.MeasurementMultiMetricResultDataDTO
	Info         stablenet.ServerInfo
	LastQueries  url.Values
}

var DefaultDevices = []stablenet.Device{
	{Obid: 9000, Name: "Bach"},
	{Obid: 9001, Name: "Fluss"},
	{Obid: 9002, Name: "Meer"},
}

var DefaultMeasurements = []stablenet.Measurement{
	{Obid: 1001, Name: "Host"},
	{Obid: 1002, Name: "Processor"},
	{Obid: 1003, Name: "Interface 1"},
}

var DefaultMetrics = []stablenet.Metric{
	{Name: "Uptime", Key: "SNMP_1"},
	{Name: "CPU 1", Key: "EXTERN_2"},
}

func floatPointer(v float64) *float64 {
	return &v
}

var defaultEntry = stablenet.MeasurementDataEntryDTO{
	Timestamp: time.Now().Local().UnixMilli(),
	Interval:  1000,
	Min:       floatPointer(5.0),
	Avg:       floatPointer(7.5),
	Max:       floatPointer(10.0),
}

var DefaultData = stablenet.MeasurementMultiMetricResultDataDTO{
	Values: []stablenet.MeasurementMetricResultDataDTO{
		{
			MetricKey: "SNMP_1",
			Data: []stablenet.MeasurementDataEntryDTO{
				defaultEntry,
			},
		},
	},
}

func CreateMockServer(username, password string) *SnServer {
	return &SnServer{
		Username:     username,
		Password:     password,
		Devices:      DefaultDevices,
		Measurements: DefaultMeasurements,
		Metrics:      DefaultMetrics,
		Data:         DefaultData,
	}
}

func (s *SnServer) getDevices(rw http.ResponseWriter, req *http.Request) {
	s.LastQueries = req.URL.Query()
	result := stablenet.DeviceQueryResult{Data: s.Devices, HasMore: false}
	payload, _ := json.Marshal(result)
	_, _ = rw.Write(payload)
}

func (s *SnServer) getMeasurements(rw http.ResponseWriter, req *http.Request) {
	s.LastQueries = req.URL.Query()
	result := stablenet.MeasurementQueryResult{Data: s.Measurements, HasMore: false}
	payload, _ := json.Marshal(result)
	_, _ = rw.Write(payload)
}

func (s *SnServer) getMetrics(rw http.ResponseWriter, req *http.Request) {
	s.LastQueries = req.URL.Query()
	payload, _ := json.Marshal(s.Metrics)
	_, _ = rw.Write(payload)
}

func (s *SnServer) getData(rw http.ResponseWriter, req *http.Request) {
	s.LastQueries = req.URL.Query()
	payload, _ := json.Marshal(s.Data)
	_, _ = rw.Write(payload)
}

func (s *SnServer) getInfo(rw http.ResponseWriter, req *http.Request) {
	s.LastQueries = req.URL.Query()
	payload, _ := xml.Marshal(s.Info)
	_, _ = rw.Write(payload)
}

func CreateHandler(server *SnServer) http.Handler {
	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			user, pass, ok := req.BasicAuth()
			if ok && user == server.Username && pass == server.Password {
				next.ServeHTTP(rw, req)
				return
			}
			http.Error(rw, "Authentication Error", http.StatusUnauthorized)
		}
	}

	r := http.NewServeMux()
	r.HandleFunc("/api/1/devices", authMiddleware(server.getDevices))
	r.HandleFunc("/api/1/measurements", authMiddleware(server.getMeasurements))
	r.HandleFunc("/api/1/measurement-data/1001/metrics", authMiddleware(server.getMetrics))
	r.HandleFunc("/api/1/measurement-data/1001", authMiddleware(server.getData))
	r.HandleFunc("/rest/info", authMiddleware(server.getInfo))

	return r
}
