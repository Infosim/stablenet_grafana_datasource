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
        this.target.dataQueries = this.target.dataQueries || [];
        this.target.statisticLink = this.target.statisticLink || '';
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
        this.target.metricRegex = this.target.metricRegex || '.*';
        this.target.measurementRegex = this.target.measurementRegex || '.*';
        this.target.metric = '';
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
        this.target.metric = '';
        localStorage.setItem(this.target.refId + "_measurements", "[]");
        this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.refId);
    }

    getMeasurements() {
        return JSON.parse(localStorage.getItem(this.target.refId + "_measurements"));
    }

    onMeasurementChange() {
        let allMeasurements = JSON.parse(localStorage.getItem(this.target.refId + "_measurements"));
        for (let i = 0; i < allMeasurements.length; i++) {
            let measurement = allMeasurements[i];
            if (measurement.value === this.target.measurement) {
                this.target.measurementRegex = measurement.text;
            }
        }
        this.onMeasurementRegexChange();
    }

    onMeasurementRegexChange() {
        localStorage.setItem(this.target.refId + "_metrics", "[]");
        let allMeasurements = JSON.parse(localStorage.getItem(this.target.refId + "_measurements"));
        let filteredMeasurements = [];
        let regex = new RegExp(this.target.measurementRegex, "i");
        for (let measurementIndex = 0; measurementIndex < allMeasurements.length; measurementIndex++) {
            let measurement = allMeasurements[measurementIndex];
            if (regex.exec(measurement.text)) {
                this.datasource.findMetricsForMeasurement(measurement.value, this.target.refId);
            }
        }
    }

    getMetrics() {
        let union = {};
        let metrics = JSON.parse(localStorage.getItem(this.target.refId + "_metrics"));
        for (let i = 0; i < metrics.length; i++) {
            union[metrics[i].text] = true;
        }
        let result = [];
        for (let [key, value] of Object.entries(union)) {
            result.push({value: key, text: key})
        }
        console.log(result);
        return Promise.resolve(result);
    }

    onMetricChange() {
        this.target.metricRegex = this.target.metric;
        this.onMetricRegexChange();
    }

    onMetricRegexChange() {
        let dataQueries = {};
        let allMetrics = JSON.parse(localStorage.getItem(this.target.refId + "_metrics"));
        let regex = new RegExp(this.target.metricRegex, "i");
        for (let metricIndex = 0; metricIndex < allMetrics.length; metricIndex++) {
            let metric = allMetrics[metricIndex];
            if (regex.exec(metric.text)) {
                if (!dataQueries[metric.measurementObid]) {
                    dataQueries[metric.measurementObid] = [];
                }
                dataQueries[metric.measurementObid].push(metric.value)
            }
        }
        this.target.dataQueries = dataQueries;
        this.onChangeInternal();
    }


    onChangeInternal() {
        this.panelCtrl.refresh(); // Asks the panel to refresh data.
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

