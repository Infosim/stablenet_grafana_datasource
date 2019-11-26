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

func Test_parseSingleTimestamp_Errors(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]string
		wantedError string
	}{
		{name: "missing timestamp", input: map[string]string{}, wantedError: "dataset did not contain a value for TIMESTAMP"},
		{name: "missing interval", input: map[string]string{"TIMESTAMP": "2019-11-15 11:56:42 +0100"}, wantedError: "dataset did not contain a value for INTERVAL"},
		{name: "wrong timestamp", input: map[string]string{"TIMESTAMP": "15.11.2019 11 Uhr 56", "INTERVAL": "00:05:00"}, wantedError: "invalid timestamp format: parsing time \"15.11.2019 11 Uhr 56\" as \"2006-01-02 15:04:05 -0700\": cannot parse \"1.2019 11 Uhr 56\" as \"2006\""},
		{name: "unparsable float", input: map[string]string{"TIMESTAMP": "2019-11-15 11:56:42 +0100", "time": "no float", "INTERVAL": "00:05:00"}, wantedError: "cannot parse value for \"no float\": strconv.ParseFloat: parsing \"nofloat\": invalid syntax"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseSingleTimestamp(test.input)
			require.Nil(t, actual, "there should be now result because an error should be returned")
			assert.EqualError(t, err, test.wantedError, "error message is not correct")
		})
	}
}

func Test_parseSingleTimestamp(t *testing.T) {
	input := map[string]string{
		"TIMESTAMP":                    "2019-11-15 11:56:42 +0100",
		"INTERVAL":                     "00:01:00",
		"OneDrive Time":                "1,277.000",
		"min OneDrive Time":            "1,200.000",
		"max OneDrive Time":            "1,300.000",
		"Script Execution Success":     "0",
		"min Script Execution Success": "0",
		"max Script Execution Success": "0",
		"Total Time":                   "1,949.000",
		"min Total Time":               "1,800.000",
		"max Total Time":               "2,000.000",
	}
	actual, err := parseSingleTimestamp(input)
	require.NoError(t, err, "no error expected")
	require.NotNil(t, actual["OneDrive Time"], "OneDrive Measurement Data is not present")
	require.NotNil(t, actual["Script Execution Success"], "Script Execution Success Measurement Data is not present")
	require.NotNil(t, actual["Total Time"], "Total time Measurement Data is not present")
	test := assert.New(t)
	oneDrive := actual["OneDrive Time"]
	script := actual["Script Execution Success"]
	totalTime := actual["Total Time"]
	testTime, _ := time.Parse(timeFormat, "2019-11-15 11:56:42 +0100")
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

func Test_getTypeAndKey(t *testing.T) {
	const testValue = 5.55
	checkMax := func(data *MetricData) bool {
		return data.Max == testValue && data.Avg == 0 && data.Min == 0
	}
	checkMin := func(data *MetricData) bool {
		return data.Max == 0 && data.Avg == 0 && data.Min == testValue
	}
	checkAvg := func(data *MetricData) bool {
		return data.Max == 0 && data.Avg == testValue && data.Min == 0
	}
	tests := []struct {
		name    string
		arg     string
		want    string
		checker func(*MetricData) bool
	}{
		{name: "check max", arg: "max OneDrive Time", want: "OneDrive Time", checker: checkMax},
		{name: "check Max", arg: "Max OneDrive Time", want: "OneDrive Time", checker: checkMax},
		{name: "check avg", arg: "OneDrive Time", want: "OneDrive Time", checker: checkAvg},
		{name: "check avg", arg: "Avg OneDrive Time", want: "Avg OneDrive Time", checker: checkAvg}, //cutting away avg is currently not necessary
		{name: "check min", arg: "min OneDrive Time", want: "OneDrive Time", checker: checkMin},
		{name: "check Min", arg: "min OneDrive Time", want: "OneDrive Time", checker: checkMin},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, consumer := getTypeAndKey(tt.arg)
			assert.Equal(t, got, tt.want, "name is not correct")
			data := MetricData{}
			consumer(&data, testValue)
			assert.True(t, tt.checker(&data), "the correct value must be set")
		})
	}
}

func Test_parseMeasurementData(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want float64
	}{
		{name: "simple", arg: "1200", want: 1200.00},
		{name: "with dot", arg: "1.200", want: 1.2},
		{name: "with comma", arg: "1,200", want: 1200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMeasurementData(tt.arg)
			require.NoError(t, err, "no error is expected")
			assert.Equal(t, tt.want, got, "parsed number is incorrect")
		})
	}
}

func Test_parseInterval(t *testing.T) {
	tests := []struct {
		arg  string
		want time.Duration
	}{
		{arg: "00:05:00", want: 5 * time.Minute},
		{arg: "12:42:13", want: 12*time.Hour + 42*time.Minute + 13*time.Second},
		{arg: "8544:09:01", want: 8544*time.Hour + 9*time.Minute + 1*time.Second},
	}
	for _, tt := range tests {
		t.Run("Test "+tt.arg, func(t *testing.T) {
			got, err := parseInterval(tt.arg)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
