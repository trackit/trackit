import ExternalReducer from '../externalReducer';
import Constants from '../../../../constants';

describe("ExternalReducer", () => {

  it("handles initial state", () => {
    expect(ExternalReducer(undefined, {})).toEqual(null);
  });

  it("handles get accounts success state", () => {
    const external = "test";
    expect(ExternalReducer(null, { type: Constants.AWS_NEW_EXTERNAL_SUCCESS, external })).toEqual(external);
  });

  it("handles get accounts fail state", () => {
    const external = "test";
    expect(ExternalReducer(external, { type: Constants.AWS_NEW_EXTERNAL_ERROR })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const external = "test";
    expect(ExternalReducer(external, { type: "" })).toEqual(external);
  });

});
