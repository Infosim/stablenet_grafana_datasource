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
        this.target.selectedDevice = this.target.selectedDevice || 'select device';
        this.target.measurement = this.target.measurement || 'select measurement';
        this.target.metric = this.target.metric || 'select metric';
        this.target.statisticLink = this.target.statisticLink || '';
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
        this.target.metricRegex = this.target.metricRegex || '.*';
    }

    getModes() {
        return [{text: 'Device', value: 'Device'}, {text: 'Statistic Link', value: 'Statistic Link'}];
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
        this.target.metric = "";
        this.datasource.findMetricsForMeasurement(this.target.measurement, this.target.refId);
    }

    getMetrics() {
        return JSON.parse(localStorage.getItem(this.target.refId + "_metrics"));
    }

    onMetricChange() {
        let allMetrics = JSON.parse(localStorage.getItem(this.target.refId + "_metrics"));
        for (let i = 0; i < allMetrics.length; i++) {
            let metric = allMetrics[i];
            if (metric.value === this.target.metric) {
                this.target.metricRegex = metric.text;
            }
        }
        this.onChangeInternal();
    }

    onMetricRegexChange(){
        let metricsList = [];
        let allMetrics = JSON.parse(localStorage.getItem(this.target.refId + "_metrics"));
        let regex = new RegExp(this.target.metricRegex, "i");
        for (let metricIndex = 0; metricIndex < allMetrics.length; metricIndex++) {
            let metric = allMetrics[metricIndex];
            if (regex.exec(metric.text)) {
                metricsList.push(metric.value);
            }
        }
        this.target.metricIds = metricsList;
        this.onChangeInternal();
    }


    onChangeInternal() {
        this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

