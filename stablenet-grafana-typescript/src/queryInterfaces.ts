interface TargetEmpty {
  refId: string;
  datasource: undefined;
}

export interface QueryOptionsEmpty {
  targets: [TargetEmpty];
  [x: string]: any;
}

export interface Target extends TargetEmpty {
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
}

export interface QueryOptions {
  targets: Target[];
  [x: string]: any;
}

export function isQOE(object: any): object is QueryOptionsEmpty {
  return !('mode' in object.targets[0]);
}
