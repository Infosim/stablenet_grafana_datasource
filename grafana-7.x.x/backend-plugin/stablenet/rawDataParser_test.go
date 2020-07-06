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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_parseSingleTimestamp(t *testing.T) {
	input := TimestampResponse{
		Interval: 60000,
		Row: []MeasurementData{
			{Min: ptr(1200), Avg: ptr(1277), Max: ptr(1300)},
			{Min: ptr(0), Avg: ptr(0), Max: ptr(0)},
			{Min: ptr(1800), Avg: ptr(1949), Max: ptr(2000)},
		},
		TimeStamp: 1574839297028,
	}
	names := []string{"OneDrive", "Script Execution Success", "Total Time"}
	actual := parseSingleTimestamp(input, names)
	require.NotNil(t, actual["OneDrive Time"], "OneDrive Measurement Data is not present")
	require.NotNil(t, actual["Script Execution Success"], "Script Execution Success Measurement Data is not present")
	require.NotNil(t, actual["Total Time"], "Total time Measurement Data is not present")
	test := assert.New(t)
	oneDrive := actual["OneDrive"]
	script := actual["Script Execution Success"]
	totalTime := actual["Total Time"]
	testTime := time.Unix(0, input.TimeStamp*int64(time.Millisecond))
	expectedOneDrive := MetricData{
		Interval: 1 * time.Minute,
		Time:     testTime,
		Min:      1200,
		Max:      1300,
		Avg:      1277,
	}
	assertMetricDataCorrect(test, expectedOneDrive, oneDrive, "One Drive")
	expectedScript := MetricData{
		Interval: 1 * time.Minute,
		Time:     testTime,
		Min:      0,
		Max:      0,
		Avg:      0,
	}
	assertMetricDataCorrect(test, expectedScript, script, "Script Execution Success")
	expectedTotalTime := MetricData{
		Interval: 1 * time.Minute,
		Time:     testTime,
		Min:      1800,
		Max:      2000,
		Avg:      1949,
	}
	assertMetricDataCorrect(test, expectedTotalTime, totalTime, "Total Time")
}

func assertMetricDataCorrect(test *assert.Assertions, expected MetricData, actual MetricData, msg string) {
	test.Equal(expected.Time, actual.Time, fmt.Sprintf("%s: time is different", msg))
	test.Equal(expected.Interval, actual.Interval, fmt.Sprintf("%s: interval is different", msg))
	test.Equal(expected.Min, actual.Min, fmt.Sprintf("%s: min is different", msg))
	test.Equal(expected.Max, actual.Max, fmt.Sprintf("%s: max is different", msg))
	test.Equal(expected.Avg, actual.Avg, fmt.Sprintf("%s: avg kis different", msg))
}

func ptr(value float64) *float64 {
	result := value
	return &result
}
