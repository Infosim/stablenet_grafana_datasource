/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { isQOE, QueryOptions, QueryOptionsEmpty, Target } from './queryInterfaces';
import {
  EntityQueryResult,
  GenericResponse,
  MetricResult,
  MetricType,
  QueryResult,
  EmptyQueryResult,
  TargetDatapoints,
  TestResult,
  TextValue,
  TSDBArg,
  TSDBResult,
} from './returnTypes';
import { DeviceQuery, MeasurementQuery, MetricQuery, Query, SingleQuery, TestOptions } from './types';

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

  queryDevices(queryString: string, refid: string): Promise<QueryResult> {
    const data: Query<DeviceQuery> = this.createDeviceQuery(queryString, refid);
    return this.doRequest<GenericResponse<EntityQueryResult>>(data).then(result => {
      const res: TextValue[] = result.data.results[refid].meta.data.map(device => {
        return {
          text: device.name,
          value: device.obid.toString(),
        };
      });
      res.unshift({
        text: 'none',
        value: '-1',
      });
      return { data: res, hasMore: result.data.results[refid].meta.hasMore };
    });
  }

  private createDeviceQuery(queryString: string, refid: string): Query<DeviceQuery> {
    const query: DeviceQuery = {
      filter: queryString,
      datasourceId: this.id,
      queryType: 'devices',
      refId: refid,
    };
    return {
      queries: [query],
    };
  }

  findMeasurementsForDevice(obid: number, input: string, refid: string): Promise<QueryResult> {
    const data: Query<MeasurementQuery> = this.createMeasurementQuery(obid, input, refid);
    return this.doRequest<GenericResponse<EntityQueryResult>>(data).then(result => {
      const res: TextValue[] = result.data.results[refid].meta.data.map(measurement => {
        return {
          text: measurement.name,
          value: measurement.obid.toString(),
        };
      });
      return {
        data: res,
        hasMore: result.data.results[refid].meta.hasMore,
      };
    });
  }

  private createMeasurementQuery(deviceObid: number, input: string, refid: string): Query<MeasurementQuery> {
    const data: MeasurementQuery = {
      refId: refid,
      datasourceId: this.id,
      queryType: 'measurements',
      deviceObid: deviceObid,
      filter: input,
    };
    return {
      queries: [data],
    };
  }

  findMetricsForMeasurement(obid: number, refid: string): Promise<MetricResult[]> {
    const data: Query<MetricQuery> = this.createMetricQuery(obid, refid);
    return this.doRequest<GenericResponse<MetricType[]>>(data).then(result => {
      return result.data.results[refid].meta.map(metric => {
        const m: MetricResult = {
          measurementObid: obid,
          value: metric.key,
          text: metric.name,
        };
        return m;
      });
    });
  }

  private createMetricQuery(mesurementObid: number, refid: string): Query<MetricQuery> {
    const data: MetricQuery = {
      refId: refid,
      datasourceId: this.id,
      queryType: 'metricNames',
      measurementObid: mesurementObid,
    };
    return {
      queries: [data],
    };
  }

  async query(options: QueryOptionsEmpty): Promise<EmptyQueryResult>;
  async query(options: QueryOptions): Promise<TSDBResult>;
  async query(options: QueryOptions | QueryOptionsEmpty): Promise<TSDBResult | EmptyQueryResult> {
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

    const data: Query<SingleQuery> = {
      from: from,
      to: to,
      queries: queries,
    };
    return await this.doRequest<TSDBArg>(data).then(handleTsdbResponse);
  }

  private doRequest<RETURN>(data: Query<any>): Promise<RETURN> {
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
