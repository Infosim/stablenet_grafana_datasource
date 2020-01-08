/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package query

import (
	"backend-plugin/stablenet"
	"fmt"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"regexp"
	"strconv"
	"strings"
)

func findMeasurementIdsInLink(link string) map[int]int {
	measurementRegex := regexp.MustCompile("[?&](\\d*)id=(\\d+)")
	idMatches := measurementRegex.FindAllStringSubmatch(link, -1)
	result := make(map[int]int)
	for _, match := range idMatches {
		index, _ := strconv.Atoi(match[1])
		measurementId, _ := strconv.Atoi(match[2])
		result[index] = measurementId
	}
	return result
}

func extractMetricKeysForMeasurements(link string) map[int][]string {
	measurementIds := findMeasurementIdsInLink(link)
	keyRegex := regexp.MustCompile("[?&](\\d*)value\\d*=(\\d+)")
	keyMatches := keyRegex.FindAllStringSubmatch(link, -1)
	result := make(map[int][]string)
	for _, match := range keyMatches {
		index, _ := strconv.Atoi(match[1])
		measurementId, ok := measurementIds[index]
		if !ok {
			continue
		}
		list, isPresent := result[measurementId]
		if !isPresent {
			list = make([]string, 0, 0)
		}
		list = append(list, match[2])
		result[measurementId] = list
	}
	return result
}

func filterWantedMetrics(fromLink []string, realMetrics []stablenet.Metric) []stablenet.Metric {
	result := make([]stablenet.Metric, 0, 0)
	for _, realMetric := range realMetrics {
		for _, requestedMetric := range fromLink {
			if strings.Contains(realMetric.Key, requestedMetric) {
				result = append(result, realMetric)
			}
		}
	}
	return result
}

type statisticLinkHandler struct {
	*StableNetHandler
}

func (s statisticLinkHandler) Process(query Query) (*datasource.QueryResult, error) {
	link, err := query.GetCustomField("statisticLink")
	if err != nil {
		return BuildErrorResult("could not extract statisticLink parameter from query", query.RefId), nil
	}
	requested := extractMetricKeysForMeasurements(link)
	if len(requested) == 0 {
		return BuildErrorResult(fmt.Sprintf("the link \"%s\" does not carry a measurement id or value ids", link), query.RefId), nil
	}
	allSeries := make([]*datasource.TimeSeries, 0, 0)
	for measurementId, metricKeys := range requested {
		realMetrics, err := s.SnClient.FetchMetricsForMeasurement(measurementId, "")
		if err != nil {
			s.Logger.Error(err.Error())
			return BuildErrorResult(fmt.Sprintf("could not fetch metrics for measurement %d, does it exist?", measurementId), query.RefId), nil
		}
		metrics := filterWantedMetrics(metricKeys, realMetrics)
		if len(metrics) == 0 {
			continue
		}

		series, err := s.fetchMetrics(query, measurementId, metrics)
		if err != nil {
			e := fmt.Errorf("could not fetch data for statistic link from server: %v", err)
			s.Logger.Error(e.Error())
			return nil, e
		}
		measurementName, err := s.SnClient.FetchMeasurementName(measurementId)
		if err != nil {
			s.Logger.Error(err.Error())
			return BuildErrorResult(fmt.Sprintf("could not fetch name of measurement %d. See Logs for more information", measurementId), query.RefId), nil
		}
		for _, singleSeries := range series {
			singleSeries.Name = *measurementName + " " + singleSeries.Name
			allSeries = append(allSeries, singleSeries)
		}
	}
	result := datasource.QueryResult{
		RefId:  query.RefId,
		Series: allSeries,
	}
	return &result, nil
}
