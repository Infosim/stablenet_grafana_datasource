package stablenet

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"time"
)

type Device struct {
	Name string `json:"name"`
	Obid int    `json:"obid"`
}

type Measurement struct {
	Name string `json:"name"`
	Obid int    `json:"obid"`
}

type Metric struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type MetricData struct {
	Interval time.Duration
	Time     time.Time
	Min      float64
	Max      float64
	Avg      float64
}

type MetricDataSeries []MetricData

func (s MetricDataSeries) MinValues() []*datasource.Point {
	return s.toValues(func(data MetricData) float64 {
		return data.Min
	})
}

func (s MetricDataSeries) toValues(selector func(MetricData) float64) []*datasource.Point {
	result := make([]*datasource.Point, 0, len(s))
	for _, metricData := range s {
		point := datasource.Point{
			Timestamp: metricData.Time.UnixNano() / int64(1000*time.Microsecond),
			Value:     selector(metricData),
		}
		result = append(result, &point)
	}
	return result
}

func (s MetricDataSeries) MaxValues() []*datasource.Point {
	return s.toValues(func(data MetricData) float64 {
		return data.Max
	})
}

func (s MetricDataSeries) AvgValues() []*datasource.Point {
	return s.toValues(func(data MetricData) float64 {
		return data.Avg
	})
}
