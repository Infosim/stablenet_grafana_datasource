/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { Target } from './query_interfaces';
import { Mode, SingleQuery, StringPair } from './types';

export class WrappedTarget {
  target: Target;
  intervalMs: number;
  dataSourceId: number;

  constructor(target: Target, intervalMs: number, dataSourceId: number) {
    this.target = target;
    this.intervalMs = intervalMs;
    this.dataSourceId = dataSourceId;
  }

  isValidStatisticLinkMode(): boolean {
    return this.target.mode === Mode.STATISTIC_LINK && this.target.statisticLink !== '';
  }

  hasEmptyMetrics(): boolean {
    return (
      !this.target.chosenMetrics ||
      Object.entries(this.target.chosenMetrics).length === 0 ||
      Object.values(this.target.chosenMetrics).filter(v => v).length === 0
    );
  }

  toStatisticLinkQuery(): SingleQuery {
    return {
      refId: this.target.refId,
      datasourceId: this.dataSourceId,
      queryType: 'statisticLink',
      statisticLink: this.target.statisticLink,
      intervalMs: this.target.useCustomAverage ? parseInt(this.target.averagePeriod, 10)* this.target.averageUnit : this.intervalMs,
      includeMinStats: this.target.includeMinStats,
      includeAvgStats: this.target.includeAvgStats,
      includeMaxStats: this.target.includeMaxStats,
    };
  }

  toDeviceQuery(): SingleQuery {
    const keys: StringPair[] = this.getRequestedMetricsAsKeys();
    const requestData: Array<{ measurementObid: number; metrics: Array<{ key: string; name: string }> }> = [];
    requestData.push({
      measurementObid: this.target.selectedMeasurement,
      metrics: keys,
    });

    return {
      refId: this.target.refId,
      datasourceId: this.dataSourceId,
      queryType: 'metricData',
      requestData: requestData,
      intervalMs: this.target.useCustomAverage ? parseInt(this.target.averagePeriod, 10) * this.target.averageUnit : this.intervalMs,
      includeMinStats: this.target.includeMinStats,
      includeAvgStats: this.target.includeAvgStats,
      includeMaxStats: this.target.includeMaxStats,
    };
  }

  private getRequestedMetricsAsKeys(): StringPair[] {
    const keys: StringPair[] = [];
    const e: Array<[string, boolean]> = Object.entries(this.target.chosenMetrics);

    for (const [key, value] of e) {
      if (value) {
        const name: string = this.target.metricPrefix + ' {MinMaxAvg} ' + this.target.metrics.filter(m => m.key === key)[0].text;
        keys.push({
          key,
          name,
        });
      }
    }
    return keys;
  }
}
