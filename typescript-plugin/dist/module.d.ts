import RocksetDatasource from './datasource';
import { StableNetQueryCtrl } from './query_ctrl';
export declare class StableNetConfigCtrl {
    static templateUrl: string;
    private passwordExists;
    private current;
    constructor();
    resetPassword(): void;
}
declare class StableNetQueryOptionsCtrl {
    static templateUrl: string;
}
declare class StableNetAnnotationsQueryCtrl {
    static templateUrl: string;
}
export { RocksetDatasource as Datasource, StableNetQueryCtrl as QueryCtrl, StableNetConfigCtrl as ConfigCtrl, StableNetQueryOptionsCtrl as QueryOptionsCtrl, StableNetAnnotationsQueryCtrl as AnnotationsQueryCtrl };
