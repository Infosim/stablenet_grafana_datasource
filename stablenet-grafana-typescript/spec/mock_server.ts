import {TestOptions} from "../src/types";

export class MockBackendServer {
    constructor() {
    }

    async request(options: TestOptions) {
        return Promise.resolve(true);
    }

    datasourceRequest(options: TestOptions) {

    }
}
