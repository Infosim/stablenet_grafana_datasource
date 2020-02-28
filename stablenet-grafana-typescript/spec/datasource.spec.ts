import {StableNetDatasource} from "../src/datasource";

const MockBackendServer = function() {
    const request = function (options: any) {};

    const datasourceRequest = function (options: any) {};

    return {
        request: request,
        datasourceRequest: datasourceRequest
    }
};

let ds: StableNetDatasource = new StableNetDatasource({id: 1},{}, MockBackendServer());

describe("constructor", function () {
    it('should create a valid object', function () {
        expect(ds).not.toBeNull();
        expect(ds).not.toBeUndefined();
    });
});

