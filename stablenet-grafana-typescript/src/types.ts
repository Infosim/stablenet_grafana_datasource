interface ArgBasic {
  refId: string;
  datasourceId: number;
  queryType: string;
}

export interface TestOptions {
  headers: object;
  url: string;
  method: string;
  data: {
    queries: ArgBasic[];
  };
}

export interface QueryDeviceOptions extends ArgBasic {
  filter: string;
}

export interface FindMeasurementOptions extends ArgBasic {
  filter: string;
  deviceObid: number;
}

export interface FindMetricsOptions extends ArgBasic {
  measurementObid: number;
}

export interface SingleQuery extends ArgBasic {
  statisticLink?: string;
  requestData?: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }>;
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
}

export interface RequestArgStandard {
  queries: QueryDeviceOptions[] | FindMeasurementOptions[] | FindMetricsOptions[];
}

export interface RequestArgQuery {
  from: string;
  to: string;
  queries: SingleQuery[];
}
