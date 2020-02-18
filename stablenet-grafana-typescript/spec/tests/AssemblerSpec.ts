import {WrappedTarget} from "../../src/data_query_assembler";

describe("A suite", function() {
    let s: WrappedTarget = new WrappedTarget({
        datasource: undefined,
        refId: "A",
        mode: 10,
        deviceQuery: 'string',
        selectedDevice: 0,
        measurementQuery: 'string',
        selectedMeasurement: 0,
        chosenMetrics: [],
        metricPrefix: 'string',
        includeMinStats: false,
        includeAvgStats: false,
        includeMaxStats: false,
        statisticLink: '',
        metrics: [],
        moreDevices: false,
        moreMeasurements: false
    },0);
    it("contains spec with an expectation", function() {
        expect(s).not.toBeNull();
    });
});
