export interface TextValue {
  text: string;
  value: number;
}

export interface TargetDatapoints {
  target: string;
  datapoints: Array<[number, number]>;
}

export interface TestResult {
  status: string;
  message: string;
  title: string;
}

export interface QueryResult {
  data: TextValue[];
  hasMore: boolean;
}

export interface MetricResult extends TextValue {
  measurementObid: number;
}

export interface QueryResultEmpty {
  data: never[];
}

interface RequestResult {
  status: number;
  headers: Function;
  config: object;
  statusText: string;
  xhrStatus: string;
}

export interface MeasurementQueryResult {
  hasMore: boolean;
  data: { name: string; obid: number };
}

export interface MetricType {
  key: string;
  name: string;
}

/**
 * An alternative return type of "doRequest()" when it is NOT called from "query()".
 * Cannot use now as only the "findMetricsForMeasurements()" path for some reason returns
 * result.data.results[refId].meta
 * instead of
 * result.data.results[refId].meta.data
 * ...as an Array. See what Backend can do.
 */
export interface RequestResultStandard<T> extends RequestResult {
  data: {
    results: {
      [x: string]: {
        refId: string;
        tables: null;
        series: never[];
        meta: T;
      };
    };
  };
}

export interface GenericResponse {
  data: { results: object };
  status: number;
  headers: any;
  config: any;
  statusText: string;
  xhrStatus: string;
}

export interface TSDBArg extends RequestResult {
  data: {
    results: {
      [x: string]: {
        refId: string;
        tables: null;
        series: Array<{ name: string; points: Array<[number, number]> }>;
      };
    };
  };
}

export interface TSDBResult extends RequestResult {
  data: TargetDatapoints[];
}
