import IntervalsReducer from '../intervalsReducer';
import Constants from '../../../constants';

describe("IntervalsReducer", () => {

  const id = "id";
  const interval = "interval";
  let state = {};
  state[id] = interval;
  let stateDefault = {};
  stateDefault[id] = "day";
  let stateDefaultPie = {};
  stateDefaultPie[id] = "month";
  let insert = {
    "id": "interval"
  };

  it("handles initial state", () => {
    expect(IntervalsReducer(undefined, {})).toEqual({});
  });

  it("handles insert interval state", () => {
    expect(IntervalsReducer({}, { type: Constants.DASHBOARD_INSERT_ITEMS_INTERVAL, intervals: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(IntervalsReducer({}, { type: Constants.DASHBOARD_ADD_ITEM, id, chartType: "bar" })).toEqual(stateDefault);
    expect(IntervalsReducer({}, { type: Constants.DASHBOARD_ADD_ITEM, id, chartType: "pie" })).toEqual(stateDefaultPie);
  });

  it("handles set interval state", () => {
    expect(IntervalsReducer({}, { type: Constants.DASHBOARD_SET_ITEM_INTERVAL, id, interval })).toEqual(state);
  });

  it("handles reset interval state", () => {
    expect(IntervalsReducer(state, { type: Constants.DASHBOARD_RESET_ITEMS_INTERVAL, id, interval })).toEqual(stateDefault);
  });

  it("handles clear interval state", () => {
    expect(IntervalsReducer(state, { type: Constants.DASHBOARD_CLEAR_ITEMS_INTERVAL })).toEqual({});
  });

  it("handles chart deletion state", () => {
    expect(IntervalsReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id })).toEqual({});
    expect(IntervalsReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id: "fakeID" })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(IntervalsReducer(state, { type: "" })).toEqual(state);
  });

});
