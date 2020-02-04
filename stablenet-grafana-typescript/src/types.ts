export interface SingleQuery {
  refId: string;
  datasourceId: number;
  queryType: string;
  statisticLink?: string;
  requestData?: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }>;
  includeMinStats: boolean;
  includeAvgStats: boolean;
  includeMaxStats: boolean;
}

interface TargetEmpty {
  refId: string;
  datasource: undefined;
}

export interface QueryOptionsEmpty {
  targets: [TargetEmpty];
  [x: string]: any;
}

export interface Target {
  refId: string;
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
  metrics: Array<{ text: string; value: string; measurementObid: number; $$hashKey: string }>;
  moreDevices: boolean;
  moreMeasurements: boolean;
  datasource: any;
}

export interface QueryOptions {
  targets: Target[];
  [x: string]: any;
}

export interface QueryResult {
  data: Array<{ target: string; datapoints: Array<[number, number]> }> | never[];
  status?: number;
  headers?: any;
  config?: any;
  statusText?: string;
  xhrStatus?: string;
}

export function isQOE(object: any): object is QueryOptionsEmpty {
  return !('mode' in object.targets[0]);
}
