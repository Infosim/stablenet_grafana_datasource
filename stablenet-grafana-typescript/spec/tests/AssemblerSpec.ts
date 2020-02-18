import {TextValue} from "../../src/returnTypes";

describe("A suite", function() {
    it("contains spec with an expectation", function() {
        let a: TextValue = {text: "a", value: 0};
        expect(a).not.toEqual({text:"b", value: 1});
    });
});
