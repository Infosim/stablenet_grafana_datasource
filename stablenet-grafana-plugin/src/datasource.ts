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
                    filter: queryString
                }
            ]
        };

        return this.doRequest(data)
            .then(result => {
                return result.data.results.A.meta.data.map(device => {
                    return {text: device.name, value: device.obid};
                })
            });
    }

    findMeasurementsForDevice(obid, refid) {
        if (obid === "select device") {
            return Promise.resolve([]);
        }

        let data = {
            queries: [
                {
                    refId: refid,
                    datasourceId: this.id,   // Required
                    queryType: "measurements",
                    deviceObid: obid
                }
            ]
        };

        return this.doRequest(data).then(result => {
            return result.data.results.A.meta.data.map(measurement => {
                let loadedMeasurements = JSON.parse(localStorage.getItem(refid+ "_measurements"));
                let object = {text: measurement.name, value: measurement.obid};
                loadedMeasurements.push(object);
                localStorage.setItem(refid + "_measurements", JSON.stringify(loadedMeasurements));
            })
        });
    }

    findMetricsForMeasurement(obid, refid) {
        if (obid === "select measurement") {
            return Promise.resolve([]);
        }

        let data = {
            queries: []
        };
        if (typeof obid === 'number') {
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
                let loadedMetrics = JSON.parse(localStorage.getItem(refid + "_metrics"));
                let object = {text: metric.name, value: metric.key, measurementObid: obid};
                loadedMetrics.push(object);
                localStorage.setItem(refid + "_metrics", JSON.stringify(loadedMetrics));
            })
        });
    }

    async query(options) {
        const from = new Date(options.range.from).getTime().toString();
        const to = new Date(options.range.to).getTime().toString();
        let queries = [];

        for (let i = 0; i < options.targets.length; i++) {
            let target = options.targets[i];

            if (target.mode === "Statistic Link" && target.statisticLink !== "") {
                queries.push({
                    refId: target.refId,
                    datasourceId: this.id,
                    queryType: "statisticLink",
                    statisticLink: target.statisticLink,
                    includeMinStats: target.includeMinStats,
                    includeAvgStats: target.includeAvgStats,
                    includeMaxStats: target.includeMaxStats
                });
                continue;
            }

            if(!target.dataQueries){
                continue;
            }
            let requestData = [];
            for (const [measurementObid, metricIds] of Object.entries(target.dataQueries)){
                requestData.push({measurementObid: parseInt(measurementObid), metricIds: metricIds})
            }
            if (requestData.length == 0) {
                continue;
            }

            queries.push({
                refId: target.refId,
                datasourceId: this.id,
                queryType: "metricData",
                requestData: requestData,
                includeMinStats: target.includeMinStats,
                includeAvgStats: target.includeAvgStats,
                includeMaxStats: target.includeMaxStats
            });
        }

        if (queries.length === 0) {
            return { data: [] };
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
}

export function checkIfRegex(text) {
    return text.charAt(0) === '/' && text.charAt(text.length - 1) === '/';
}

export function filterTextValuePair(pair, filterValue) {
    return checkIfRegex(filterValue) ?
        pair.text.match(new RegExp(filterValue.substring(1).slice(0, -1), "i"))
        :
        pair.text.toLocaleLowerCase().indexOf(filterValue.toLocaleLowerCase()) !== -1;
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
