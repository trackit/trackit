import DatesReducer from '../datesReducer';
import Constants from '../../../../constants';
import moment from 'moment';

describe("DatesReducer", () => {

  it("handles initial state", () => {
    expect(DatesReducer(undefined, {})).toEqual(null);
  });

  it("handles set dates state", () => {
    const dates = {
      startDate: moment().startOf('month'),
      endDate: moment()
    };
    expect(DatesReducer(null, { type: Constants.AWS_SET_COSTS_DATES, dates })).toEqual(dates);
  });

  it("handles clear dates state", () => {
    const dates = {
      startDate: moment().startOf('month'),
      endDate: moment()
    };
    expect(DatesReducer(dates, { type: Constants.AWS_CLEAR_COSTS_DATES })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const dates = {
      startDate: moment().startOf('month'),
      endDate: moment()
    };
    expect(DatesReducer(dates, { type: "" })).toEqual(dates);
  });

});
