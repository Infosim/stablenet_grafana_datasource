package stablenet

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"math"
	"sort"
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

func (s MetricDataSeries) ExpandWithMissingValues() MetricDataSeries {
	if len(s) < 2 {
		return s
	}
	result := make([]MetricData, 0, len(s))
	s.sortAscAfterTimestamp()
	currentIndex := len(s) - 1
	for currentIndex >= 0 {
		currentInterval := s[currentIndex].Interval
		result = append(result, s[currentIndex])
		threshold := s[currentIndex].Time.Add(-currentInterval)
		currentIndex = currentIndex - 1
		if currentIndex >= 0 && s[currentIndex].Time.Before(threshold) {
			result = append(result, MetricData{Time: threshold.Add(currentInterval), Avg: math.NaN(), Min: math.NaN(), Max: math.NaN()})
			for currentIndex >= 0 && s[currentIndex].Time.Before(threshold) {
				threshold = threshold.Add(-currentInterval)
			}
			result = append(result, MetricData{Time: threshold.Add(currentInterval), Avg: math.NaN(), Min: math.NaN(), Max: math.NaN()})
		}
	}
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func (s MetricDataSeries) sortAscAfterTimestamp() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Time.Before(s[j].Time)
	})
}
