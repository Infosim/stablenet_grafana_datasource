/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package stablenet

import (
	"math"
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

func parseSingleTimestamp(data TimestampResponse, metricKeys []string) map[string]MetricData {
	measurementTime := time.Unix(0, data.TimeStamp*int64(time.Millisecond))

	interval := time.Duration(data.Interval) * time.Millisecond

	result := make(map[string]MetricData)
	for index, row := range data.Row {
		value := MetricData{
			Min:      math.NaN(),
			Avg:      math.NaN(),
			Max:      math.NaN(),
			Interval: interval,
			Time:     measurementTime,
		}

		if row.Max != nil {
			value.Max = *row.Max
		}
		if row.Min != nil {
			value.Min = *row.Min
		}
		if row.Avg != nil {
			value.Avg = *row.Avg
		}

		result[metricKeys[index]] = value
	}

	return result
}
