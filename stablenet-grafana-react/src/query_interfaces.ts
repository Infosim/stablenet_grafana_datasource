/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */

import {DataQuery} from "@grafana/data";

export interface Target extends DataQuery {
  mode: number;
  deviceQuery: string;
  selectedDevice: number;
  measurementQuery: string;
  selectedMeasurement: number;
  chosenMetrics: object;
  metricPrefix: string;
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
  statisticLink: string;
  averagePeriod: string;
  averageUnit: number;
  useCustomAverage: boolean;
  metrics: Array<{ text: string; key: string; measurementObid: number }>;
  moreDevices: boolean;
  moreMeasurements: boolean;
}
