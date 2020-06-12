/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataQueryRequest, DataQueryResponse, DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import {
  LabelValue,
  MetricResult,
  QueryResult,
  SingleQuery,
  StableNetConfigOptions,
  Target,
  TestResult,
} from './Types';
import { EMPTY } from 'rxjs';
import { Observable } from 'rxjs';
import { WrappedTarget } from './DataQueryAssembler';

export class DataSource extends DataSourceWithBackend<Target, StableNetConfigOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<StableNetConfigOptions>) {
    super(instanceSettings);
  }

  async testDatasource(): Promise<TestResult> {
    return super
      .testDatasource()
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
    return super.postResource('devices', { filter: queryString }).then(result => {
      const res: LabelValue[] = result.data.map(device => {
        return {
          label: device.name,
          value: device.obid,
        };
      });
      res.push({
        label: 'none',
        value: -1,
      });
      return { data: res, hasMore: result.hasMore };
    });
  }

  async findMeasurementsForDevice(obid: number): Promise<QueryResult> {
    return super.postResource('measurements', { deviceObid: obid }).then(result => {
      const res: LabelValue[] = result.data.map(measurement => {
        return {
          label: measurement.name,
          value: measurement.obid,
        };
      });
      return {
        data: res,
        hasMore: result.hasMore,
      };
    });
  }

  async findMetricsForMeasurement(obid: number): Promise<MetricResult[]> {
    return super.postResource('metrics', { measurementObid: obid }).then(result =>
      result.map(metric => {
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
}
