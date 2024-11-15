/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2021
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { LabelValue, MetricResult, QueryResult, StableNetConfigOptions, Target } from './Types';

interface CollectionDTO<T> {
  hasMore: boolean;
  data: T[];
}

interface Device {
  obid: number;
  name: string;
}

interface Measurement {
  obid: number;
  name: string;
}

interface Metric {
  obid: number;
  key: string;
  name: string;
}

export class DataSource extends DataSourceWithBackend<Target, StableNetConfigOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<StableNetConfigOptions>) {
    super(instanceSettings);
  }

  async queryDevices(queryString: string): Promise<QueryResult> {
    const { data, hasMore }: CollectionDTO<Device> = await super.getResource('devices', { filter: queryString });

    const res: LabelValue[] = data.map(({ name, obid }) => ({ label: name, value: obid, }));

    res.push({ label: 'none', value: -1 });
    return { hasMore, data: res };
  }

  async findMeasurementsForDevice(obid: number, input: string): Promise<QueryResult> {
    const { data, hasMore }: CollectionDTO<Measurement> = await super.getResource('measurements', { deviceObid: obid, filter: input });

    return { hasMore, data: data.map(({ obid, name }) => ({ value: obid, label: name, })) };
  }

  async findMetricsForMeasurement(obid: number): Promise<MetricResult[]> {
    const result: Metric[] = await super.getResource('metrics', { measurementObid: obid });

    return result.map(({ obid, key, name }) => ({ measurementObid: obid, key, text: name, }));
  }
}
