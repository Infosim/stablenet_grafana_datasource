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
        this.type = instanceSettings.type;
        this.url = instanceSettings.url;
        this.name = instanceSettings.name;
        this.q = $q;
        this.backendSrv = backendSrv;
        this.templateSrv = templateSrv;
        this.withCredentials = instanceSettings.withCredentials;
        this.headers = {'Content-Type': 'application/json'};
        if (typeof instanceSettings.basicAuth === 'string' && instanceSettings.basicAuth.length > 0) {
            this.headers['Authorization'] = instanceSettings.basicAuth;
        }

        console.log("CONSTRUCTOR");
    }

    /**
     * Is called when "Save and test" is clicked in the plugin Config page
     */
    testDatasource() {
        //console.log("TEST_DS");

        return this.doRequest({
            url: TEST,
            method: 'GET'
        }).then(response => {
            //console.log("Test DS Response: ", response)
            if (response.status === 200) {
                return {status: "success", message: "Data source is working", title: "Success"};
            }
        });
    }

    queryAllDevices(server, filter) {
        let data = {
            queries: [
                {
                    refId: "A",
                    datasourceId: 8,   // Required
                    queryType: "devices"
                }
            ]
        };
        return this.doRequest({
            url: '/api/tsdb/query',
            data: data,
            method: 'POST'
        }).then(result => {
            return result.data.results.A.meta.map(device => {
                return {text: device.name, value: device.obid};
            })
        });
    }

    findMeasurementsForDevice(server, filter, obid) {
        if(obid === "select option"){
            return []
        }
        let data = {
            queries: [
                {
                    refId: "A",
                    datasourceId: 8,   // Required
                    queryType: "measurements",
                    deviceObid: obid
                }
            ]
        };
        return this.doRequest({
            url: '/api/tsdb/query',
            data: data,
            method: 'POST'
        }).then(result => {
            return result.data.results.A.meta.map(measurement => {
                return {text: measurement.name, value: measurement.obid};
            })
        });
    }

    findMetricsForMeasurement(server, filter, obid) {
        if(obid === "select measurement"){
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
                    datasourceId: 8,   // Requiredma
                    queryType: "metricNames",
                    measurementObid: parseInt(obid)
                }
            ]
        };
        return this.doRequest({
            url: '/api/tsdb/query',
            data: data,
            method: 'POST'
        }).then(result => {
            return result.data.results.A.meta.map(metric => {
                return {text: metric, value: metric};
            })
        });
    }

    /**
     * Is called when the 5th dropdown menu is OPENED
     * and
     * after "query()" when a metric is PICKED
     */
    metricFindQueryq(server, filter, dot, query) {
        //console.log("METRIC_FIND_QUERY");
        console.log("MFQ Options: ", query);

        if (server == 'select server' || query == 'select measurement') {
            return [{text: 'select metric', value: 'select metric'}]
        }

        return this.doRequest({
            url: PROXY + '/search',
            data: {server: server, filter: filter, dot: dot, obid: query}, //interpolated,
            method: 'POST',
        }).then(result => {
            console.log("MFQ Response: ", result)
            return this.mapToTextValue(result)
        });
    }



    async metricFindQuery(server, filter, dot, query) {
        //console.log("METRIC_FIND_QUERY");
        console.log("MFQ Options: ", dot, query);

        if (server == 'select server' || query == 'select measurement') {
            return [{text: 'select metric', value: 'select metric'}]
        }

        let ids = typeof query === 'number' ?
            [query]
            :
            await this.doRequest({
                url: (filter === 'tag') ? MSM_LIST_POST : DEVICES_MSM,
                data: {server: server, obid: dot},
                method: 'POST',
            })
                .then(response => response.data)
                .then(xml => measurementListToJS(xml))                                                      //returns array of {text:.. , value:..}
                .then(m => m.filter(x => isRegex(query) ?
                    x.text.match(new RegExp(query.substring(1).slice(0, -1)))
                    :
                    x.text.indexOf(query) !== -1)
                    .map(x => x.value));

        //ids is [number] OR ( [number1, number2, ...] OR [] if nothing has left )

        let arr = [];

        for (let i = 0; i < ids.length; i++) {
            arr = await this.doRequest({
                url: SERVLET,
                data: {server: server, obid: ids[i]},
                method: 'POST',
            })
                .then(response => response.data)
                .then(d => toMetricResp(d))       //returns array of [metricname1, metricname2, ...]
                .then(stringArray => arr.concat(stringArray));        //in JS, .concat() returns a NEW array without changing the old ones
        }

        arr = [...new Set(arr)];

        arr = arr.sort();

        console.log("MFQ Response: ", arr);

        return await this.mapToTextValue({data: arr, status: 201, statusText: "Created"});      //mTTV expects an object with a 'data' property which is an array
    }

    doTsdbRequest() {
        let data = {
            from: "1555324640782",  // Optional, time range from
            to: "1555328240782",    // Optional, time range to
            queries: [
                {
                    datasourceId: 8,   // Required
                    refId: "A",         // Optional, default is "A"
                    maxDataPoints: 100, // Optional, default is 100
                    intervalMs: 1000,   // Optional, default is 1000

                    myFieldFoo: "bar",  // Any other fields,
                    myFieldBar: "baz",  // defined by user
                }
            ]
        };
        return this.backendSrv.datasourceRequest({
            url: '/api/tsdb/query',
            method: 'POST',
            data: data
        });
    }

    /**
     * No idea. Not going to delete tho.
     * @param {*} options
     */
    annotationQuery(options) {
        console.log("ANNOT_QUERY");

        var query = this.templateSrv.replace(options.annotation.query, {}, 'glob');
        var annotationQuery = {
            range: options.range,
            annotation: {
                name: options.annotation.name,
                datasource: options.annotation.datasource,
                enable: options.annotation.enable,
                iconColor: options.annotation.iconColor,
                query: query
            },
            rangeRaw: options.rangeRaw
        };

        return this.doRequest({
            url: ANNOT_QUERY,
            method: 'POST',
            data: annotationQuery
        }).then(function (result) {
            return result.data;
        });
    }

    doRequest(options) {
        options.withCredentials = this.withCredentials;
        options.headers = this.headers;

        return this.backendSrv.datasourceRequest(options);
    }

    /**
     * Transforms something into something else.            <----Quality commenting right here
     * Called right before query().
     * @param {*} options
     */
    buildQueryParameters(options) {
        //console.log("BUILD_QUERY_PARAMS", options)

        //remove placeholder targets
        options.targets = _.filter(options.targets, target => {
            return target.target !== 'select metric';
        });

        var targets = _.map(options.targets, target => {
            return {
                target: this.templateSrv.replace(target.target, options.scopedVars, 'regex'),
                refId: target.refId,
                data: target,           //the whole target,filter,measurement,... package
                hide: target.hide,
                type: target.type || 'timeserie'
            };
        });

        options.targets = targets;

        return options;
    }

    /**
     * Maps an array to an array of {text:... , value:...}
     * @param {*} result array (hopefully)
     * @returns {{text:string,value:any}[]} another array
     */
    mapToTextValue(result) {
        //console.log("MAP_TO_TV");
        return _.map(result.data, (d, i) => {
            if (d && d.text && d.value) {
                return {text: d.text, value: d.value};
            } else if (_.isObject(d)) {
                return {text: d, value: i};
            }
            return {text: d, value: d};
        });
    }

}

/**
 * Takes a xml from the /rest/devices/list or rest/tag-filters/list page and reduces it to
 * the names and the object ids.
 * @param {string} xml the xml in string form
 * @param regex either /<device/ or /<tagfilterelement/
 * @returns {{text:string, value: number}[]} an array of objects with properties 'text', 'value'
 */
function deviceListToJS(xml, regex) {
    return xml.replace(/.*<collection>/, '')                                                       //remove everything until end of opening collection-tag
        .split(regex)                                                                       //get all devices incl. tags
        .map(v => v.substring(0, v.indexOf('>')).trim())                                     //take infos until the one '>' that closes the device-tag, thus ignoring tags. Also trim!
        .filter(v => v !== "")                                                              //remove the first one...
        .map(v => v.split(/" /)                                                             //split into separate device properties
            .filter(w => w.indexOf('name') !== -1 || w.indexOf('obid') !== -1)      //take only the ones with name or obid
            .map(w => w.replace(/.*=/, '')))                                         //leave only the values, e.g., 'name="pizza' gives us "pizza
        .map(v => {
            return {
                text: v[0].substr(1),                                                            //substr because there is one " at the start and we removed the closing one with the split
                value: parseInt(v[1].substring(1).slice(0, -1))                                   //and this guy is packed in additional set of "" because of... reasons
            }
        })
}

/**
 * Takes a xml from /rest/devices/measurements/{id} and transforms it
 * @param {string} xml the xml in string form
 * @returns {{text: string, value: number}[]} an array of objects with text, value
 */
function measurementListToJS(xml) {
    let pure = xml.replace(/.*<collection>/, '');                        //remove everything until end of opening collection-tag
    let arr = [];
    while (pure.indexOf('obid') !== -1) {                                 //as long as we have obid-s...
        arr.push(pure.substring(1, pure.indexOf('>')))                     //add each measurement with the properties, i.e., until the end of the opening tag
        let name = pure.substring(1, pure.indexOf(' '));
        pure = pure.substring(pure.indexOf('</' + name + '>') + 3 + name.length);
    }

    return arr.map(v => v.split(/" /)                                                             //take the separate properties
        .map(w => " " + w)                                                      //to have sth like " name", for the regex not to mistake it with e.g. "templatename"
        .filter(w => w.indexOf(' name') != -1 || w.indexOf('obid') != -1)       //take only the obid and the name
        .map(w => w.replace(/.*=/, '')))                                         //reduce it to a string/string<number>
        .map(v => {
            return {
                text: v[0].substring(1),
                value: parseInt(v[1].substring(1, v[1].length - 1))
            }
        });
}

/**
 * Extracts the metrics present in the given Servlet
 * @param value a StatisticServlet as an array of objects
 * @returns an array with the metrics names
 */
function toMetricResp(value) {
    if (value.length === 0)
        return [];
    return Object.keys(value[0]).filter(prop => prop != "Time");
}

/**
 * Transforms StatisticServlet data to the format required by Grafana
 * @param target name of the metric we need
 * @param value an array with one object per measurement from the Servlet
 * @returns an array with numeric touples [value, unixEpochTime], e.g. [3.5, 1567411608]
 */
function toDatapoints(target, value) {
    return value.map(metric => {
        return {Time: new Date(metric['Time']), Opt: metric[target]}
    })
        .map(metric => [parseFloat(metric['Opt']), metric['Time'].getTime()])
}

/**
 * Checks if a string is a Regex, as per our definition in the datasource
 * @param text string to be checked
 * @returns 'true' if the string starts and ends with a '/', and 'false' otherwise
 */
function isRegex(text) {
    return text.charAt(0) === '/' && text.charAt(text.length - 1) === '/';
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