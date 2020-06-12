/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataQueryRequest, DataQueryResponse, DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { SingleQuery, StableNetConfigOptions, TestOptions } from './Types';
import { EMPTY } from 'rxjs';
import { Target } from './QueryInterfaces';
import {
  EntityQueryResult,
  LabelValue,
  MetricQueryResult,
  MetricResult,
  QueryResult,
  TargetDatapoints,
  TestResult,
  TSDBArg,
  TSDBResult,
} from './ReturnTypes';
import { Observable } from 'rxjs';
import { WrappedTarget } from './DataQueryAssembler';

export class DataSource extends DataSourceWithBackend<Target, StableNetConfigOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<StableNetConfigOptions>, private backendSrv) {
    super(instanceSettings);
  }

  async testDatasource(): Promise<TestResult> {
    const options: TestOptions = {
      headers: { 'Content-Type': 'application/json' },
      url: '/api/datasources/' + this.id + '/resources/test',
      method: 'GET',
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

  async queryDevices(queryString: string): Promise<QueryResult> {
    return this.doResourceRequest<EntityQueryResult>('devices', { filter: queryString }).then(result => {
      const res: LabelValue[] = result.data.data.map(device => {
        return {
          label: device.name,
          value: device.obid,
        };
      });
      res.push({
        label: 'none',
        value: -1,
      });
      return { data: res, hasMore: result.data.hasMore };
    });
  }

  async findMeasurementsForDevice(obid: number): Promise<QueryResult> {
    const data = { deviceObid: obid, filter: '' };
    return this.doResourceRequest<EntityQueryResult>('measurements', data).then(result => {
      const res: LabelValue[] = result.data.data.map(measurement => {
        return {
          label: measurement.name,
          value: measurement.obid,
        };
      });
      return {
        data: res,
        hasMore: result.data.hasMore,
      };
    });
  }

  async findMetricsForMeasurement(obid: number): Promise<MetricResult[]> {
    const data = { measurementObid: obid };
    return this.doResourceRequest<MetricQueryResult>('metrics', data).then(result =>
      result.data.map(metric => {
        const m: MetricResult = {
          measurementObid: obid,
          key: metric.key,
          text: metric.name,
        };
        return m;
      })
    );
  }

  query(request: DataQueryRequest<Target>): Observable<DataQueryResponse> {
    const { targets } = request;
    const queries: SingleQuery[] = [];
    if (!('statisticLink' in request.targets[0]) && !('chosenMetrics' in request.targets[0])) {
      return EMPTY;
    }

    for (let i = 0; i < targets.length; i++) {
      const target: WrappedTarget = new WrappedTarget(targets[i], request.intervalMs!);

      if (target.isValidStatisticLinkMode()) {
        queries.push(target.toStatisticLinkQuery());
        continue;
      }

      if (target.hasEmptyMetrics()) {
        continue;
      }

      queries.push(target.toDeviceQuery());
    }

    if (queries.length === 0) {
      return EMPTY;
    }

    const req: DataQueryRequest = {
      ...request,
      targets: queries,
    };

    return super.query(req);
  }

  private doResourceRequest<RETURN>(resource: string, data: any): Promise<RETURN> {
    const options: TestOptions = {
      headers: { 'Content-Type': 'application/json' },
      url: '/api/datasources/' + this.id + '/resources/' + resource,
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
