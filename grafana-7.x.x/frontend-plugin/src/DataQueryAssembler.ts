/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { Target } from './QueryInterfaces';
import { Mode, SingleQuery, StringPair } from './Types';

export class WrappedTarget {
  constructor(private target: Target, private intervalMs: number) {}

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
      queryType: 'statisticLink',
      statisticLink: this.target.statisticLink,
      intervalMs: this.target.useCustomAverage
        ? parseInt(this.target.averagePeriod, 10) * this.target.averageUnit
        : this.intervalMs,
      includeMinStats: this.target.includeMinStats,
      includeAvgStats: this.target.includeAvgStats,
      includeMaxStats: this.target.includeMaxStats,
    };
  }

  toDeviceQuery(): SingleQuery {
    const keys: StringPair[] = this.getRequestedMetricsAsKeys();
    return {
      refId: this.target.refId,
      measurementObid: this.target.selectedMeasurement.value,
      metrics: keys,
      intervalMs: this.target.useCustomAverage
        ? parseInt(this.target.averagePeriod, 10) * this.target.averageUnit
        : this.intervalMs,
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
        const name: string = this.target.metricPrefix + ' ' + this.target.metrics.filter(m => m.key === key)[0].text;
        keys.push({
          key,
          name,
        });
      }
    }
    return keys;
  }
}
