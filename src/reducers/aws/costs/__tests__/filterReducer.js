import FilterReducer from '../filterReducer';
import Constants from '../../../../constants';
import IntervalReducer from "../intervalReducer";

describe("FilterReducer", () => {

  const id = "id";
  const filter = "filter";
  let state = {};
  state[id] = filter;
  let stateDefault = {};
  stateDefault[id] = "product";
  let insert = {
    "id": "filter"
  };

  it("handles initial state", () => {
    expect(FilterReducer(undefined, {})).toEqual({});
  });

  it("handles insert filter state", () => {
    expect(FilterReducer({}, { type: Constants.AWS_INSERT_COSTS_FILTER, filter: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(FilterReducer({}, { type: Constants.AWS_ADD_CHART, id })).toEqual(stateDefault);
  });

  it("handles set filter state", () => {
    expect(FilterReducer({}, { type: Constants.AWS_SET_COSTS_FILTER, id, filter })).toEqual(state);
  });

  it("handles reset filter state", () => {
    expect(FilterReducer(state, { type: Constants.AWS_RESET_COSTS_FILTER, id, filter })).toEqual(stateDefault);
  });

  it("handles clear filter state", () => {
    expect(FilterReducer(state, { type: Constants.AWS_CLEAR_COSTS_FILTER })).toEqual({});
  });

  it("handles chart deletion state", () => {
    expect(FilterReducer(state, { type: Constants.AWS_REMOVE_CHART, id })).toEqual({});
    expect(FilterReducer(state, { type: Constants.AWS_REMOVE_CHART, id: "fakeID" })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(FilterReducer(state, { type: "" })).toEqual(state);
  });

});
