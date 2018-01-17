import FilterReducer from '../filterReducer';
import Constants from '../../../../constants';

describe("FilterReducer", () => {

  const id = "id";
  const filter = "filter";
  let state = {};
  state[id] = filter;

  it("handles initial state", () => {
    expect(FilterReducer(undefined, {})).toEqual({});
  });

  it("handles set filter state", () => {
    expect(FilterReducer({}, { type: Constants.AWS_SET_COSTS_FILTER, id, filter })).toEqual(state);
  });

  it("handles clear filter state", () => {
    expect(FilterReducer(state, { type: Constants.AWS_CLEAR_COSTS_FILTER })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(FilterReducer(state, { type: "" })).toEqual(state);
  });

});
