package stablenet

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseRawData(data []map[string]string) (map[string]MetricDataSeries, error) {
	return nil, nil
}

const timestampKey = "TIMESTAMP"
const intervalKey = "INTERVAL"
const timeFormat = "2006-01-02 15:04:05 -0700"
const intervalFormat = "15:04:05"

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
	_, intervalErr := time.Parse(intervalFormat, data[intervalKey])
	if intervalErr != nil {
		return nil, fmt.Errorf("invalid interval format: %v", intervalErr)
	}
	result := make(map[string]MetricData)
	for key, stringVal := range data {
		if key == timestampKey || key == intervalKey {
			continue
		}
		value, formatErr := parseMeasurementData(stringVal)
		if formatErr != nil {
			return nil, fmt.Errorf("cannot parse value for \"%s\": %v", stringVal, formatErr)
		}
		name, consumer := getTypeAndKey(key)
		if _, ok := result[name]; !ok {
			result[name] = MetricData{
				Time: measurementTime,
			}
		}
		measurementData := result[name]
		consumer(&measurementData, value)
		result[name] = measurementData
	}
	return result, nil
}

func getTypeAndKey(key string) (string, func(*MetricData, float64)) {
	lowerkey := strings.ToLower(key)
	if strings.HasPrefix(lowerkey, "min") {
		return key[4:], func(data *MetricData, f float64) {
			data.Min = f
		}
	} else if strings.HasPrefix(lowerkey, "max") {
		return key[4:], func(data *MetricData, f float64) {
			data.Max = f
		}
	}
	return key, func(data *MetricData, f float64) {
		data.Avg = f
	}
}

func parseMeasurementData(value string) (float64, error) {
	value = strings.Replace(value, " ", "", -1)
	value = strings.Replace(value, ",", "", -1)
	return strconv.ParseFloat(value, 64)
}
