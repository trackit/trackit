import DatesReducer from '../datesReducer';
import Constants from '../../../../constants';

describe("DatesReducer", () => {

  const id = "id";
  const dates = "dates";
  let state = {};
  state[id] = dates;

  it("handles initial state", () => {
    expect(DatesReducer(undefined, {})).toEqual({});
  });

  it("handles set dates state", () => {
    expect(DatesReducer({}, { type: Constants.AWS_SET_COSTS_DATES, id, dates })).toEqual(state);
  });

  it("handles clear dates state", () => {
    expect(DatesReducer(state, { type: Constants.AWS_CLEAR_COSTS_DATES })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(DatesReducer(state, { type: "" })).toEqual(state);
  });

});
