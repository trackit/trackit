import ExternalReducer from '../externalReducer';
import Constants from '../../../../constants';

describe("ExternalReducer", () => {

  it("handle initial state", () => {
    expect(ExternalReducer(undefined, {})).toEqual(null);
  });

  it("handle get accounts success state", () => {
    const external = "test";
    expect(ExternalReducer(null, { type: Constants.AWS_NEW_EXTERNAL_SUCCESS, external })).toEqual(external);
  });

  it("handle get accounts fail state", () => {
    const external = "test";
    expect(ExternalReducer(external, { type: Constants.AWS_NEW_EXTERNAL_ERROR })).toEqual(null);
  });

  it("handle wrong type state", () => {
    const external = "test";
    expect(ExternalReducer(external, { type: "" })).toEqual(external);
  });

});