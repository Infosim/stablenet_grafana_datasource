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

//func Test_parseSingleTimestamp_Errors(t *testing.T) {
//	tests := []struct {
//		name        string
//		input       map[string]string
//		wantedError string
//	}{
//		{name: "missing timestamp", input: map[string]string{}, wantedError: "dataset did not contain a value for TIMESTAMP"},
//		{name: "missing interval", input: map[string]string{"TIMESTAMP": "2019-11-15 11:56:42 +0100"}, wantedError: "dataset did not contain a value for INTERVAL"},
//		{name: "wrong timestamp", input: map[string]string{"TIMESTAMP": "15.11.2019 11 Uhr 56", "INTERVAL": "00:05:00"}, wantedError: "invalid timestamp format: parsing time \"15.11.2019 11 Uhr 56\" as \"2006-01-02 15:04:05 -0700\": cannot parse \"1.2019 11 Uhr 56\" as \"2006\""},
//		{name: "unparsable float", input: map[string]string{"TIMESTAMP": "2019-11-15 11:56:42 +0100", "time": "no float", "INTERVAL": "00:05:00"}, wantedError: "cannot parse value for \"no float\": strconv.ParseFloat: parsing \"nofloat\": invalid syntax"},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			actual, err := parseSingleTimestamp(test.input)
//			require.Nil(t, actual, "there should be now result because an error should be returned")
//			assert.EqualError(t, err, test.wantedError, "error message is not correct")
//		})
//	}
//}

func pointer(value float64) *float64 {
	result := value
	return &result
}

func Test_parseSingleTimestamp(t *testing.T) {
	input := timestampResponse{
		Interval: 60000,
		Row: []measurementData{
			{Min: pointer(1200), Avg: pointer(1277), Max: pointer(1300)},
			{Min: pointer(0), Avg: pointer(0), Max: pointer(0)},
			{Min: pointer(1800), Avg: pointer(1949), Max: pointer(2000)},
		},
		TimeStamp: 1574839297028,
	}
	names := []string{"OneDrive", "Script Execution Success", "Total Time"}
	actual, err := parseSingleTimestamp(input, names)
	require.NoError(t, err, "no error expected")
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
