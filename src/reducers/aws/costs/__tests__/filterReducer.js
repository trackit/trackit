import FilterReducer from '../filterReducer';
import Constants from '../../../../constants';

describe("FilterReducer", () => {

  it("handles initial state", () => {
    expect(FilterReducer(undefined, {})).toEqual(null);
  });

  it("handles set filter state", () => {
    const filter = "filter";
    expect(FilterReducer(null, { type: Constants.AWS_SET_COSTS_FILTER, filter })).toEqual(filter);
  });

  it("handles clear filter state", () => {
    const filter = "filter";
    expect(FilterReducer(filter, { type: Constants.AWS_CLEAR_COSTS_FILTER })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const filter = "filter";
    expect(FilterReducer(filter, { type: "" })).toEqual(filter);
  });

});
