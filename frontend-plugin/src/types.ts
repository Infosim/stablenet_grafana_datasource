/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataSourceJsonData, SelectableValue } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface Metric {
  text: string;
  key: string;
  measurementObid: number;
}

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
  metrics: Metric[];
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

export enum Mode {
  MEASUREMENT = 0,
  STATISTIC_LINK = 10,
}

export enum Unit {
  SECONDS = 1_000,
  MINUTES = 60_000,
  HOURS = 3_600_000,
  DAYS = 86_400_000,
}

/** These are options configured for each StableNetDataSource instance */
export interface StableNetConfigOptions extends DataSourceJsonData {
  ip?: string;
  port?: number;
}

/** Value that is used in the backend, but never sent over HTTP to the frontend */
export interface StableNetSecureJsonData {
  password?: string;
}
