import DatesReducer from '../datesReducer';
import Constants from '../../../../constants';
import moment from "moment/moment";

describe("DatesReducer", () => {

  const id = "id";
  const dates = "dates";
  let state = {};
  state[id] = dates;
  let stateDefault = {};
  stateDefault[id] = {
    startDate: moment().subtract(1, 'month').startOf('month'),
    endDate: moment().subtract(1, 'month').endOf('month')
  };
  let insert = {
    "id": {
      startDate: moment().subtract(1, 'month').startOf('month'),
      endDate: moment().subtract(1, 'month').endOf('month')
    }
  };

  it("handles initial state", () => {
    expect(DatesReducer(undefined, {})).toEqual({});
  });

  it("handles insert dates state", () => {
    expect(DatesReducer({}, { type: Constants.AWS_INSERT_COSTS_DATES, dates: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(DatesReducer({}, { type: Constants.AWS_ADD_CHART, id })).toEqual(stateDefault);
  });

  it("handles set dates state", () => {
    expect(DatesReducer({}, { type: Constants.AWS_SET_COSTS_DATES, id, dates })).toEqual(state);
  });

  it("handles reset dates state", () => {
    expect(DatesReducer(state, { type: Constants.AWS_RESET_COSTS_DATES, id, dates })).toEqual(stateDefault);
  });

  it("handles clear dates state", () => {
    expect(DatesReducer(state, { type: Constants.AWS_CLEAR_COSTS_DATES })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(DatesReducer(state, { type: "" })).toEqual(state);
  });

});
