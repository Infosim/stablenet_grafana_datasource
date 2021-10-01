/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataQuery, DataQueryResponse, DataSourceJsonData, SelectableValue } from '@grafana/data';
import { of, Observable } from 'rxjs';

/**
 * This interface's structure is optimized for the config panel (it represents its state). However, we send this whole
 * object to the server because Grafana doesn't allow having different types for the data query and the config panel.
 * Ideally, the data should be organized and stripped before sending the request, but that isn't possible. As a workaround,
 * we do the stripping and reorganizing in the server instead.
 */
export interface Target extends DataQuery {
  mode: number;
  selectedDevice: LabelValue;
  selectedMeasurement: LabelValue;
  measurementFilter: string;
  chosenMetrics: string[];
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
