/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
export interface BasicQuery {
  refId: string;
  datasourceId: number;
  queryType: string;
}

export interface TestOptions {
  headers: object;
  url: string;
  method: string;
  data: Query<BasicQuery>;
}

export interface DeviceQuery extends BasicQuery {
  filter: string;
}

export interface MeasurementQuery extends BasicQuery {
  filter: string;
  deviceObid: number;
}

export interface MetricQuery extends BasicQuery {
  measurementObid: number;
}

export interface Query<T> {
  from?: string;
  to?: string;
  queries: T[];
}

export interface SingleQuery extends BasicQuery {
  statisticLink?: string;
  requestData?: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }>;
  intervalMs: number;
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
}

export interface RequestArgQuery {
  from: string;
  to: string;
  queries: SingleQuery[];
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
  SECONDS = 1,
  MINUTES = 60,
  HOURS = 3600,
  DAYS = 86400,
}
