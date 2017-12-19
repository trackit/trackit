import IntervalReducer from '../intervalReducer';
import Constants from '../../../../constants';

describe("IntervalReducer", () => {

  it("handles initial state", () => {
    expect(IntervalReducer(undefined, {})).toEqual(null);
  });

  it("handles set interval state", () => {
    const interval = "interval";
    expect(IntervalReducer(null, { type: Constants.AWS_SET_COSTS_INTERVAL, interval })).toEqual(interval);
  });

  it("handles clear interval state", () => {
    const interval = "interval";
    expect(IntervalReducer(interval, { type: Constants.AWS_CLEAR_COSTS_INTERVAL })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const interval = "interval";
    expect(IntervalReducer(interval, { type: "" })).toEqual(interval);
  });

});
