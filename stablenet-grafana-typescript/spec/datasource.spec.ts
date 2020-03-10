import {handleTsdbResponse, StableNetDatasource} from "../src/datasource";
import {MockBackendServer} from "./mock_server";
import {DeviceQuery, MeasurementQuery, MetricQuery, Query, TestOptions} from "../src/types";
import {MetricResult, QueryResult, TestResult, TSDBArg, TSDBResult} from "../src/returnTypes";
import {QueryOptions} from "../src/query_interfaces";

let backendSrv: MockBackendServer;
let datasource: StableNetDatasource;

describe("handleTsdbResponse()", () => {
    let response: TSDBArg;

    beforeEach(() => {
        response = {
            config: {},
            data: {
                results: {}
            },
            headers: () => {},
            status: 10,
            statusText: "success",
            xhrStatus: ""
        }
    });

    it('should return empty data on empty results', () => {
        expect(handleTsdbResponse(response).data).toEqual([]);
    });

    it('should return empty data on empty series', () => {
        response.data.results = {
            "A": {
                refId: "A",
                series: [],
                tables: null
            }
        };

        expect(handleTsdbResponse(response).data).toEqual([]);
    });

    it('should parse correctly data series', () => {
        response.data.results = {
            "A": {
                refId: "A",
                series: [
                    {
                        name: "CPU",
                        points: [[0,1], [0,2],[0,3]]
                    }
                ],
                tables: null,
            }
        };

        let expected: TSDBResult = {
            config: {},
            data: [
                {
                    target: "CPU",
                    datapoints: [[0,1], [0,2],[0,3]]
                }
            ],
            headers: response.headers,
            status: 10,
            statusText: "success",
            xhrStatus: ""
        };

        expect(handleTsdbResponse(response)).toEqual(expected);
    });

    it('should parse correctly multiple data series', () => {
        response.data.results = {
            "A": {
                refId: "A",
                series: [
                    {
                        name: "CPU",
                        points: [[0,1], [1,2],[2,3]]
                    },
                    {
                        name: "Ping",
                        points: [[0.5, 1.5], [1.5, 2.5], [2.5, 3.5]]
                    }
                ],
                tables: null,
            }
        };

        let expected: TSDBResult = {
            config: {},
            data: [
                {
                    target: "CPU",
                    datapoints: [[0,1], [1,2],[2,3]]
                },
                {
                    target: "Ping",
                    datapoints: [[0.5, 1.5], [1.5, 2.5], [2.5, 3.5]]
                }
            ],
            headers: response.headers,
            status: 10,
            statusText: "success",
            xhrStatus: ""
        };

        expect(handleTsdbResponse(response)).toEqual(expected);
    });
});

describe("constructor", () => {
    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);
    });

    it('should create a valid object', () => {
        expect(datasource).not.toBeNull();
        expect(datasource).not.toBeUndefined();
    });

    it('should assign correct id', () => {
        expect(datasource.id).toBe(1);
    });
});

describe("testDatasource()", () => {
    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);
        spyOn(backendSrv,'request').and.callThrough();
        datasource.testDatasource();
    });

    it('should call the backend server', () => {
        expect(backendSrv.request).toHaveBeenCalled();
    });

    it('should apply correct arguments', () => {
        let expected: TestOptions = {
            data: {
                queries: [
                    {
                        refId: 'UNUSED',
                        datasourceId: 1,
                        queryType: 'testDatasource',
                    },
                ],
            },
            headers: {"Content-Type": "application/json"},
            method: "POST",
            url: "/api/tsdb/query"
        };
        expect(backendSrv.request).toHaveBeenCalledWith(expected);
    });

    it('should return successfully on fulfilled Promise', async () => {
        let expected: TestResult = {
            status: 'success',
            message: 'Data source is working and can connect to StableNet®.',
            title: 'Success',
        };
        expect(await datasource.testDatasource()).toEqual(expected);
    });
});

describe("queryDevices()", () => {
    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);

        spyOn(backendSrv, 'forDeviceQuery').and.callThrough();
        spyOn(datasource, 'queryDevices').and.callThrough();
        datasource.queryDevices("","A");
    });

    it('should use correct type and call forDeviceQuery', () => {
        expect(backendSrv.forDeviceQuery).toHaveBeenCalled();
    });

    it('should make a call with correct argument', () => {
        let expected: Query<DeviceQuery> = {
            queries: [
                {
                    filter: "",
                    datasourceId: 1,
                    queryType: 'devices',
                    refId: "A",
                }
            ]
        };

        expect(backendSrv.forDeviceQuery).toHaveBeenCalledWith(expected);
    });

    it('should always have the empty id first', async () => {
        expect((await datasource.queryDevices("","A")).data[0]).toEqual({text: 'none', value: -1});
    });

    it('should return correctly', async () => {
        let expected: QueryResult = {
            data: [
                {text: 'none', value: -1},
                {text: 'ada', value: 1008},
                {text: 'adfsserver2.root.infosim.net', value: 1039},
                {text: 'agentl.infosim.net', value: 1068},
            ],
            hasMore: false
        };

        expect(await datasource.queryDevices("","A")).toEqual(expected);
    });
});

describe("findMeasurementsForDevice()", () => {
    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);

        spyOn(backendSrv, 'forMeasurementQuery').and.callThrough();
        spyOn(datasource, 'findMeasurementsForDevice').and.callThrough();
        datasource.findMeasurementsForDevice(1083, "","A");
    });

    it('should use correct type and call forMeasurementQuery', () => {
        expect(backendSrv.forMeasurementQuery).toHaveBeenCalled();
    });

    it('should make a call with correct argument', () => {
        let expected: Query<MeasurementQuery> = {
            queries: [
                {
                    refId: "A",
                    datasourceId: 1,
                    queryType: "measurements",
                    deviceObid: 1083,
                    filter: ""
                }
            ]
        };

        expect(backendSrv.forMeasurementQuery).toHaveBeenCalledWith(expected);
    });

    it('should return correctly', async () => {
        let expected: QueryResult = {
            hasMore: false,
            data: [
                {text: "Berlin", value: 3691},
                {text: "Berlin Cisco CPU 1", value: 3701},
                {text: "Berlin Cisco Dial-In Statistics", value: 3707}
            ]
        };

        expect(await datasource.findMeasurementsForDevice(1083,"","A")).toEqual(expected);
    });
});

describe("findMetricsForMeasurement()", () => {
    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);

        spyOn(backendSrv, 'forMetricQuery').and.callThrough();
        spyOn(datasource, 'findMetricsForMeasurement').and.callThrough();
        datasource.findMetricsForMeasurement(3701,"A");
    });

    it('should use correct type and call forMetricQuery', () => {
        expect(backendSrv.forMetricQuery).toHaveBeenCalled();
    });

    it('should make a call with correct argument', () => {
        let expected: Query<MetricQuery> = {
            queries: [
                {
                    refId: "A",
                    datasourceId: 1,
                    queryType: 'metricNames',
                    measurementObid: 3701,
                }
            ]
        };

        expect(backendSrv.forMetricQuery).toHaveBeenCalledWith(expected);
    });

    it('should return correctly', async () => {
        let expected: MetricResult[] = [
            {key: "SNMP_1000", text: "System Users", measurementObid: 3701},
            {key: "SNMP_1001", text: "System Processes", measurementObid: 3701},
            {key: "SNMP_1002", text: "System Uptime", measurementObid: 3701}
        ];

        expect(await datasource.findMetricsForMeasurement(3701, "A")).toEqual(expected);
    });
});

describe("query()", () => {
    let arg: QueryOptions;

    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);

        spyOn(backendSrv, 'forSingleQuery').and.callThrough();

        arg = {
            range: {
                from: 1583838568548,
                to: 1583858568549
            },
            intervalMs: 20000,
            targets: [
                {
                    refId: "A",
                    mode: 0,
                    deviceQuery: "",
                    selectedDevice: 1083,
                    measurementQuery: "",
                    selectedMeasurement: 3701,
                    chosenMetrics: {SNMP_1: true},
                    metricPrefix: "Berlin Cisco CPU 1",
                    includeMinStats: false,
                    includeAvgStats: true,
                    includeMaxStats: false,
                    statisticLink: "",
                    averagePeriod: "",
                    averageUnit: 60000,
                    useCustomAverage: false,
                    metrics: [
                        {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
                        {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}
                    ],
                    moreDevices: false,
                    moreMeasurements: false,
                    datasource: undefined
                }
            ]
        };
    });

    it('should use correct type and call forSingleQuery', () => {
        datasource.query(arg);
        expect(backendSrv.forSingleQuery).toHaveBeenCalled();
    });

    it('should produce proper args for one query', () => {
        let expected = {
            from: '1583838568548',
            to:   '1583858568549',
            queries: [ {
                refId: 'A',
                datasourceId: 1,
                queryType: 'metricData',
                requestData: [
                    { measurementObid: 3701,
                        metrics: [
                            {
                                key: 'SNMP_1',
                                name: 'Berlin Cisco CPU 1 {MinMaxAvg} cpu-load 1min'
                            }
                        ]
                    }
                ],
                intervalMs: 20000,
                includeMinStats: false,
                includeAvgStats: true,
                includeMaxStats: false
            } ]
        };

        datasource.query(arg);
        expect(backendSrv.forSingleQuery).toHaveBeenCalledWith(expected);
    });

    it('should work as expected for invalid options', async () => {
        delete arg.targets[0].mode;
        datasource.query(arg);
        expect(backendSrv.forSingleQuery).not.toHaveBeenCalled();
        expect(await datasource.query(arg)).toEqual({ data: []});
    });

    it('should work as expected when no metrics are chosen', async () => {
        arg.targets[0].chosenMetrics = {};
        datasource.query(arg);
        expect(backendSrv.forSingleQuery).not.toHaveBeenCalled();
        expect(await datasource.query(arg)).toEqual({ data: []});
    });

    it('should call backend in statistic link mode', () => {
        arg.targets[0].statisticLink = "?id=1234";
        arg.targets[0].mode = 10;
        arg.targets[0].chosenMetrics = {};
        datasource.query(arg);
        expect(backendSrv.forSingleQuery).toHaveBeenCalled();
    });
});
