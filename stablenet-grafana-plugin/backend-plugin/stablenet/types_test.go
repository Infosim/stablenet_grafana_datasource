package stablenet

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestMetricDataSeries(t *testing.T) {
	const length = 10
	series := make([]MetricData, 0, length)
	wantedMin := make([]*datasource.Point, 0, length)
	wantedMax := make([]*datasource.Point, 0, length)
	wantedAvg := make([]*datasource.Point, 0, length)
	aMoment := time.Now()
	for i := 0; i < length; i++ {
		min := float64(i * 1000)
		avg := float64(i * 2000)
		max := float64(i * 3000)
		series = append(series, MetricData{
			Time: aMoment,
			Min:  min,
			Avg:  avg,
			Max:  max,
		})
		unix := aMoment.UnixNano() / int64(time.Millisecond)
		wantedMin = append(wantedMin, &datasource.Point{Timestamp: unix, Value: min})
		wantedAvg = append(wantedAvg, &datasource.Point{Timestamp: unix, Value: avg})
		wantedMax = append(wantedMax, &datasource.Point{Timestamp: unix, Value: max})
		aMoment = aMoment.Add(5 * time.Minute)
	}
	test := assert.New(t)
	dataSeries := MetricDataSeries(series)
	test.Equal(wantedMin, dataSeries.MinValues(), "min Values differ")
	test.Equal(wantedAvg, dataSeries.AvgValues(), "avg Values differ")
	test.Equal(wantedMax, dataSeries.MaxValues(), "max Values differ")
}

func TestMetricDataSeries_ExpandWithMissingValues(t *testing.T) {
	uptimes, err := ioutil.ReadFile("./test-data/uptimes.json")
	require.NoError(t, err)
	statisticData, err := parseStatisticByteSlice(uptimes)
	require.NoError(t, err)
	uptimeStatistic := statisticData["System Uptime"]
	expanded := uptimeStatistic.ExpandWithMissingValues()
	assert.Equal(t, 247, len(expanded))
}
