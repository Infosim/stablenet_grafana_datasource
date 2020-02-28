/*

///<reference path="D:\git\stablenet_grafana_integration\stablenet-grafana-typescript\node_modules\@types\grafana\app\plugins\sdk.d.ts" />
// @ts-ignore
import { QueryCtrl } from 'grafana/app/plugins/sdk';
import {Target} from "../src/query_interfaces";
import {Mode, Unit} from "../src/types";
import {StableNetQueryCtrl} from "../src/query_ctrl";

let queryCtrl: StableNetQueryCtrl;

describe("constructor", () => {
    beforeEach(() => {
        queryCtrl = new StableNetQueryCtrl({},{});
    });

   it("should create an object", () => {
       expect(queryCtrl).not.toBeUndefined();
       expect(queryCtrl).not.toBeNull();
   });

    it('should initialise fields correctly', () => {
        let expectedTarget: Target = {
            averagePeriod: "",
            averageUnit: Unit.MINUTES,
            chosenMetrics: {},
            datasource: undefined,
            deviceQuery: "",
            includeAvgStats: true,
            includeMaxStats: false,
            includeMinStats: false,
            measurementQuery: "",
            metricPrefix: "",
            metrics: [],
            moreDevices: false,
            moreMeasurements: false,
            refId: "A",
            selectedDevice: -1,
            selectedMeasurement: -1,
            statisticLink: "",
            useCustomAverage: false,
            mode: Mode.MEASUREMENT
        };
        expect(queryCtrl.target).toEqual(expectedTarget);
    });
});

*/
