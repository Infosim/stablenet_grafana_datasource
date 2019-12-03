/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import {QueryCtrl} from 'app/plugins/sdk';
import './css/query-editor.css!'

export class GenericDatasourceQueryCtrl extends QueryCtrl {
    constructor($scope, $injector) {
        super($scope, $injector);
        this.target.mode = this.target.mode || 'Device';
        this.target.deviceQuery = this.target.deviceQuery || '';
        this.target.selectedDevice = this.target.selectedDevice || 'none';
        this.target.measurementQuery = this.target.measurementQuery || '';
        this.target.selectedMeasurement = this.target.selectedMeasurement || '';
        this.target.metrics = this.target.metrics || [];
        this.target.chosenMetrics = this.target.chosenMetrics || {};
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true  : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
        this.target.statisticLink = this.target.statisticLink || '';
    }

    getModes() {
        return [{text: 'Device', value: 'Device'}, {text: 'Statistic Link', value: 'Statistic Link'}];
    }

    onDeviceQueryChange() {
        this.target.selectedDevice = "none";
        this.target.measurementQuery = '';
        this.target.selectedMeasurement = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    getDevices() {
        return this.datasource.queryDevices(this.target.deviceQuery, this.target.refId);
    }

    onDeviceChange() {
        this.target.measurementQuery = '';
        this.target.selectedMeasurement = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
    }

    getMeasurements() {
        return this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId);
    }

    onMeasurementRegexChange() {
        this.target.selectedMeasurement = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    async onMeasurementChange() {
        this.target.metrics = await this.datasource.findMetricsForMeasurement(this.target.selectedMeasurement, this.target.refId);
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    onChangeInternal() {
        this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

