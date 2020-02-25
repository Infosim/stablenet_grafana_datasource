import { WrappedTarget } from "../src/data_query_assembler";
import {Target} from "../src/query_interfaces";
import {Mode, SingleQuery, Unit} from "../src/types";

let t: Target = {
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
let wt: WrappedTarget;

describe("constructor", function () {
    it('should work properly', function () {
        wt = new WrappedTarget(t,0,0);
        expect(wt).not.toBeNull();
        expect(wt).not.toBeUndefined();
    });
});

describe("isValidStatisticLinkMode()", function() {
    it('should accept correct Targets', function () {
        t.statisticLink = 'https://127.0.0.1:5443/PlotServlet?id=3406&chart=5504&last=0,1440&offset=0,0&tz=Europe%2FBerlin&interval=60000&marker0=0.1%2C-16711936%2C0%2C0.1%2C1.0E10%2C%2C0%2C0%2C1%2C&marker1=-1.0E10%2C-16711936%2C0%2C0.1%2C0.0%2C%2C0%2C0%2C1%2C&marker2=0.0%2C-65536%2C0%2C0.1%2C0.1%2C%2C0%2C0%2C1%2C&quality=-1.0&dns=1&multiplecharts=0&multicharttype=0&width=1172&height=575&value0=1002&token=5B1388ECADB562B2011EE59F3E3861E5FC605DAED129BA869869A0826E3BAFC66A2F3D79A997AD2540D86BC23E212B2D1DF1868B3BB78B34072D49E8CB0009CA';
        t.mode = Mode.STATISTIC_LINK;
        wt = new WrappedTarget(t,0,0);
        expect(wt.isValidStatisticLinkMode()).toBe(true);
    });

    it('should check correct mode', function () {
        t.mode = Mode.MEASUREMENT;
        t.statisticLink = 'https://127.0.0.1:5443/PlotServlet?id=3406&chart=5504&last=0,1440&offset=0,0&tz=Europe%2FBerlin&interval=60000&marker0=0.1%2C-16711936%2C0%2C0.1%2C1.0E10%2C%2C0%2C0%2C1%2C&marker1=-1.0E10%2C-16711936%2C0%2C0.1%2C0.0%2C%2C0%2C0%2C1%2C&marker2=0.0%2C-65536%2C0%2C0.1%2C0.1%2C%2C0%2C0%2C1%2C&quality=-1.0&dns=1&multiplecharts=0&multicharttype=0&width=1172&height=575&value0=1002&token=5B1388ECADB562B2011EE59F3E3861E5FC605DAED129BA869869A0826E3BAFC66A2F3D79A997AD2540D86BC23E212B2D1DF1868B3BB78B34072D49E8CB0009CA';
        wt = new WrappedTarget(t,0,0);
        expect(wt.isValidStatisticLinkMode()).toBe(false);
    });

    it("should check non-empty statisticLink", function() {
        t.mode = Mode.STATISTIC_LINK;
        t.statisticLink = '';
        wt = new WrappedTarget(t,0,0);
        expect(wt.isValidStatisticLinkMode()).toBe(false);
    });
});

describe("hasEmptyMetrics()", function () {
    it('should check invalid chosenMetrics', function () {
        t.chosenMetrics = {};
        wt = new WrappedTarget(t,0,0);
        expect(wt.hasEmptyMetrics()).toBe(true);
    });

    it('should check whether any metrics were chosen', function () {
        t.chosenMetrics = {"SNMP": false, "TLS": false};
        wt = new WrappedTarget(t,0,0);
        expect(wt.hasEmptyMetrics()).toBe(true);
    });

    it('should accept Targets with at least one metric', function () {
        t.chosenMetrics = {"SNMP": false, "TLS": true};
        wt = new WrappedTarget(t,0,0);
        expect(wt.hasEmptyMetrics()).toBe(false);
    });
});

describe("toStatisticLinkQuery()",function () {
    it('should return correctly with default interval', function () {
        t.useCustomAverage = false;
        let other: SingleQuery = {
            refId: t.refId,
            datasourceId: 0,
            queryType: 'statisticLink',
            statisticLink: t.statisticLink,
            intervalMs: 0,
            includeMinStats: t.includeMinStats,
            includeAvgStats: t.includeAvgStats,
            includeMaxStats: t.includeMaxStats
        };
        wt = new WrappedTarget(t,0,0);
        expect(wt.toStatisticLinkQuery()).toEqual(other);
    });

    it('should return correctly with custom interval', function () {
        t.useCustomAverage = true;
        t.averagePeriod = '10';
        t.averageUnit = Unit.MINUTES;
        let other: SingleQuery = {
            refId: t.refId,
            datasourceId: 0,
            queryType: 'statisticLink',
            statisticLink: t.statisticLink,
            intervalMs: 600000,
            includeMinStats: t.includeMinStats,
            includeAvgStats: t.includeAvgStats,
            includeMaxStats: t.includeMaxStats
        };
        wt = new WrappedTarget(t,0,0);
        expect(wt.toStatisticLinkQuery()).toEqual(other);
    });

    it('should contain NaN for invalid custom interval', function () {
        t.useCustomAverage = true;
        t.averagePeriod = 'stablenet';
        t.averageUnit = Unit.MINUTES;
        wt = new WrappedTarget(t,0,0);
        expect(wt.toStatisticLinkQuery().intervalMs).toBeNaN();
    });
});

describe("toDeviceQuery()", function () {
    it('should return correctly with default interval', function () {
        t.chosenMetrics = {"SNMP_1": true};
        t.selectedMeasurement = 3701;
        t.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        t.metricPrefix = "XY";
        t.useCustomAverage = false;

        let other: SingleQuery = {
            refId: t.refId,
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
        wt = new WrappedTarget(t,0,0);
        expect(wt.toDeviceQuery()).toEqual(other);
    });

    it('should return correctly with custom interval', function () {
        t.chosenMetrics = {"SNMP_1": true};
        t.selectedMeasurement = 3701;
        t.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        t.metricPrefix = "XY";
        t.useCustomAverage = true;
        t.averagePeriod = '10';
        t.averageUnit = Unit.SECONDS;

        let other: SingleQuery = {
            refId: t.refId,
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
        wt = new WrappedTarget(t,0,0);
        expect(wt.toDeviceQuery()).toEqual(other);
    });

    it('should contain NaN for invalid custom interval', function () {
        t.chosenMetrics = {"SNMP_1": true};
        t.selectedMeasurement = 3701;
        t.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        t.metricPrefix = "XY";
        t.useCustomAverage = true;
        t.averagePeriod = 'stablenet';
        t.averageUnit = Unit.SECONDS;
        wt = new WrappedTarget(t,0,0);
        expect(wt.toDeviceQuery().intervalMs).toBeNaN();
    });

    it('should return correctly when no metrics are chosen', function () {
        t.chosenMetrics = {"SNMP_1": false};
        t.selectedMeasurement = 3701;
        t.metrics = [
            {measurementObid: 3701, key: "SNMP_1", text: "cpu-load 1min"},
            {measurementObid: 3701, key: "SNMP_2", text: "cpu-load 5min"}];
        t.metricPrefix = "XY";
        t.useCustomAverage = false;

        let other: SingleQuery = {
            refId: t.refId,
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
        wt = new WrappedTarget(t,0,0);
        expect(wt.toDeviceQuery()).toEqual(other);
    });
});
