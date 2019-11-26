/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import _ from "lodash";

const BACKEND_URL = '/api/tsdb/query';
const DEFAULT_REFID = 'A';

export class GenericDatasource {

    constructor(instanceSettings, $q, backendSrv, templateSrv) {
        this.id = instanceSettings.id;
        this.backendSrv = backendSrv;
        this.templateSrv = templateSrv;
    }

    testDatasource() {
        let options = {
            headers: {'Content-Type': 'application/json'},
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
        }

        return this.backendSrv.request(options)
            .then(response => {
                if (response.message !== null) {
                    return {
                        status: "success",
                        message: "Data source is working and can connect to StableNet®.",
                        title: "Success"
                    };
                } else {
                    return {
                        status: "error",
                        message: "Datasource cannot connect to StableNet®.",
                        title: "Failure"
                    };
                }
            });
    }

    queryDevices(queryString) {
        let data = {
            queries: [
                {
                    refId: DEFAULT_REFID,
                    datasourceId: this.id,   // Required
                    queryType: "devices",
                    deviceQuery: queryString
                }
            ]
        };

        return this.doRequest(data)
            .then(result => {
                return result.data.results.A.meta.map(device => {
                    return {text: device.name, value: device.obid};
                })
            });
    }

    findMeasurementsForDevice(obid) {
        if (obid === "select device") {
            return [];
        }

        let data = {
            queries: [
                {
                    refId: DEFAULT_REFID,
                    datasourceId: this.id,   // Required
                    queryType: "measurements",
                    deviceObid: obid
                }
            ]
        };

        return this.doRequest(data).then(result => {
            return result.data.results.A.meta.map(measurement => {
                return {text: measurement.name, value: measurement.obid};
            })
        });
    }

    findMetricsForMeasurement(obid) {
        if (obid === "select measurement") {
            return [];
        }

        let data = {
            queries: []
        };
        if (typeof obid === 'number'){
            data.queries.push({
                refId: DEFAULT_REFID,
                datasourceId: this.id,
                queryType: "metricNames",
                measurementObid: obid
            })
        } else {
            //@TODO: find a way to ask (POST) Backend for /rest/devices/measurements : deviceId
        }

        return this.doRequest(data).then(result => {
            return result.data.results.A.meta.map(metric => {
                return {text: metric.name, value: metric.id};
            })
        });
    }

    async query(options) {
        const from = new Date(options.range.from).getTime().toString();
        const to = new Date(options.range.to).getTime().toString();
        let queries = [];
        let id = this.id;

        for(let i = 0; i < options.targets.length; i++){
            let target = options.targets[i];

            if (target.mode === "Statistic Link" && target.statisticLink !== "") {
                queries.push({
                    refId: target.refId,
                    datasourceId: id,
                    queryType: "statisticLink",
                    statisticLink: target.statisticLink,
                    includeMinStats: target.includeMinStats,
                    includeAvgStats: target.includeAvgStats,
                    includeMaxStats: target.includeMaxStats
                });
                continue;
            }

            if (!target.metric || target.metric === "select metric") {
                continue;
            }

            //WIP
            let measurementsList;
            if (typeof target.measurement === 'number'){
                measurementsList = [target.measurement];
            } else {
                measurementsList = await this.findMeasurementsForDevice(target.device)
                                            .then(response => response.filter(m => this.checkIfRegex(target.measurement) ?
                                                                                m.text.match(new RegExp(target.measurement.substring(1).slice(0,-1), "i"))
                                                                                :
                                                                                m.text.toLocaleLowerCase().indexOf(target.measurement.toLocaleLowerCase()) !== -1)
                                                                      .map(m => m.value));
            }
            //WIP end

            let metricsList;
            if (typeof target.metric === 'number'){
                metricsList =  [target.metric];
            } else {
                metricsList = await this.findMetricsForMeasurement(target.measurement)
                                        .then(response => response.filter(m => this.checkIfRegex(target.metric) ?
                                                                                m.text.match(new RegExp(target.metric.substring(1).slice(0,-1), "i"))
                                                                                :
                                                                                m.text.toLocaleLowerCase().indexOf(target.metric.toLocaleLowerCase()) !== -1)
                                                                  .map(m => m.value));
            }

            queries.push({
                refId: target.refId,
                datasourceId: id,
                queryType: "metricData",
                requestData: [{measurementObid: parseInt(target.measurement), metricIds: [...metricsList]}],
                includeMinStats: target.includeMinStats,
                includeAvgStats: target.includeAvgStats,
                includeMaxStats: target.includeMaxStats
            });
        }

        if (queries.length === 0) {
            return {data: []};
        }

        let data = {
            from: from,
            to: to,
            queries: queries
        };
        return await this.doRequest(data)
            .then(handleTsdbResponse);
    }

    doRequest(data) {
        let options = {
            headers: {'Content-Type': 'application/json'},
            url: BACKEND_URL,
            method: 'POST',
            data: data
        }
        return this.backendSrv.datasourceRequest(options);
    }

    checkIfRegex(text){
        return text.charAt(0) === '/' && text.charAt(text.length-1) === '/';
    }
}

export function handleTsdbResponse(response) {
    const res = [];
    _.forEach(response.data.results, r => {
        _.forEach(r.series, s => {
            res.push({target: s.name, datapoints: s.points});
        });
        _.forEach(r.tables, t => {
            t.type = 'table';
            t.refId = r.refId;
            res.push(t);
        });
    });

    response.data = res;
    return response;
}