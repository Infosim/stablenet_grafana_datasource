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
	valueKeys := make(map[int][]string)
	for _, match := range keyMatches {
		index, _ := strconv.Atoi(match[1])
		list, isPresent := valueKeys[index]
		if !isPresent {
			list = make([]string, 0, 0)
		}
		list = append(list, match[2])
		valueKeys[index] = list
	}
	result := make(map[int][]string)
	for index, measurementId := range measurementIds {
		if _, ok := result[measurementId]; !ok {
			result[measurementId] = make([]string, 0, 0)
		}
		list, isPresent := valueKeys[index]
		if !isPresent {
			result[measurementId] = make([]string, 0, 0)
		} else {
			result[measurementId] = list
		}
	}
	return result
}

func filterWantedMetrics(fromLink []string, realMetrics []stablenet.Metric) []StringPair {
	result := make([]StringPair, 0, 0)
	for _, realMetric := range realMetrics {
		for _, requestedMetric := range fromLink {
			if len(fromLink) == 0 || strings.Contains(realMetric.Key, requestedMetric) {
				result = append(result, StringPair{Key: realMetric.Key, Name: realMetric.Name})
			}
		}
	}
	return result
}

func ParseStatisticLink(originalQuery MetricQuery, metricSupplier func(int) ([]stablenet.Metric, error)) ([]MetricQuery, error) {
	requested := extractMetricKeysForMeasurements(*originalQuery.StatisticLink)
	if len(requested) == 0 {
		return nil, fmt.Errorf("the link \"%s\" does not carry at least a measurement id", *originalQuery.StatisticLink)
	}
	allQueries := make([]MetricQuery, 0, 0)
	for measurementId, metricKeys := range requested {
		realMetrics, err := metricSupplier(measurementId)
		if err != nil {
			return nil, fmt.Errorf("could not fetch metrics for measurement %d: %v", measurementId, err)
		}
		metrics := filterWantedMetrics(metricKeys, realMetrics)
		if len(metrics) == 0 {
			continue
		}

		query := originalQuery.ShallowClone()
		query.Metrics = metrics
		query.MeasurementObid = measurementId

		allQueries = append(allQueries, query)
	}
	return allQueries, nil
}
