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
  metrics: Array<{ text: string; key: string; measurementObid: number; $$hashKey: string }>;
  moreDevices: boolean;
  moreMeasurements: boolean;
}

export interface QueryOptions {
  targets: Target[];
  [x: string]: any;
}
