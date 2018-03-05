import DatesReducer from '../datesReducer';
import Constants from '../../../constants';
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
    expect(DatesReducer({}, { type: Constants.DASHBOARD_INSERT_ITEMS_DATES, dates: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(DatesReducer({}, { type: Constants.DASHBOARD_ADD_ITEM, id })).toEqual(stateDefault);
  });

  it("handles set dates state", () => {
    expect(DatesReducer({}, { type: Constants.DASHBOARD_SET_ITEM_DATES, id, dates })).toEqual(state);
  });

  it("handles reset dates state", () => {
    expect(DatesReducer(state, { type: Constants.DASHBOARD_RESET_ITEMS_DATES, id, dates })).toEqual(stateDefault);
  });

  it("handles clear dates state", () => {
    expect(DatesReducer(state, { type: Constants.DASHBOARD_CLEAR_ITEMS_DATES })).toEqual({});
  });

  it("handles chart deletion state", () => {
    expect(DatesReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id })).toEqual({});
    expect(DatesReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id: "fakeID" })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(DatesReducer(state, { type: "" })).toEqual(state);
  });

});
