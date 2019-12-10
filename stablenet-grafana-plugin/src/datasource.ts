/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import _ from "lodash";

const BACKEND_URL = '/api/tsdb/query';

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

    queryDevices(queryString, refid) {
        let data = {
            queries: [
                {
                    refId: refid,
                    datasourceId: this.id,   // Required
                    queryType: "devices",
                    filter: queryString
                }
            ]
        };

        return this.doRequest(data)
            .then(result => {
                let res =  result.data.results[refid].meta.data.map(device => {
                    return {text: device.name, value: device.obid};
                });
                res.unshift({text: "none", value: -1});
                return {data: res, hasMore: result.data.results[refid].meta.hasMore};
            });
    }

    findMeasurementsForDevice(obid, input, refid) {
        if (obid === "none") {
            return Promise.resolve([]);
        }

        let data = {queries: []};

        if (input === undefined){
            data.queries.push({
                refId: refid,
                datasourceId: this.id,   // Required
                queryType: "measurements",
                deviceObid: obid,
            });
        } else {
            data.queries.push({
                refId: refid,
                datasourceId: this.id,   // Required
                queryType: "measurements",
                deviceObid: obid,
                filter: input,
            });
        }

        return this.doRequest(data).then(result => {
            let res = result.data.results[refid].meta.data.map(measurement => {
                return {text: measurement.name, value: measurement.obid};
            });
            return {data: res, hasMore: result.data.results[refid].meta.hasMore};
        });
    }

    findMetricsForMeasurement(obid, refid) {
        if (obid === -1) {
            return Promise.resolve([]);
        }

        let data = {
            queries: []
        };

        data.queries.push({
            refId: refid,
            datasourceId: this.id,
            queryType: "metricNames",
            measurementObid: obid
        })

        return this.doRequest(data).then(result => {
            return result.data.results[refid].meta.map(metric => {
                return {text: metric.name, value: metric.dataId, measurementObid: obid};
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

            if (!target.chosenMetrics || (Object.entries(target.chosenMetrics).length === 0) 
                                      || (Object.values(target.chosenMetrics).filter(v => v).length === 0)){
                continue;
            }

            let requestData = [];
            let ids = [];
            let e = Object.entries(target.chosenMetrics);
            
            for (let [key, value] of e){
                if (value){
                    ids.push(parseInt(key));
                }
            }

            requestData.push({measurementObid: parseInt(target.selectedMeasurement), metricIds: ids});

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
