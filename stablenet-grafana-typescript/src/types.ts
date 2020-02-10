interface BasicQuery {
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
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
}

export interface RequestArgQuery {
  from: string;
  to: string;
  queries: SingleQuery[];
}
