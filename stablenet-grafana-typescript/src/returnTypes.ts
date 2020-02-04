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

export interface FindResult {
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

/**
 * An alternative return type of "doRequest()" when it is NOT called from "query()".
 * Cannot use now as only the "findMetricsForMeasurements()" path for some reason returns
 * result.data.results[refId].meta
 * instead of
 * result.data.results[refId].meta.data
 * ...as an Array. See what Backend can do.
 */
export interface RequestResultStandard extends RequestResult {
  data: {
    results: {
      [x: string]: {
        refId: string;
        tables: null;
        series: never[];
        meta: {
          hasMore: boolean;
          data: Array<{ name: string; obid: number }> | Array<{ key: string; name: string }>;
        };
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
