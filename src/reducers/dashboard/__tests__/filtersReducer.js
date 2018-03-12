import FilterReducer from '../filtersReducer';
import Constants from '../../../constants';

describe("FilterReducer", () => {

  const id = "id";
  const filter = "filter";
  let state = {};
  state[id] = filter;
  let stateDefault = {};
  stateDefault[id] = "product";
  let stateNull = {};
  stateNull[id] = null;
  let insert = {
    "id": "filter"
  };

  const props = {
    type: "cb_bar"
  };

  it("handles initial state", () => {
    expect(FilterReducer(undefined, {})).toEqual({});
  });

  it("handles insert filter state", () => {
    expect(FilterReducer({}, { type: Constants.DASHBOARD_INSERT_ITEMS_FILTER, filters: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(FilterReducer({}, { type: Constants.DASHBOARD_ADD_ITEM, id, props })).toEqual(stateDefault);
  });

  it("handles set filter state", () => {
    expect(FilterReducer({}, { type: Constants.DASHBOARD_SET_ITEM_FILTER, id, filter })).toEqual(state);
  });

  it("handles reset filter state", () => {
    expect(FilterReducer(state, { type: Constants.DASHBOARD_RESET_ITEMS_FILTER, id, filter })).toEqual(stateNull);
  });

  it("handles clear filter state", () => {
    expect(FilterReducer(state, { type: Constants.DASHBOARD_CLEAR_ITEMS_FILTER })).toEqual({});
  });

  it("handles chart deletion state", () => {
    expect(FilterReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id })).toEqual({});
    expect(FilterReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id: "fakeID" })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(FilterReducer(state, { type: "" })).toEqual(state);
  });

});
