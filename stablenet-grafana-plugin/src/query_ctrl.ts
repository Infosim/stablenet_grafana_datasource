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
        this.target.chosenMetrics = this.target.chosenMetrics || {};
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
        this.target.statisticLink = this.target.statisticLink || '';

        this.target.metrics = this.target.metrics || [];
            //normally metrics should not be stored within this.target (they can be fetched any time given measurement obid), 
            //but we need the variable to make ng-repeat in query-editor.html (and thus the checkboxes) work
        this.target.moreDevices = typeof this.target.moreDevices === 'undefined' ? false : this.target.moreDevices;
        this.target.moreMeasurements = typeof this.target.moreMeasurements === 'undefined' ? false : this.target.moreMeasurements;
            //these two do not belong in this.target either, but the ng-ifs in the optional tooltips have to be bound to something
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
        let result = this.datasource.queryDevices(this.target.deviceQuery, this.target.refId);
        result.then(r => this.target.moreDevices = r.hasMore)
        return result.then(r => r.data);
    }

    onDeviceChange() {
        this.target.measurementQuery = '';
        this.target.selectedMeasurement = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
    }

    getMeasurements() {
        let result = this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId);
        result.then(r => this.target.moreMeasurements = r.hasMore)
        return result.then(r => r.data);
    }

    onMeasurementRegexChange() {
        this.target.selectedMeasurement = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    onMeasurementChange() {
        this.datasource.findMetricsForMeasurement(this.target.selectedMeasurement, this.target.refId).then(res => this.target.metrics = res);
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    onChangeInternal() {
        this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

