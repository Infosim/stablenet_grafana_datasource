/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
System.register([], function(exports_1) {
    var BACKEND_URL, RocksetDatasource;
    function handleTsdbResponse(response) {
        var res = [];
        Object.values(response.data.results).forEach(function (r) {
            if (r.series) {
                r.series.forEach(function (s) {
                    res.push({
                        target: s.name,
                        datapoints: s.points
                    });
                });
            }
            if (r.tables) {
                r.tables.forEach(function (t) {
                    t.type = 'table';
                    t.refId = r.refId;
                    res.push(t);
                });
            }
        });
        response.data = res;
        return response;
    }
    exports_1("handleTsdbResponse", handleTsdbResponse);
    return {
        setters:[],
        execute: function() {
            ///<reference path="../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
            BACKEND_URL = '/api/tsdb/query';
            RocksetDatasource = (function () {
                /** @ngInject */
                function RocksetDatasource(instanceSettings, $q, backendSrv, templateSrv) {
                    this.$q = $q;
                    this.backendSrv = backendSrv;
                    this.templateSrv = templateSrv;
                    this.id = instanceSettings.id;
                }
                RocksetDatasource.prototype.testDatasource = function () {
                    var options = {
                        headers: { 'Content-Type': 'application/json' },
                        url: BACKEND_URL,
                        method: 'POST',
                        data: {
                            queries: [
                                {
                                    datasourceId: this.id,
                                    queryType: "testDatasource"
                                }
                            ]
                        }
                    };
                    return this.backendSrv.request(options)
                        .then(function (response) {
                        if (response.message !== null) {
                            return {
                                status: "success",
                                message: "Data source is working and can connect to StableNet®.",
                                title: "Success"
                            };
                        }
                        else {
                            return {
                                status: "error",
                                message: "Datasource cannot connect to StableNet®.",
                                title: "Failure"
                            };
                        }
                    });
                };
                RocksetDatasource.prototype.queryDevices = function (queryString, refid) {
                    var data = {
                        queries: [
                            {
                                refId: refid,
                                datasourceId: this.id,
                                queryType: "devices",
                                filter: queryString
                            }
                        ]
                    };
                    return this.doRequest(data)
                        .then(function (result) {
                        var res = result.data.results[refid].meta.data.map(function (device) {
                            return {
                                text: device.name,
                                value: device.obid
                            };
                        });
                        res.unshift({
                            text: "none",
                            value: -1
                        });
                        return { data: res,
                            hasMore: result.data.results[refid].meta.hasMore
                        };
                    });
                };
                RocksetDatasource.prototype.findMeasurementsForDevice = function (obid, input, refid) {
                    if (obid === "none") {
                        return Promise.resolve([]);
                    }
                    var data = { queries: [] };
                    if (input === undefined) {
                        data.queries.push({
                            refId: refid,
                            datasourceId: this.id,
                            queryType: "measurements",
                            deviceObid: obid,
                        });
                    }
                    else {
                        data.queries.push({
                            refId: refid,
                            datasourceId: this.id,
                            queryType: "measurements",
                            deviceObid: obid,
                            filter: input,
                        });
                    }
                    return this.doRequest(data)
                        .then(function (result) {
                        var res = result.data.results[refid].meta.data.map(function (measurement) {
                            return {
                                text: measurement.name,
                                value: measurement.obid
                            };
                        });
                        return {
                            data: res,
                            hasMore: result.data.results[refid].meta.hasMore
                        };
                    });
                };
                RocksetDatasource.prototype.findMetricsForMeasurement = function (obid, refid) {
                    if (obid === -1) {
                        return Promise.resolve([]);
                    }
                    var data = {
                        queries: []
                    };
                    data.queries.push({
                        refId: refid,
                        datasourceId: this.id,
                        queryType: "metricNames",
                        measurementObid: obid
                    });
                    return this.doRequest(data)
                        .then(function (result) {
                        return result.data.results[refid].meta.map(function (metric) {
                            return {
                                text: metric.name,
                                value: metric.key,
                                measurementObid: obid
                            };
                        });
                    });
                };
                RocksetDatasource.prototype.query = function (options) {
                    var from = new Date(options.range.from).getTime().toString();
                    var to = new Date(options.range.to).getTime().toString();
                    var queries = [];
                };
                RocksetDatasource.prototype.doRequest = function (data) {
                    var options = {
                        headers: { 'Content-Type': 'application/json' },
                        url: BACKEND_URL,
                        method: 'POST',
                        data: data
                    };
                    return this.backendSrv.datasourceRequest(options);
                };
                return RocksetDatasource;
            })();
            exports_1("default", RocksetDatasource);
        }
    }
});
//# sourceMappingURL=datasource.js.map