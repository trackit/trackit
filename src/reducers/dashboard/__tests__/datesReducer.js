import DatesReducer from '../datesReducer';
import Constants from '../../../constants';
import moment from "moment/moment";

describe("DatesReducer", () => {

  const dates = {
    startDate: moment().subtract(1, 'week').startOf('week'),
    endDate: moment().subtract(1, 'week').endOf('week')
  };

  const state = dates;

  const defaultState = {
    startDate: moment().subtract(1, 'month').startOf('month'),
    endDate: moment().subtract(1, 'month').endOf('month')
  };

  it("handles initial state", () => {
    expect(DatesReducer(undefined, {})).toEqual({});
  });

  it("handles insert dates state", () => {
    expect(DatesReducer({}, { type: Constants.DASHBOARD_INSERT_DATES, dates })).toEqual(state);
  });

  it("handles set dates state", () => {
    expect(DatesReducer({}, { type: Constants.DASHBOARD_SET_DATES, dates })).toEqual(state);
  });

  it("handles reset dates state", () => {
    expect(DatesReducer(state, { type: Constants.DASHBOARD_RESET_DATES })).toEqual(defaultState);
  });

  it("handles clear dates state", () => {
    expect(DatesReducer(state, { type: Constants.DASHBOARD_CLEAR_DATES })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(DatesReducer(state, { type: "" })).toEqual(state);
  });

});
