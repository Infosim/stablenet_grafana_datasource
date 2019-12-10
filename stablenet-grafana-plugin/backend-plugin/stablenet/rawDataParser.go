/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package stablenet

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const timestampKey = "Time"
const intervalKey = "Interval"
const timeFormat = "2006-01-02 15:04:05 -0700"

func parseSingleTimestamp(data map[string]string) (map[string]MetricData, error) {
	if _, ok := data[timestampKey]; !ok {
		return nil, fmt.Errorf("dataset did not contain a value for %s", timestampKey)
	}
	if _, ok := data[intervalKey]; !ok {
		return nil, fmt.Errorf("dataset did not contain a value for %s", intervalKey)
	}
	measurementTime, timeErr := time.Parse(timeFormat, data[timestampKey])
	if timeErr != nil {
		return nil, fmt.Errorf("invalid timestamp format: %v", timeErr)
	}
	interval, intervalErr := parseInterval(data[intervalKey])
	if intervalErr != nil {
		return nil, fmt.Errorf("invalid interval format: %v", intervalErr)
	}
	result := make(map[string]MetricData)
	for key, stringVal := range data {
		if key == timestampKey || key == intervalKey || measurementTime == time.Unix(0, 0) {
			continue
		}
		value, formatErr := parseMeasurementData(stringVal)
		if formatErr != nil {
			return nil, fmt.Errorf("cannot parse value for \"%s\": %v", stringVal, formatErr)
		}
		name, consumer := getTypeAndKey(key)
		if _, ok := result[name]; !ok {
			result[name] = MetricData{
				Time:     measurementTime,
				Interval: interval,
			}
		}
		measurementData := result[name]
		consumer(&measurementData, value)
		result[name] = measurementData
	}
	return result, nil
}

func getTypeAndKey(key string) (string, func(*MetricData, float64)) {
	lowerKey := strings.ToLower(key)
	if strings.HasPrefix(lowerKey, "min") {
		return key[4:], func(data *MetricData, f float64) {
			data.Min = f
		}
	} else if strings.HasPrefix(lowerKey, "max") {
		return key[4:], func(data *MetricData, f float64) {
			data.Max = f
		}
	}
	return key, func(data *MetricData, f float64) {
		data.Avg = f
	}
}

var durationRegex = regexp.MustCompile("(\\d+):(\\d\\d):(\\d\\d)")

func parseInterval(value string) (time.Duration, error) {
	millis, _ := strconv.Atoi(value)
		return time.Duration(millis)*time.Millisecond, nil
	
}

func parseMeasurementData(value string) (float64, error) {
	value = strings.Replace(value, " ", "", -1)
	value = strings.Replace(value, ",", "", -1)
	if value == "" {
		return math.NaN(), nil
	}
	return strconv.ParseFloat(value, 64)
}
