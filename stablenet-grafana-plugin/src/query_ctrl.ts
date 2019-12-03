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
        this.target.measurement = this.target.measurement || '';
        this.target.dataQueries = this.target.dataQueries || [];
        this.target.statisticLink = this.target.statisticLink || '';
        this.target.includeMinStats = typeof this.target.includeMinStats === 'undefined' ? false : this.target.includeMinStats;
        this.target.includeAvgStats = typeof this.target.includeAvgStats === 'undefined' ? true : this.target.includeAvgStats;
        this.target.includeMaxStats = typeof this.target.includeMaxStats === 'undefined' ? false : this.target.includeMaxStats;
            this.target.metricRegex = this.target.metricRegex || '.*';
        this.target.metric = this.target.metric || [];
        this.target.chosenMetrics = this.target.chosenMetrics || [];
            this.target.testlist = ["a", "b", "x"]
    }

    getModes() {
        return [{text: 'Device', value: 'Device'}, {text: 'Statistic Link', value: 'Statistic Link'}];
    }

    onDeviceQueryChange() {
        this.target.selectedDevice = "none";
        this.target.measurement = 'none';
        this.target.metric = 'select metric';
        this.onChangeInternal();
    }

    getDevices() {
        return this.datasource.queryDevices(this.target.deviceQuery);
    }

    onDeviceChange() {
        this.target.measurement = 'none';
        this.target.metric = [];
    }

    getMeasurements() {
        return this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId);
    }

    async onMeasurementChange() {
        let pizza = await this.datasource.findMetricsForMeasurement(this.target.measurement, this.target.refId);
        console.log(this.target.metric);
        pizza.forEach(m => this.target.metric.push(m));
        this.target.metric.push(this.target.metric.length)
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

