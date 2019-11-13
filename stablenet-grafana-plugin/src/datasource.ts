/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
import _ from "lodash";

const PROXY = 'http://localhost:3001';
const TEST = PROXY + '/';
const DEVICES_LIST = PROXY + '/devicesList';
const TAGS_LIST = PROXY + '/tagsList';
const DEVICES_MSM = PROXY + '/devicesMsm';
const MSM_LIST_GET = PROXY + '/msmListGet';
const MSM_LIST_POST = PROXY + '/msmListPost';
const SERVLET = PROXY + '/servlet';

const ANNOT_QUERY = PROXY + '/annotations';

export class GenericDatasource {

    constructor(instanceSettings, $q, backendSrv, templateSrv) {
        this.id = instanceSettings.id;
        this.backendSrv = backendSrv;
        this.templateSrv = templateSrv;

        console.log(instanceSettings)
    }

    testDatasource() {
        return this.doRequest({
            url: TEST,
            method: 'GET'
        }).then(response => {
            if (response.status === 200) {
                return {status: "success", message: "Data source is working", title: "Success"};
            }
        });
    }

    queryAllDevices(filter) {
        let data = {
            queries: [
                {
                    refId: "A",
                    datasourceId: this.id,   // Required
                    queryType: "devices"
                }
            ]
        };
        return this.doRequest(data).then(result => {
            return result.data.results.A.meta.map(device => {
                return {text: device.name, value: device.obid};
            })
        });
    }

    findMeasurementsForDevice(filter, obid) {
        if (obid === "select option") {
            return []
        }
        let data = {
            queries: [
                {
                    refId: "A",
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

    findMetricsForMeasurement(filter, obid) {
        if (obid === "select measurement") {
            return []
        }
        const from = this.templateSrv.timeRange.from.valueOf().toString();
        const to = this.templateSrv.timeRange.to.valueOf().toString();
        let data = {
            from: from,
            to: to,
            queries: [
                {
                    refId: "A",
                    datasourceId: this.id,   // Requiredma
                    queryType: "metricNames",
                    measurementObid: parseInt(obid)
                }
            ]
        };
        return this.doRequest(data).then(result => {
            return result.data.results.A.meta.map(metric => {
                return {text: metric, value: metric};
            })
        });
    }

    query(options) {
        const from = options.range.from.valueOf().toString();
        const to = options.range.to.valueOf().toString();
        console.log("Hello World");
        console.log(options);
        let queries = [];
        let id = this. id;
        options.targets.forEach(function (target) {
            queries.push({
                refId: target.refId,
                datasourceId: id,
                queryType: "metricData",
                measurementObid: parseInt(target.measurement),
                metricName: target.metricName
            }) ;
        });
        let data = {
            from: from,
            to: to,
            queries: queries
        };
        return this.doRequest(data).then(handleTsdbResponse);
    }

    doRequest(data) {
        let options = {
            headers: {'Content-Type': 'application/json'},
            url: '/api/tsdb/query',
            method: 'POST',
            data: data
        }
        return this.backendSrv.datasourceRequest(options);
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
    console.log(res);
    return response;
}