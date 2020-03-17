import {DeviceQuery, MeasurementQuery, MetricQuery, SingleQuery, TestOptions} from "../src/types";
import {EntityQueryResult, GenericResponse, MetricType, TSDBArg} from "../src/returnTypes";

export class MockBackendServer {
    constructor() {}

    async request(options: TestOptions): Promise<boolean> {
        return Promise.resolve(true);
    }

    datasourceRequest(options: TestOptions): any {
        let data = options.data;

        if (MockBackendServer.instanceofDeviceQuery(data)){
            return this.forDeviceQuery(data);
        }

        if (MockBackendServer.instanceofMeasurementQuery(data)){
            return this.forMeasurementQuery(data);
        }

        if (MockBackendServer.instanceofMetricQuery(data)){
            return this.forMetricQuery(data);
        }

        if (MockBackendServer.instanceofSingleQuery(data)){
            return this.forSingleQuery(data);
        }

        throw new Error("Calling backend with invalid data");
    }

    forDeviceQuery(data: DeviceQuery): Promise<GenericResponse<EntityQueryResult>> {
        let res: GenericResponse<EntityQueryResult> = {
            config: {},
            data: {
                results: {
                    "A": {
                        refId: "A",
                        meta: {
                            data: [
                                {name: "ada", obid: 1008},
                                {name: "adfsserver2.root.infosim.net", obid: 1039},
                                {name: "agentl.infosim.net", obid: 1068}
                            ],
                            hasMore: false
                        }
                    }
                }
            },
            headers: () => {},
            status: 0,
            statusText: "",
            xhrStatus: ""
        };

        return Promise.resolve(res)
    }

    forMeasurementQuery(data: MeasurementQuery): Promise<GenericResponse<EntityQueryResult>> {
        let res: GenericResponse<EntityQueryResult> = {
            config: {},
            data: {
                results: {
                    "A": {
                        refId: "A",
                        meta: {
                            hasMore: false,
                            data: [
                                {name: "Berlin", obid: 3691},
                                {name: "Berlin Cisco CPU 1", obid: 3701},
                                {name: "Berlin Cisco Dial-In Statistics", obid: 3707}
                            ]
                        }
                    }
                }
            },
            headers: () => {},
            status: 0,
            statusText: "",
            xhrStatus: ""
        };

        return Promise.resolve(res);
    }

    forMetricQuery(data: MetricQuery): Promise<GenericResponse<MetricType[]>> {
        let res: GenericResponse<MetricType[]> = {
            config: {},
            data: {
                results: {
                    "A": {
                        refId: "A",
                        tables: null,
                        series: [],
                        meta: [
                            {key: "SNMP_1000", name: "System Users"},
                            {key: "SNMP_1001", name: "System Processes"},
                            {key: "SNMP_1002", name: "System Uptime"}
                        ]
                    }
                }
            },
            headers: () => {},
            status: 0,
            statusText: "",
            xhrStatus: ""
        };

        return Promise.resolve(res);
    }

    forSingleQuery(data: SingleQuery): Promise<TSDBArg> {
        let res: TSDBArg = {
            config: {},
            data: {
                results: {
                    "A": {
                        refId: "A",
                        series: [
                            {
                                name: "Berlin Cisco CPU 1 Avg cpu-load 1min",
                                points: [
                                    [7, 1583839989933],
                                    [33, 1583840289933],
                                    [8, 1583840589933],
                                    [7, 1583840889933],
                                    [20, 1583841189933],
                                    [7, 1583841489933]
                                ]
                            }
                        ]
                    }
                }
            },
            headers: () => {},
            status: 0,
            statusText: "",
            xhrStatus: ""
        };

        return Promise.resolve(res);
    }

    private static instanceofSingleQuery(data: any): data is SingleQuery {
        return 'from' in data;
    }

    private static instanceofMetricQuery(data: any): data is MetricQuery {
        return 'queries' in data && data.queries[0].queryType === 'metricNames';
    }

    private static instanceofMeasurementQuery(data: any): data is MeasurementQuery {
        return 'queries' in data && data.queries[0].queryType === 'measurements';
    }

    private static instanceofDeviceQuery(data: any): data is DeviceQuery {
        return 'queries' in data && data.queries[0].queryType === 'devices';
    }
}
