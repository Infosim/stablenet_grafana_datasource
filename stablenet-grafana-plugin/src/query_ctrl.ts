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
        this.target.includeMinStats = this.target.includeMinStats || true;
        this.target.includeAvgStats = this.target.includeAvgStats || true;
        this.target.includeMaxStats = this.target.includeMaxStats || true;
    }

    getModes(){
        return [{value:'Device', text:'Device'}, {value:'Statistic Link', text:'Statistic Link'}];
    }

    onDeviceQueryChange() {
        this.target.devices = this.datasource.queryDevices(this.target.deviceQuery);
        this.target.selectedDevice = "select device";
    }

    getDevices() {
        return this.target.devices || [];
    }

    onDeviceChange() {
        this.target.measurements = this.datasource.findMeasurementsForDevice(this.target.selectedDevice);
        this.target.measurement = 'select measurement';
    }

    getMeasurements() {
        return this.target.measurements || [];
    }

    onMeasurementChange() {
        this.target.metrics = this.datasource.findMetricsForMeasurement(this.target.measurement) || [];
        this.target.metric = 'select metric';
    }

    getMetrics() {
        return this.target.metrics;
    }

    /**
     * Following bug:
     *
     * When using the migrated metricFindQuery(), once a metric is chosen, the native 'this.panelCtrl.refresh()' function
     * sets the value of the dropdown menu text (not the menu items!!) to something internal before this internal thing is
     * updated. Such an update happens once metricFindQuery() returns. Therefore the shown text is always one choice 'behind'
     * the current one, although datapoints are correctly represented.
     *
     * To tackle this, the refresh() function is called with a 0.5s delay, so that metricFindQuery() has time to terminate.
     * This solution is of course temporary, until an alternative is found.
     */
    onChangeInternal() {
        setTimeout(() => {
            this.panelCtrl.refresh(); // Asks the panel to refresh data.
        }, 500)
    }
}

GenericDatasourceQueryCtrl.templateUrl = 'partials/query.editor.html';

