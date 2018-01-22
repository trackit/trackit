import IntervalReducer from '../intervalReducer';
import Constants from '../../../../constants';

describe("IntervalReducer", () => {

  const id = "id";
  const interval = "interval";
  let state = {};
  state[id] = interval;

  it("handles initial state", () => {
    expect(IntervalReducer(undefined, {})).toEqual({});
  });

  it("handles set interval state", () => {
    expect(IntervalReducer({}, { type: Constants.AWS_SET_COSTS_INTERVAL, id, interval })).toEqual(state);
  });

  it("handles clear interval state", () => {
    expect(IntervalReducer(state, { type: Constants.AWS_CLEAR_COSTS_INTERVAL })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(IntervalReducer(state, { type: "" })).toEqual(state);
  });

});
