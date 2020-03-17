import { WrappedTarget } from "../src/data_query_assembler";
import {Target} from "../src/query_interfaces";
import {Mode, SingleQuery, Unit} from "../src/types";

let target: Target;
let wrappedTarget: WrappedTarget;

describe("constructor", () => {
    beforeEach(() => {
        target = {
            datasource: undefined,
            refId: '',
            mode: -1,
            deviceQuery: '',
            selectedDevice: -1,
            measurementQuery: '',
            selectedMeasurement: -1,
            chosenMetrics: {},
            metricPrefix: '',
            includeMinStats: false,
            includeAvgStats: false,
            includeMaxStats: false,
            statisticLink: '',
            metrics: [],
            averagePeriod: '',
            averageUnit: -1,
            useCustomAverage: false,
            moreDevices: false,
            moreMeasurements: false
        };
    });

    it('should work properly', () => {
        wrappedTarget = new WrappedTarget(target,0,0);
        expect(wrappedTarget).not.toBeNull();
        expect(wrappedTarget).not.toBeUndefined();
    });
});

describe("isValidStatisticLinkMode()", () => {
    beforeEach(() => {
        target = {
            datasource: undefined,
            refId: '',
            mode: -1,
            deviceQuery: '',
            selectedDevice: -1,
            measurementQuery: '',
            selectedMeasurement: -1,
            chosenMetrics: {},
            metricPrefix: '',
            includeMinStats: false,
            includeAvgStats: false,
            includeMaxStats: false,
            statisticLink: '',
            metrics: [],
            averagePeriod: '',
            averageUnit: -1,
            useCustomAverage: false,
            moreDevices: false,
            moreMeasurements: false
        };
        wrappedTarget = new WrappedTarget(target,0,0);
    });

    it('should accept correct Targets', () => {
        target.statisticLink = 'https://127.0.0.1:5443/PlotServlet?id=3406&chart=5504&last=0,1440&offset=0,0&tz=Europe%2FBerlin&interval=60000&marker0=0.1%2C-16711936%2C0%2C0.1%2C1.0E10%2C%2C0%2C0%2C1%2C&marker1=-1.0E10%2C-16711936%2C0%2C0.1%2C0.0%2C%2C0%2C0%2C1%2C&marker2=0.0%2C-65536%2C0%2C0.1%2C0.1%2C%2C0%2C0%2C1%2C&quality=-1.0&dns=1&multiplecharts=0&multicharttype=0&width=1172&height=575&value0=1002&token=5B1388ECADB562B2011EE59F3E3861E5FC605DAED129BA869869A0826E3BAFC66A2F3D79A997AD2540D86BC23E212B2D1DF1868B3BB78B34072D49E8CB0009CA';
        target.mode = Mode.STATISTIC_LINK;
        expect(wrappedTarget.isValidStatisticLinkMode()).toBe(true);
    });

    it('should accept any link', () => {
        target.statisticLink = 'any link';
        target.mode = Mode.STATISTIC_LINK;
        expect(wrappedTarget.isValidStatisticLinkMode()).toBe(true);
    });

    it('should check correct mode', () => {
        target.mode = Mode.MEASUREMENT;
        target.statisticLink = 'https://127.0.0.1:5443/PlotServlet?id=3406&chart=5504&last=0,1440&offset=0,0&tz=Europe%2FBerlin&interval=60000&marker0=0.1%2C-16711936%2C0%2C0.1%2C1.0E10%2C%2C0%2C0%2C1%2C&marker1=-1.0E10%2C-16711936%2C0%2C0.1%2C0.0%2C%2C0%2C0%2C1%2C&marker2=0.0%2C-65536%2C0%2C0.1%2C0.1%2C%2C0%2C0%2C1%2C&quality=-1.0&dns=1&multiplecharts=0&multicharttype=0&width=1172&height=575&value0=1002&token=5B1388ECADB562B2011EE59F3E3861E5FC605DAED129BA869869A0826E3BAFC66A2F3D79A997AD2540D86BC23E212B2D1DF1868B3BB78B34072D49E8CB0009CA';
        expect(wrappedTarget.isValidStatisticLinkMode()).toBe(false);
    });

    it("should check non-empty statisticLink", () => {
        target.mode = Mode.STATISTIC_LINK;
        target.statisticLink = '';
        expect(wrappedTarget.isValidStatisticLinkMode()).toBe(false);
    });
});

describe("hasEmptyMetrics()", () => {
    beforeEach(() => {
        target = {
            datasource: undefined,
            refId: '',
            mode: -1,
            deviceQuery: '',
            selectedDevice: -1,
            measurementQuery: '',
            selectedMeasurement: -1,
            chosenMetrics: {},
            metricPrefix: '',
            includeMinStats: false,
            includeAvgStats: false,
            includeMaxStats: false,
            statisticLink: '',
            metrics: [],
            averagePeriod: '',
            averageUnit: -1,
            useCustomAverage: false,
            moreDevices: false,
            moreMeasurements: false
        };
        wrappedTarget = new WrappedTarget(target,0,0);
    });

    it('should check invalid chosenMetrics', () => {
        target.chosenMetrics = {};
        expect(wrappedTarget.hasEmptyMetrics()).toBe(true);
    });

    it('should check whether any metrics were chosen', () => {
        target.chosenMetrics = {"SNMP": false, "TLS": false};
        expect(wrappedTarget.hasEmptyMetrics()).toBe(true);
    });

    it('should accept Targets with at least one metric', () => {
        target.chosenMetrics = {"SNMP": false, "TLS": true};
        expect(wrappedTarget.hasEmptyMetrics()).toBe(false);
    });
});

describe("toStatisticLinkQuery()", () => {
    beforeEach(() => {
        target = {
            datasource: undefined,
            refId: '',
            mode: -1,
            deviceQuery: '',
            selectedDevice: -1,
            measurementQuery: '',
            selectedMeasurement: -1,
            chosenMetrics: {},
            metricPrefix: '',
            includeMinStats: false,
            includeAvgStats: false,
            includeMaxStats: false,
            statisticLink: '',
            metrics: [],
            averagePeriod: '',
            averageUnit: -1,
            useCustomAverage: false,
            moreDevices: false,
            moreMeasurements: false
        };
        wrappedTarget = new WrappedTarget(target,0,0);
    });

    it('should return correctly with default interval', () => {
        target.useCustomAverage = false;
        let expected: SingleQuery = {
            refId: target.refId,
            datasourceId: 0,
            queryType: 'statisticLink',
            statisticLink: target.statisticLink,
            intervalMs: 0,
            includeMinStats: target.includeMinStats,
            includeAvgStats: target.includeAvgStats,
            includeMaxStats: target.includeMaxStats
        };
        expect(wrappedTarget.toStatisticLinkQuery()).toEqual(expected);
    });

    it('should return correctly with custom interval', () => {
        target.useCustomAverage = true;
        target.averagePeriod = '10';
        target.averageUnit = Unit.MINUTES;
        let expected: SingleQuery = {
            refId: target.refId,
            datasourceId: 0,
            queryType: 'statisticLink',
            statisticLink: target.statisticLink,
            intervalMs: 600000,
            includeMinStats: target.includeMinStats,
            includeAvgStats: target.includeAvgStats,
            includeMaxStats: target.includeMaxStats
        };
        expect(wrappedTarget.toStatisticLinkQuery()).toEqual(expected);
    });

    it('should contain NaN for invalid custom interval', () => {
        target.useCustomAverage = true;
        target.averagePeriod = 'stablenet';
        target.averageUnit = Unit.MINUTES;
        expect(wrappedTarget.toStatisticLinkQuery().intervalMs).toBeNaN();
    });
});

describe("toDeviceQuery()", () => {
    beforeEach(() => {
        target = {
            datasource: undefined,
            refId: '',
            mode: -1,
            deviceQuery: '',
            selectedDevice: -1,
            measurementQuery: '',
            selectedMeasurement: -1,
            chosenMetrics: {},
            metricPrefix: '',
            includeMinStats: false,
            includeAvgStats: false,
            includeMaxStats: false,
            statisticLink: '',
            metrics: [],
            averagePeriod: '',
            averageUnit: -1,
            useCustomAverage: false,
            moreDevices: false,
            moreMeasurements: false
        };
        wrappedTarget = new WrappedTarget(target,0,0);
    });

    it('should return correctly with default interval', () => {
        target.chosenMetrics = {"SNMP_1": true};
        target.selectedMeasurement = 3701;
        target.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        target.metricPrefix = "XY";
        target.useCustomAverage = false;

        let expected: SingleQuery = {
            refId: target.refId,
            datasourceId: 0,
            queryType: 'metricData',
            includeAvgStats:false,
            includeMaxStats:false,
            includeMinStats:false,
            intervalMs:0,
            requestData: [{
                measurementObid: 3701,
                metrics:[{key:"SNMP_1", name:"XY {MinMaxAvg} cpu-load 1min"}]
            }]
        };
        wrappedTarget = new WrappedTarget(target,0,0);
        expect(wrappedTarget.toDeviceQuery()).toEqual(expected);
    });

    it('should return correctly with custom interval', () => {
        target.chosenMetrics = {"SNMP_1": true};
        target.selectedMeasurement = 3701;
        target.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        target.metricPrefix = "XY";
        target.useCustomAverage = true;
        target.averagePeriod = '10';
        target.averageUnit = Unit.SECONDS;

        let expected: SingleQuery = {
            refId: target.refId,
            datasourceId: 0,
            queryType: 'metricData',
            includeAvgStats:false,
            includeMaxStats:false,
            includeMinStats:false,
            intervalMs:10000,
            requestData: [{
                measurementObid: 3701,
                metrics:[{key:"SNMP_1", name:"XY {MinMaxAvg} cpu-load 1min"}]
            }]
        };
        wrappedTarget = new WrappedTarget(target,0,0);
        expect(wrappedTarget.toDeviceQuery()).toEqual(expected);
    });

    it('should contain NaN for invalid custom interval', () => {
        target.chosenMetrics = {"SNMP_1": true};
        target.selectedMeasurement = 3701;
        target.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        target.metricPrefix = "XY";
        target.useCustomAverage = true;
        target.averagePeriod = 'stablenet';
        target.averageUnit = Unit.SECONDS;
        wrappedTarget = new WrappedTarget(target,0,0);
        expect(wrappedTarget.toDeviceQuery().intervalMs).toBeNaN();
    });

    it('should return correctly when no metrics are chosen', () => {
        target.chosenMetrics = {"SNMP_1": false};
        target.selectedMeasurement = 3701;
        target.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        target.metricPrefix = "XY";
        target.useCustomAverage = false;

        let expected: SingleQuery = {
            refId: target.refId,
            datasourceId: 0,
            queryType: 'metricData',
            includeAvgStats:false,
            includeMaxStats:false,
            includeMinStats:false,
            intervalMs:0,
            requestData: [{
                measurementObid: 3701,
                metrics:[]
            }]
        };
        wrappedTarget = new WrappedTarget(target,0,0);
        expect(wrappedTarget.toDeviceQuery()).toEqual(expected);
    });
});
