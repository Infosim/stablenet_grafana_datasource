import { QueryCtrl } from 'app/plugins/sdk';
export declare class StableNetQueryCtrl extends QueryCtrl {
    static templateUrl: string;
    constructor($scope: any, $injector: any);
    getModes(): {
        text: string;
        value: number;
    }[];
    onModeChange(): void;
    onDeviceQueryChange(): void;
    getDevices(): any;
    onDeviceChange(): void;
    getMeasurements(): any;
    onMeasurementRegexChange(): void;
    onMeasurementChange(): void;
    onChangeInternal(): void;
}
