/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
interface EmptyTarget {
  refId: string;
  datasource: undefined;
}

export interface Target extends EmptyTarget {
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
  metrics: Array<{ text: string; key: string; measurementObid: number }>;
  moreDevices: boolean;
  moreMeasurements: boolean;
}

export interface QueryOptions {
  targets: Target[];
  [x: string]: any;
}
