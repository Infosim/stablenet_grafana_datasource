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
)

type SnServer struct {
	Username     string
	Password     string
	Devices      []stablenet.Device
	Measurements []stablenet.Measurement
	Metrics      []stablenet.Metric
	Data         []stablenet.TimestampResponse
	Info         stablenet.ServerInfo
	LastQueries  url.Values
}

func CreateMockServer(username, password string) *SnServer {
	five := 5.0
	ten := 10.0
	avg := 7.5

	return &SnServer{
		Username: username,
		Password: password,
		Devices: []stablenet.Device{
			{Obid: 9000, Name: "Bach"},
			{Obid: 9001, Name: "Fluss"},
			{Obid: 9002, Name: "Meer"},
		},
		Measurements: []stablenet.Measurement{
			{Obid: 1001, Name: "Host"},
			{Obid: 1002, Name: "Processor"},
			{Obid: 1003, Name: "Interface 1"},
		},
		Metrics: []stablenet.Metric{
			{Name: "Uptime", Key: "SNMP_1"},
			{Name: "CPU 1", Key: "EXTERN_2"},
		},
		Data: []stablenet.TimestampResponse{
			{
				TimeStamp: 100,
				Row: []stablenet.MeasurementData{
					{Min: &five, Max: &ten, Avg: &avg},
				},
			},
		},
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
