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
    const options = {
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

    return this.backendSrv.request(options).then(response => {
      if (response.message !== null) {
        return {
          status: 'success',
          message: 'Data source is working and can connect to StableNet®.',
          title: 'Success',
        };
      } else {
        return {
          status: 'error',
          message: 'Datasource cannot connect to StableNet®.',
          title: 'Failure',
        };
      }
    });
  }

  queryDevices(queryString, refid): Promise<{ data: Array<{ text: string; value: number }>; hasMore: boolean }> {
    const data = {
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
      const res = result.data.results[refid].meta.data.map(device => {
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

  findMeasurementsForDevice(obid, input, refid): Promise<{ data: Array<{ text: string; value: number }>; hasMore: boolean }> | Promise<never[]> {
    if (obid === 'none') {
      return Promise.resolve([]);
    }

    const data: any = { queries: [] };

    if (input === undefined) {
      data.queries.push({
        refId: refid,
        datasourceId: this.id, // Required
        queryType: 'measurements',
        deviceObid: obid,
      });
    } else {
      data.queries.push({
        refId: refid,
        datasourceId: this.id, // Required
        queryType: 'measurements',
        deviceObid: obid,
        filter: input,
      });
    }

    return this.doRequest(data).then(result => {
      const res = result.data.results[refid].meta.data.map(measurement => {
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

  findMetricsForMeasurement(obid, refid): Promise<Array<{ text: string; value: string; measurementObid: number }>> | Promise<never[]> {
    if (obid === -1) {
      return Promise.resolve([]);
    }

    const data: any = {
      queries: [],
    };

    data.queries.push({
      refId: refid,
      datasourceId: this.id,
      queryType: 'metricNames',
      measurementObid: obid,
    });

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

  async query(options) {
    const from = new Date(options.range.from).getTime().toString();
    const to = new Date(options.range.to).getTime().toString();
    const queries: any = [];

    for (let i = 0; i < options.targets.length; i++) {
      const target = options.targets[i];

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

      const requestData: any = [];
      const keys: any = [];
      const e = Object.entries(target.chosenMetrics);

      for (const [key, value] of e) {
        if (value) {
          const text = target.metricPrefix + ' {MinMaxAvg} ' + target.metrics.filter(m => m.value === key)[0].text;
          keys.push({
            key: key,
            name: text,
          });
        }
      }

      requestData.push({
        measurementObid: parseInt(target.selectedMeasurement, 10),
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

    const data = {
      from: from,
      to: to,
      queries: queries,
    };
    return await this.doRequest(data).then(handleTsdbResponse);
  }

  doRequest(data): any {
    const options = {
      headers: { 'Content-Type': 'application/json' },
      url: BACKEND_URL,
      method: 'POST',
      data: data,
    };
    return this.backendSrv.datasourceRequest(options);
  }
}

export function handleTsdbResponse(response) {
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
