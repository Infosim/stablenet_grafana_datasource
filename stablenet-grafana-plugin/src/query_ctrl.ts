/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import { QueryCtrl } from 'app/plugins/sdk';
import './css/query-editor.css!'

export class GenericDatasourceQueryCtrl extends QueryCtrl {
    constructor($scope, $injector) {
        super($scope, $injector);
        this.target.mode = this.target.mode || 'Device';
        this.target.deviceQuery = this.target.deviceQuery || '';
        this.target.selectedDevice = this.target.selectedDevice || 'select device';
        this.target.measurement = this.target.measurement || 'select measurement';
        this.target.metric = this.target.metric || 'select metric';
        this.target.statisticLink = this.target.statisticLink || '';
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
        this.target.deviceStorage = this.target.deviceStorage || 'select device';
        this.target.measurementStorage = this.target.measurementStorage || 'select measurement';
        this.target.metricStorage = this.target.metricStorage || 'select metric';
    }

    getModes() {
        return [{ text: 'Device', value: 'Device' }, { text: 'Statistic Link', value: 'Statistic Link' }];
    }

    onDeviceQueryChange() {
        this.target.selectedDevice = "select device";
        this.target.measurement = 'select measurement';
        this.target.metric = 'select metric';
        this.onChangeInternal();
    }

    getDevices() {
        return this.datasource.queryDevices(this.target.deviceQuery);
    }

    onDeviceChange() {
        this.target.measurement = 'select measurement';
        this.target.metric = 'select metric';
        this.onChangeInternal();
    }

    getMeasurements() {
        return this.datasource.findMeasurementsForDevice(this.target.selectedDevice);
    }

    onMeasurementChange() {
        this.target.metric = 'select metric';
        this.onChangeInternal();
    }

    getMetrics() {
        return this.datasource.findMetricsForMeasurement(this.target.measurement);
    }

    onChangeInternal() {
        this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

