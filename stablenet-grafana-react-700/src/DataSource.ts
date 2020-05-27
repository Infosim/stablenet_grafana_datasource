/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataQueryRequest, DataQueryResponse, DataSourceApi, DataSourceInstanceSettings } from '@grafana/data';

import {
  StableNetConfigOptions,
  DeviceQuery,
  Query,
  MeasurementQuery,
  MetricQuery,
  BasicQuery,
  TestOptions,
  SingleQuery,
} from './Types';
import { Target } from './QueryInterfaces';
import {
  EmptyQueryResult,
  EntityQueryResult,
  GenericResponse,
  LabelValue,
  MetricResult,
  MetricType,
  QueryResult,
  TargetDatapoints,
  TestResult,
  TSDBArg,
  TSDBResult,
} from './ReturnTypes';
import { WrappedTarget } from './DataQueryAssembler';

const BACKEND_URL = '/api/tsdb/query';

export class DataSource extends DataSourceApi<Target, StableNetConfigOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<StableNetConfigOptions>, $q, private backendSrv) {
    super(instanceSettings);
  }

  async testDatasource(): Promise<TestResult> {
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

  async queryDevices(queryString: string, refid: string): Promise<QueryResult> {
    const data: Query<DeviceQuery> = this.createDeviceQuery(queryString, refid);
    return this.doRequest<GenericResponse<EntityQueryResult>>(data).then(result => {
      const res: LabelValue[] = result.data.results[refid].meta.data.map(device => {
        return {
          label: device.name,
          value: device.obid,
        };
      });
      res.unshift({
        label: 'none',
        value: -1,
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

  async findMeasurementsForDevice(obid: number, refid: string): Promise<QueryResult> {
    const data: Query<MeasurementQuery> = this.createMeasurementQuery(obid, refid);
    return this.doRequest<GenericResponse<EntityQueryResult>>(data).then(result => {
      const res: LabelValue[] = result.data.results[refid].meta.data.map(measurement => {
        return {
          label: measurement.name,
          value: measurement.obid,
        };
      });
      return {
        data: res,
        hasMore: result.data.results[refid].meta.hasMore,
      };
    });
  }

  private createMeasurementQuery(deviceObid: number, refid: string): Query<MeasurementQuery> {
    const data: MeasurementQuery = {
      refId: refid,
      datasourceId: this.id,
      queryType: 'measurements',
      deviceObid: deviceObid,
      filter: '',
    };
    return {
      queries: [data],
    };
  }

  async findMetricsForMeasurement(obid: number, refid: string): Promise<MetricResult[]> {
    const data: Query<MetricQuery> = this.createMetricQuery(obid, refid);
    return this.doRequest<GenericResponse<MetricType[]>>(data).then(result =>
      result.data.results[refid].meta.map(metric => {
        const m: MetricResult = {
          measurementObid: obid,
          key: metric.key,
          text: metric.name,
        };
        return m;
      })
    );
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

  async query(options: DataQueryRequest<Target>): Promise<DataQueryResponse | EmptyQueryResult> {
    const { range } = options;
    const from = range!.from.valueOf().toString(10);
    const to = range!.to.valueOf().toString(10);

    const { targets } = options;
    const queries: SingleQuery[] = [];
    if (!('statisticLink' in options.targets[0]) && !('chosenMetrics' in options.targets[0])) {
      return { data: [] };
    }

    for (let i = 0; i < targets.length; i++) {
      const target: WrappedTarget = new WrappedTarget(targets[i], options.intervalMs!, this.id);

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
      return { data: [] };
    }

    const data: Query<SingleQuery> = {
      from: from,
      to: to,
      queries: queries,
    };

    return await this.doRequest<TSDBArg>(data).then(handleTsdbResponse);
  }

  private doRequest<RETURN>(data: Query<BasicQuery>): Promise<RETURN> {
    const options: TestOptions = {
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
