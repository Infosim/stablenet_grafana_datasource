/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
const BACKEND_URL = '/api/tsdb/query';

export class StableNetDatasource {
  id: number;

  constructor(instanceSettings, $q, private backendSrv) {
    this.id = instanceSettings.id;
  }

  testDatasource(): Promise<{ status: string; message: string; title: string }> {
    const options: {
      headers: object;
      url: string;
      method: string;
      data: { queries: Array<{ datasourceId: number; queryType: string }> };
    } = {
      headers: { 'Content-Type': 'application/json' },
      url: BACKEND_URL,
      method: 'POST',
      data: {
        queries: [
          {
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

  queryDevices(queryString: string, refid: string): Promise<{ data: Array<{ text: string; value: number }>; hasMore: boolean }> {
    const data: {
      queries: Array<{
        refId: string;
        datasourceId: number;
        queryType: string;
        filter: string;
      }>;
    } = {
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
      const res: Array<{ text: string; value: number }> = result.data.results[refid].meta.data.map(device => {
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

  findMeasurementsForDevice(obid: number, input: string, refid: string): Promise<{ data: Array<{ text: string; value: number }>; hasMore: boolean }> {
    const data: {
      queries: Array<{
        refId: string;
        datasourceId: number;
        queryType: string;
        deviceObid: number;
        filter: string;
      }>;
    } = {
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
      const res: Array<{ text: string; value: number }> = result.data.results[refid].meta.data.map(measurement => {
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

  findMetricsForMeasurement(obid: number, refid: string): Promise<Array<{ text: string; value: string; measurementObid: number }>> {
    const data: { queries: Array<{ refId: string; datasourceId: number; queryType: string; measurementObid: number }> } = {
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

  async query(
    options: any
  ): Promise<{
    data: Array<{ target: string; datapoints: Array<[number, number]> }> | never[];
    status?: number;
    headers?: any;
    config?: any;
    statusText?: string;
    xhrStatus?: string;
  }> {
    const from: string = new Date(options.range.from).getTime().toString();
    const to: string = new Date(options.range.to).getTime().toString();
    const queries: Array<{
      refId: string;
      datasourceId: number;
      queryType: string;
      statisticLink?: string;
      requestData?: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }>;
      includeMinStats: boolean;
      includeAvgStats: boolean;
      includeMaxStats: boolean;
    }> = [];

    for (let i = 0; i < options.targets.length; i++) {
      const target: {
        refId: string;
        mode?: number;
        deviceQuery?: string;
        selectedDevice?: number;
        measurementQuery?: string;
        selectedMeasurement?: number;
        chosenMetrics?: object;
        metricPrefix?: string;
        includeMinStats?: boolean;
        includeAvgStats?: boolean;
        includeMaxStats?: boolean;
        statisticLink?: string;
        metrics?: Array<{ text: string; value: string; measurementObid: number; $$hashKey: string }>;
        moreDevices?: boolean;
        moreMeasurements?: boolean;
        datasource: any;
      } = options.targets[i];

      if (target.mode === 10 && target.statisticLink !== '') {
        queries.push({
          refId: target.refId,
          datasourceId: this.id,
          queryType: 'statisticLink',
          statisticLink: target.statisticLink,
          //@ts-ignore
          includeMinStats: target.includeMinStats,
          //@ts-ignore
          includeAvgStats: target.includeAvgStats,
          //@ts-ignore
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
          // @ts-ignore
          const text: string = target.metricPrefix + ' {MinMaxAvg} ' + target.metrics.filter(m => m.value === key)[0].text;
          keys.push({
            key: key,
            name: text,
          });
        }
      }

      requestData.push({
        // @ts-ignore
        measurementObid: target.selectedMeasurement,
        metrics: keys,
      });

      queries.push({
        refId: target.refId,
        datasourceId: this.id,
        queryType: 'metricData',
        requestData: requestData,
        //@ts-ignore
        includeMinStats: target.includeMinStats,
        //@ts-ignore
        includeAvgStats: target.includeAvgStats,
        //@ts-ignore
        includeMaxStats: target.includeMaxStats,
      });
    }

    if (queries.length === 0) {
      return { data: [] };
    }

    const data: {
      from: string;
      to: string;
      queries: Array<{
        refId: string;
        datasourceId: number;
        queryType: string;
        statisticLink?: string;
        requestData?: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }>;
        includeMinStats: boolean;
        includeAvgStats: boolean;
        includeMaxStats: boolean;
      }>;
    } = {
      from: from,
      to: to,
      queries: queries,
    };
    return await this.doRequest(data).then(handleTsdbResponse);
  }

  private doRequest(data: {
    from?: string;
    to?: string;
    queries: any[];
  }): Promise<{ data: { results: object }; status: number; headers?: any; config: any; statusText: string; xhrStatus: string }> {
    const options = {
      headers: { 'Content-Type': 'application/json' },
      url: BACKEND_URL,
      method: 'POST',
      data: data,
    };
    return this.backendSrv.datasourceRequest(options);
  }
}

export function handleTsdbResponse(response: {
  data: any;
  status: number;
  headers?: any;
  config: any;
  statusText: string;
  xhrStatus: string;
}): {
  data: Array<{ target: string; datapoints: Array<[number, number]> }>;
  status: number;
  headers?: any;
  config: any;
  statusText: string;
  xhrStatus: string;
} {
  const res: any = [];
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

  response.data = res;
  return response;
}
