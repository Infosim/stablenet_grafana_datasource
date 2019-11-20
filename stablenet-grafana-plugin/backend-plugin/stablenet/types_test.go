package stablenet

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/stretchr/testify/assert"
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
	startTime,_ := time.Parse(timeFormat, "2019-11-15 12:00:00 +0100")
	data := []MetricData{
		{Time: startTime, Avg: 42},
		{Time: startTime.Add(1 * time.Hour), Interval: 1 * time.Hour, Avg: 42},
		{Time: startTime.Add(1*time.Hour + 55*time.Minute), Interval: 55 * time.Minute, Avg: 42},
		{Time: startTime.Add(8 * time.Hour), Interval: 1 * time.Hour, Avg: 42},
		{Time: startTime.Add(9 * time.Hour), Interval: 1 * time.Hour, Avg: 42},
	}
	actual := MetricDataSeries(data).ExpandWithMissingValues()
	assert.Equal(t, 11, len(actual))
}
