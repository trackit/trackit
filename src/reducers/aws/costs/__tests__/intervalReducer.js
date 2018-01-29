import IntervalReducer from '../intervalReducer';
import Constants from '../../../../constants';
import ValuesReducer from "../valuesReducer";
import DatesReducer from "../datesReducer";

describe("IntervalReducer", () => {

  const id = "id";
  const interval = "interval";
  let state = {};
  state[id] = interval;
  let stateDefault = {};
  stateDefault[id] = "day";
  let insert = {
    "id": "interval"
  };

  it("handles initial state", () => {
    expect(IntervalReducer(undefined, {})).toEqual({});
  });

  it("handles insert interval state", () => {
    expect(IntervalReducer({}, { type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(IntervalReducer({}, { type: Constants.AWS_ADD_CHART, id })).toEqual(stateDefault);
  });

  it("handles set interval state", () => {
    expect(IntervalReducer({}, { type: Constants.AWS_SET_COSTS_INTERVAL, id, interval })).toEqual(state);
  });

  it("handles reset interval state", () => {
    expect(IntervalReducer(state, { type: Constants.AWS_RESET_COSTS_INTERVAL, id, interval })).toEqual(stateDefault);
  });

  it("handles clear interval state", () => {
    expect(IntervalReducer(state, { type: Constants.AWS_CLEAR_COSTS_INTERVAL })).toEqual({});
  });

  it("handles chart deletion state", () => {
    expect(IntervalReducer(state, { type: Constants.AWS_REMOVE_CHART, id })).toEqual({});
    expect(IntervalReducer(state, { type: Constants.AWS_REMOVE_CHART, id: "fakeID" })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(IntervalReducer(state, { type: "" })).toEqual(state);
  });

});
