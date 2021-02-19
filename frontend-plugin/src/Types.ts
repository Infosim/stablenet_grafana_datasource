/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataQuery, DataQueryResponse, DataSourceJsonData, SelectableValue } from '@grafana/data';
import { of, Observable } from 'rxjs';

export interface Target extends DataQuery {
  mode: number;
  selectedDevice: LabelValue;
  selectedMeasurement: LabelValue;
  measurementFilter: string;
  chosenMetrics: object;
  metricPrefix: string;
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
  statisticLink: string;
  averagePeriod: string;
  averageUnit: number;
  useCustomAverage: boolean;
  measurements: LabelValue[];
  metrics: Array<{ text: string; key: string; measurementObid: number }>;
  moreDevices: boolean;
  moreMeasurements: boolean;
}

export interface LabelValue extends SelectableValue<number> {
  label: string;
  value: number;
}

export interface TestResult {
  status: string;
  message: string;
  title: string;
}

export interface QueryResult {
  data: LabelValue[];
  hasMore: boolean;
}

export interface MetricResult {
  key: string;
  text: string;
  measurementObid: number;
}

export const EmptyResult: Observable<DataQueryResponse> = of({ data: [] });

export interface SingleQuery extends DataQuery {
  statisticLink?: string;
  measurementObid?: number;
  metrics?: Array<{ key: string; name: string }>;
  // Do not use the name intervalMs here because this property gets overridden by Grafana.
  // We want to use our own average period.
  customInterval: number;
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
}

export interface StringPair {
  key: string;
  name: string;
}

export enum Mode {
  MEASUREMENT = 0,
  STATISTIC_LINK = 10,
}

export enum Unit {
  SECONDS = 1000,
  MINUTES = 60000,
  HOURS = 3600000,
  DAYS = 86400000,
}

/**
 * These are options configured for each StableNetDataSource instance
 */
export interface StableNetConfigOptions extends DataSourceJsonData {
  snip?: string;
  snport?: string;
  snusername?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface StableNetSecureJsonData {
  snpassword?: string;
}
