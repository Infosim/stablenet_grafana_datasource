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
        this.target.metricPrefix = this.target.metricPrefix || '';
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
        this.target.statisticLink = this.target.statisticLink || '';
        //normally metrics should not be stored within this.target (they can be fetched any time given measurement obid), 
        //but we need the variable to make ng-repeat in query-editor.html (and thus the checkboxes) work
        this.target.metrics = this.target.metrics || [];
        //the following two do not belong in this.target either, but the ng-ifs in the optional tooltips have to be bound to something
        this.target.moreDevices = typeof this.target.moreDevices === 'undefined' ? false : this.target.moreDevices;
        this.target.moreMeasurements = typeof this.target.moreMeasurements === 'undefined' ? false : this.target.moreMeasurements;
    }

    getModes() {
        return [{text: 'Device', value: 'Device'}, {text: 'Statistic Link', value: 'Statistic Link'}];
    }

    onDeviceQueryChange() {
        this.target.selectedDevice = "none";
        this.target.measurementQuery = '';
        this.target.selectedMeasurement = '';
        this.target.metricPrefix = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    getDevices() {
        return this.datasource.queryDevices(this.target.deviceQuery, this.target.refId)
                                .then(r => {
                                    this.target.moreDevices = r.hasMore;
                                    return r.data;
                                });
    }

    onDeviceChange() {
        this.target.measurementQuery = '';
        this.target.selectedMeasurement = '';
        this.target.metricPrefix = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
    }

    getMeasurements() {
        return this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId)
                                .then(r => {
                                    this.target.moreMeasurements = r.hasMore;
                                    return r.data;
                                })
    }

    onMeasurementRegexChange() {
        this.target.selectedMeasurement = '';
        this.target.metrics = [];
        this.target.chosenMetrics = {};
        this.onChangeInternal();
    }

    onMeasurementChange() {
        this.datasource.findMetricsForMeasurement(this.target.selectedMeasurement, this.target.refId)
                        .then(res => this.target.metrics = res);
        this.target.chosenMetrics = {};
        //We need the name of the Measurement with this id, but it is not saved on our side except somewhere in the dropdown box
        //which is basically an ul with one li per measurement and a small hell to work with. 1.5h wasted.
        this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId)
                                .then(r => r.data)
                                .then(r => r.filter(m => m.value === this.target.selectedMeasurement)[0])
                                .then(r => this.target.metricPrefix = r.text)
        //this.target.metricPrefix = this.target.selectedMeasurement; //@TODO: use name, not obid
        this.onChangeInternal();
    }

    onChangeInternal() {
        this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

