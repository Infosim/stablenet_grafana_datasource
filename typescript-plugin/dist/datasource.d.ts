/// <reference path="../node_modules/grafana-sdk-mocks/app/headers/common.d.ts" />
export default class RocksetDatasource {
    private $q;
    private backendSrv;
    private templateSrv;
    headers: Object;
    id: number;
    name: string;
    url: string;
    apiKey: string;
    /** @ngInject */
    constructor(instanceSettings: any, $q: any, backendSrv: any, templateSrv: any);
    testDatasource(): any;
    queryDevices(queryString: any, refid: any): any;
    findMeasurementsForDevice(obid: any, input: any, refid: any): any;
    findMetricsForMeasurement(obid: any, refid: any): any;
    query(options: any): void;
    doRequest(data: any): any;
}
export declare function handleTsdbResponse(response: any): any;
