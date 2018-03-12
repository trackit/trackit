import DatesReducer from '../datesReducer';
import Constants from '../../../../constants';
import moment from "moment/moment";

describe("DatesReducer", () => {

  const dates = {
    startDate: moment().subtract(1, 'month').startOf('month'),
    endDate: moment().subtract(1, 'month').endOf('month')
  };
  const insert = dates;
  const state = dates;

  it("handles initial state", () => {
    expect(DatesReducer(undefined, {})).toEqual({});
  });

  it("handles insert dates state", () => {
    expect(DatesReducer({}, { type: Constants.AWS_INSERT_S3_DATES, dates: insert })).toEqual(state);
  });

  it("handles set dates state", () => {
    expect(DatesReducer({}, { type: Constants.AWS_SET_S3_DATES, dates })).toEqual(state);
  });

  it("handles clear dates state", () => {
    expect(DatesReducer(state, { type: Constants.AWS_CLEAR_S3_DATES })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(DatesReducer(state, { type: "" })).toEqual(state);
  });

});
