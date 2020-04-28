import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
  constant: number;
}

export const defaultQuery: Partial<MyQuery> = {
  constant: 6.5,
};

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
