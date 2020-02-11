export interface TextValue {
  text: string;
  value: string;
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

export interface EmptyQueryResult {
  data: never[];
}

interface RequestResult {
  status: number;
  headers: Function;
  config: object;
  statusText: string;
  xhrStatus: string;
}

export interface EntityQueryResult {
  hasMore: boolean;
  data: Array<{ name: string; obid: number }>;
}

export interface MetricType {
  key: string;
  name: string;
}

export interface GenericResponse<T> extends RequestResult {
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
