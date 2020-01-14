System.register(['app/plugins/sdk', './css/query_editor.css!'], function(exports_1) {
    var __extends = (this && this.__extends) || function (d, b) {
        for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
    var sdk_1;
    var StableNetQueryCtrl;
    return {
        setters:[
            function (sdk_1_1) {
                sdk_1 = sdk_1_1;
            },
            function (_1) {}],
        execute: function() {
            StableNetQueryCtrl = (function (_super) {
                __extends(StableNetQueryCtrl, _super);
                function StableNetQueryCtrl($scope, $injector) {
                    _super.call(this, $scope, $injector);
                    this.target.mode = this.target.mode || 0;
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
                    //normally metrics should not be stored within this.target (they can be fetched any time given measurement obid),
                    //but we need the variable to make ng-repeat in query-editor.html (and thus the checkboxes) work
                    this.target.metrics = this.target.metrics || [];
                    //the following two do not belong in this.target either, but the ng-ifs in the optional tooltips have to be bound to something
                    this.target.moreDevices = typeof this.target.moreDevices === 'undefined' ? false : this.target.moreDevices;
                    this.target.moreMeasurements = typeof this.target.moreMeasurements === 'undefined' ? false : this.target.moreMeasurements;
                }
                StableNetQueryCtrl.prototype.getModes = function () {
                    return [{ text: 'Measurement', value: 0 }, { text: 'Statistic Link', value: 10 }];
                };
                StableNetQueryCtrl.prototype.onModeChange = function () {
                    this.target.includeMinStats = false;
                    this.target.includeAvgStats = true;
                    this.target.includeMaxStats = false;
                };
                StableNetQueryCtrl.prototype.onDeviceQueryChange = function () {
                    var _this = this;
                    this.datasource.queryDevices(this.target.deviceQuery, this.target.refId)
                        .then(function (r) { return r.data; })
                        .then(function (r) { return r ? r.map(function (el) { return el.value; }) : []; })
                        .then(function (r) {
                        if (!r.includes(_this.target.selectedDevice)) {
                            _this.target.selectedDevice = -1;
                            _this.target.measurementQuery = '';
                            _this.target.selectedMeasurement = '';
                            _this.target.metricPrefix = '';
                            _this.target.metrics = [];
                            _this.target.chosenMetrics = {};
                        }
                        return r;
                    })
                        .then(function () { return _this.onChangeInternal(); });
                };
                StableNetQueryCtrl.prototype.getDevices = function () {
                    var _this = this;
                    return this.datasource.queryDevices(this.target.deviceQuery, this.target.refId)
                        .then(function (r) {
                        _this.target.moreDevices = r.hasMore;
                        return r.data;
                    });
                };
                StableNetQueryCtrl.prototype.onDeviceChange = function () {
                    this.target.measurementQuery = '';
                    this.target.selectedMeasurement = '';
                    this.target.metricPrefix = '';
                    this.target.metrics = [];
                    this.target.chosenMetrics = {};
                };
                StableNetQueryCtrl.prototype.getMeasurements = function () {
                    var _this = this;
                    return this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId)
                        .then(function (r) {
                        _this.target.moreMeasurements = r.hasMore;
                        return r.data;
                    });
                };
                StableNetQueryCtrl.prototype.onMeasurementRegexChange = function () {
                    var _this = this;
                    this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId)
                        .then(function (r) { return r.data; })
                        .then(function (r) { return r ? r.map(function (el) { return el.value; }) : []; })
                        .then(function (r) {
                        if (!r.includes(_this.target.selectedMeasurement)) {
                            _this.target.selectedMeasurement = '';
                            _this.target.metrics = [];
                            _this.target.chosenMetrics = {};
                        }
                        return r;
                    })
                        .then(function () { return _this.onChangeInternal(); });
                };
                StableNetQueryCtrl.prototype.onMeasurementChange = function () {
                    var _this = this;
                    this.datasource.findMetricsForMeasurement(this.target.selectedMeasurement, this.target.refId)
                        .then(function (res) { return _this.target.metrics = res; });
                    this.target.chosenMetrics = {};
                    this.datasource.findMeasurementsForDevice(this.target.selectedDevice, this.target.measurementQuery, this.target.refId)
                        .then(function (r) { return r.data; })
                        .then(function (r) { return r.filter(function (m) { return m.value === _this.target.selectedMeasurement; })[0]; })
                        .then(function (r) { return _this.target.metricPrefix = r.text; });
                    this.onChangeInternal();
                };
                StableNetQueryCtrl.prototype.onChangeInternal = function () {
                    this.panelCtrl.refresh(); // Asks the panel to refresh data.
                };
                StableNetQueryCtrl.templateUrl = 'partials/query.editor.html';
                return StableNetQueryCtrl;
            })(sdk_1.QueryCtrl);
            exports_1("StableNetQueryCtrl", StableNetQueryCtrl);
        }
    }
});
//# sourceMappingURL=query_ctrl.js.map