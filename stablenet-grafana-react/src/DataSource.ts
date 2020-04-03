import defaults from 'lodash/defaults';

import { DataQueryRequest, DataQueryResponse, DataSourceApi, DataSourceInstanceSettings, MutableDataFrame, FieldType } from '@grafana/data';

import {defaultQuery, StableNetConfigOptions, MyQuery, TestOptions} from './types';
import {TestResult} from "./returnTypes";

const BACKEND_URL = '/api/tsdb/query';

export class DataSource extends DataSourceApi<MyQuery, StableNetConfigOptions> {

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

  async query(options: DataQueryRequest<MyQuery>): Promise<DataQueryResponse> {
    const { range } = options;
    const from = range!.from.valueOf();
    const to = range!.to.valueOf();

    // Return a constant for each query.
    const data = options.targets.map(target => {
      const query = defaults(target, defaultQuery);
      return new MutableDataFrame({
        refId: query.refId,
        fields: [
          { name: 'Time', values: [from, to], type: FieldType.time },
          { name: 'Value', values: [query.constant, query.constant], type: FieldType.number },
        ],
      });
    });

    return { data };
  }

}
