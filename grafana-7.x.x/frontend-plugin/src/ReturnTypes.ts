/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { SelectableValue } from '@grafana/data';

export interface LabelValue extends SelectableValue<number> {
  label: string;
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
  data: LabelValue[];
  hasMore: boolean;
}

export interface MetricResult {
  key: string;
  text: string;
  measurementObid: number;
}

export interface EmptyQueryResult {
  data: never[] | undefined[];
}

interface RequestResult {
  status: number;
  headers: Function;
  config: object;
  statusText: string;
  xhrStatus: string;
}

export interface EntityQueryResult {
  data: {
    hasMore: boolean;
    data: Array<{ name: string; obid: number }>;
  };
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
        meta: T;
        tables?: null;
        series?: [];
      };
    };
  };
}

export interface TSDBArg extends RequestResult {
  data: {
    results: {
      [x: string]: {
        refId: string;
        series: Array<{ name: string; points: Array<[number, number]> }>;
        tables?: null;
      };
    };
  };
}

export interface TSDBResult extends RequestResult {
  data: TargetDatapoints[];
}
