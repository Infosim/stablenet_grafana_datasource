/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { QueryCtrl } from 'grafana/app/plugins/sdk';
import './css/query-editor.css';
import { StableNetDatasource } from './datasource';
import { TextValue } from './returnTypes';
import { Mode, Unit } from './types';

/** @ngInject */
export class StableNetQueryCtrl extends QueryCtrl {
  static templateUrl = 'partials/query.editor.html';

  constructor($scope: any, $injector: any) {
    super($scope, $injector);
    this.target.mode = this.target.mode || Mode.MEASUREMENT;
    this.target.deviceQuery = this.target.deviceQuery || '';
    this.target.selectedDevice = this.target.selectedDevice || -1;
    this.target.measurementQuery = this.target.measurementQuery || '';
    this.target.selectedMeasurement = this.target.selectedMeasurement || '';
    this.target.chosenMetrics = this.target.chosenMetrics || {};
    this.target.metricPrefix = this.target.metricPrefix || '';
    this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
    this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
    this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
    this.target.statisticLink = this.target.statisticLink || '';
    this.target.averagePeriod = this.target.averagePeriod || '';
    this.target.averageUnit = this.target.averageUnit || Unit.MINUTES;
    this.target.useCustomAverage = typeof this.target.useCustomAverage === 'undefined' ? false : this.target.useCustomAverage;
    //normally metrics should not be stored within this.target (they can be fetched any time given measurement obid),
    //but we need the variable to make ng-repeat in query-editor.html (and thus the checkboxes) work
    this.target.metrics = this.target.metrics || [];
    //the following two do not belong in this.target either, but the ng-ifs in the optional tooltips have to be bound to something
    this.target.moreDevices = typeof this.target.moreDevices === 'undefined' ? false : this.target.moreDevices;
    this.target.moreMeasurements = typeof this.target.moreMeasurements === 'undefined' ? false : this.target.moreMeasurements;
  }

  getModes(): TextValue[] {
    return [
      { text: 'Measurement', value: Mode.MEASUREMENT },
      { text: 'Statistic Link', value: Mode.STATISTIC_LINK },
    ];
  }

  getUnits(): TextValue[] {
    return [
      { text: 'sec', value: Unit.SECONDS },
      { text: 'min', value: Unit.MINUTES },
      { text: 'hrs', value: Unit.HOURS },
      { text: 'days', value: Unit.DAYS },
    ];
  }

  onModeChange(): void {
    this.target.includeMinStats = false;
    this.target.includeAvgStats = true;
    this.target.includeMaxStats = false;
  }

  onDeviceQueryChange(): void {
    (this.datasource as StableNetDatasource)
      .queryDevices(this.target.deviceQuery, this.target.refId, this.$scope)
      .then(r => r.data)
      .then(r => (r ? r.map(el => el.value) : []))
      .then(r => {
        if (!r.includes(this.target.selectedDevice)) {
          this.target.selectedDevice = -1;
          this.target.measurementQuery = '';
          this.target.selectedMeasurement = '';
          this.target.metricPrefix = '';
          this.target.metrics = [];
          this.target.chosenMetrics = {};
        }
        return r;
      })
      .then(() => this.onChangeInternal());
  }

  getDevices(): Promise<TextValue[]> {
    return (this.datasource as StableNetDatasource).queryDevices(this.target.deviceQuery, this.target.refId, this.$scope).then(r => {
      this.target.moreDevices = r.hasMore;
      return r.data;
    });
  }

  onDeviceChange(): void {
    this.target.measurementQuery = '';
    this.target.selectedMeasurement = '';
    this.target.metricPrefix = '';
    this.target.metrics = [];
    this.target.chosenMetrics = {};
    setTimeout(() => {
      this.onChangeInternal();
    }, 50);
  }

  getMeasurements(): Promise<TextValue[]> {
    return (this.datasource as StableNetDatasource)
      .findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId, this.$scope)
      .then(r => {
        this.target.moreMeasurements = r.hasMore;
        return r.data;
      });
  }

  onMeasurementRegexChange(): void {
    (this.datasource as StableNetDatasource)
      .findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId, this.$scope)
      .then(r => r.data)
      .then(r => (r ? r.map(el => el.value) : []))
      .then(r => {
        if (!r.includes(this.target.selectedMeasurement)) {
          this.target.selectedMeasurement = '';
          this.target.metrics = [];
          this.target.chosenMetrics = {};
        }
        return r;
      })
      .then(() => this.onChangeInternal());
  }

  onMeasurementChange(): void {
    (this.datasource as StableNetDatasource)
      .findMetricsForMeasurement(this.target.selectedMeasurement, this.target.refId, this.$scope)
      .then(res => (this.target.metrics = res));
    this.target.chosenMetrics = {};
    (this.datasource as StableNetDatasource)
      .findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId, this.$scope)
      .then(r => r.data)
      .then(r => r.filter(m => m.value === this.target.selectedMeasurement)[0])
      .then(r => (this.target.metricPrefix = r.text));
    this.onChangeInternal();
  }

  onChangeInternal(): void {
    this.panelCtrl.refresh(); // Asks the panel to refresh data.
  }
}
