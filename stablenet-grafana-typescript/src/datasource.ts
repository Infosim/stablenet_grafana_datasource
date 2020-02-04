/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import {
  FindMeasurementOptions,
  FindMetricsOptions,
  QueryDeviceOptions,
  RequestArgQuery,
  RequestArgStandard,
  SingleQuery,
  TestOptions,
} from './types';
import { isQOE, QueryOptions, QueryOptionsEmpty, Target } from './queryInterfaces';
import { FindResult, MetricResult, QueryResultEmpty, TargetDatapoints, TestResult, TextValue, TSDBArg, TSDBResult } from './returnTypes';

const BACKEND_URL = '/api/tsdb/query';

export class StableNetDatasource {
  id: number;

  constructor(instanceSettings, $q, private backendSrv) {
    this.id = instanceSettings.id;
  }

  testDatasource(): Promise<TestResult> {
    const options: TestOptions = {
      headers: { 'Content-Type': 'application/json' },
      url: BACKEND_URL,
      method: 'POST',
      data: {
        queries: [
          {
            refId: 'UNUSED',
            datasourceId: this.id,
            queryType: 'testDatasource',
          },
        ],
      },
    };

    return this.backendSrv
      .request(options)
      .then(() => {
        return {
          status: 'success',
          message: 'Data source is working and can connect to StableNet®.',
          title: 'Success',
        };
      })
      .catch(err => {
        return {
          status: 'error',
          message: err.data.message,
          title: 'Failure',
        };
      });
  }

  queryDevices(queryString: string, refid: string): Promise<FindResult> {
    const data: { queries: QueryDeviceOptions[] } = {
      queries: [
        {
          refId: refid,
          datasourceId: this.id, // Required
          queryType: 'devices',
          filter: queryString,
        },
      ],
    };

    return this.doRequest(data).then(result => {
      const res: TextValue[] = result.data.results[refid].meta.data.map(device => {
        return {
          text: device.name,
          value: device.obid,
        };
      });
      res.unshift({
        text: 'none',
        value: -1,
      });
      return { data: res, hasMore: result.data.results[refid].meta.hasMore };
    });
  }

  findMeasurementsForDevice(obid: number, input: string, refid: string): Promise<FindResult> {
    const data: { queries: FindMeasurementOptions[] } = {
      queries: [
        {
          refId: refid,
          datasourceId: this.id, // Required
          queryType: 'measurements',
          deviceObid: obid,
          filter: input,
        },
      ],
    };

    return this.doRequest(data).then(result => {
      const res: TextValue[] = result.data.results[refid].meta.data.map(measurement => {
        return {
          text: measurement.name,
          value: measurement.obid,
        };
      });
      return {
        data: res,
        hasMore: result.data.results[refid].meta.hasMore,
      };
    });
  }

  findMetricsForMeasurement(obid: number, refid: string): Promise<MetricResult[]> {
    const data: { queries: FindMetricsOptions[] } = {
      queries: [
        {
          refId: refid,
          datasourceId: this.id,
          queryType: 'metricNames',
          measurementObid: obid,
        },
      ],
    };

    return this.doRequest(data).then(result => {
      return result.data.results[refid].meta.map(metric => {
        return {
          text: metric.name,
          value: metric.key,
          measurementObid: obid,
        };
      });
    });
  }

  async query(options: QueryOptionsEmpty): Promise<QueryResultEmpty>;
  async query(options: QueryOptions): Promise<TSDBResult>;
  async query(options: QueryOptions | QueryOptionsEmpty): Promise<TSDBResult | QueryResultEmpty> {
    const from: string = new Date(options.range.from).getTime().toString();
    const to: string = new Date(options.range.to).getTime().toString();
    const queries: SingleQuery[] = [];

    if (isQOE(options)) {
      return { data: [] };
    }

    for (let i = 0; i < options.targets.length; i++) {
      const target: Target = options.targets[i];

      if (target.mode === 10 && target.statisticLink !== '') {
        queries.push({
          refId: target.refId,
          datasourceId: this.id,
          queryType: 'statisticLink',
          statisticLink: target.statisticLink,
          includeMinStats: target.includeMinStats,
          includeAvgStats: target.includeAvgStats,
          includeMaxStats: target.includeMaxStats,
        });
        continue;
      }

      if (
        !target.chosenMetrics ||
        Object.entries(target.chosenMetrics).length === 0 ||
        Object.values(target.chosenMetrics).filter(v => v).length === 0
      ) {
        continue;
      }

      const keys: Array<{ key: string; name: string }> = [];
      const requestData: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }> = [];
      const e: Array<[string, boolean]> = Object.entries(target.chosenMetrics);

      for (const [key, value] of e) {
        if (value) {
          const text: string = target.metricPrefix + ' {MinMaxAvg} ' + target.metrics.filter(m => m.value === key)[0].text;
          keys.push({
            key: key,
            name: text,
          });
        }
      }

      requestData.push({
        measurementObid: target.selectedMeasurement,
        metrics: keys,
      });

      queries.push({
        refId: target.refId,
        datasourceId: this.id,
        queryType: 'metricData',
        requestData: requestData,
        includeMinStats: target.includeMinStats,
        includeAvgStats: target.includeAvgStats,
        includeMaxStats: target.includeMaxStats,
      });
    }

    if (queries.length === 0) {
      return { data: [] };
    }

    const data: RequestArgQuery = {
      from: from,
      to: to,
      queries: queries,
    };
    return await this.doRequest(data).then(handleTsdbResponse);
  }

  private doRequest(data: RequestArgStandard); //for the queryDevices(), findMeasurements() and findMetrics() functions
  private doRequest(data: RequestArgQuery): Promise<TSDBArg>;
  private doRequest(
    data: RequestArgQuery | RequestArgStandard
  ): Promise<{ data: { results: object }; status: number; headers: any; config: any; statusText: string; xhrStatus: string }> {
    const options = {
      headers: { 'Content-Type': 'application/json' },
      url: BACKEND_URL,
      method: 'POST',
      data: data,
    };
    return this.backendSrv.datasourceRequest(options);
  }
}

export function handleTsdbResponse(response: TSDBArg): TSDBResult {
  const res: TargetDatapoints[] = [];
  Object.values(response.data.results).forEach((r: any) => {
    if (r.series) {
      r.series.forEach(s => {
        res.push({
          target: s.name,
          datapoints: s.points,
        });
      });
    }
    if (r.tables) {
      r.tables.forEach(t => {
        t.type = 'table';
        t.refId = r.refId;
        res.push(t);
      });
    }
  });
  return {
    status: response.status,
    headers: response.headers,
    config: response.config,
    statusText: response.statusText,
    xhrStatus: response.xhrStatus,
    data: res,
  };
}
