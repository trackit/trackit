import moment from 'moment';
import S3Reducer from '../s3Reducer';
import Constants from '../../../constants';

const s3Data = ["s3Data1", "s3Data2"];
const s3View = {
  startDate: moment(),
  endDate: moment()
};

const emptyState = {
  data: [],
  view: {}
};

const state = {
  data: s3Data,
  view: s3View
};

describe("S3Reducer", () => {

  it("handles initial state", () => {
    expect(S3Reducer(undefined, {})).toEqual({});
  });

  it("handles get data success state", () => {
    expect(S3Reducer(emptyState, { type: Constants.AWS_GET_S3_DATA_SUCCESS, s3Data })).toEqual({ ...emptyState, data: s3Data});
  });

  it("handles get data fail state", () => {
    expect(S3Reducer(state, { type: Constants.AWS_GET_S3_DATA_ERROR })).toEqual({...state, data: []});
  });

  it("handles set view dates success state", () => {
    expect(S3Reducer(emptyState, { type: Constants.AWS_SET_S3_VIEW_DATES, ...s3View })).toEqual({ ...emptyState, view: s3View});
  });

  it("handles wrong type state", () => {
    expect(S3Reducer(state, { type: "" })).toEqual(state);
  });

});
