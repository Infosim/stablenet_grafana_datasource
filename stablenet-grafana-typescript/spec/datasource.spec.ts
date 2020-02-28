import {StableNetDatasource} from "../src/datasource";
import {MockBackendServer} from "./mock_server";
import {TestOptions} from "../src/types";
import {TestResult} from "../src/returnTypes";

let backendSrv: MockBackendServer;
let datasource: StableNetDatasource;

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

        spyOn(backendSrv,'request').and.returnValue(Promise.resolve(true));
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
        expect(await datasource.testDatasource()).toEqual(expected)
    });
});

describe("queryDevices()", () => {
    beforeEach(() => {
        backendSrv = new MockBackendServer();
        datasource = new StableNetDatasource({id:1},null,backendSrv);

        spyOn(backendSrv, 'datasourceRequest');
    });
});
