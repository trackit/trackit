import S3Reducer from '../s3Reducer';
import Constants from '../../../constants';

describe("S3Reducer", () => {

  it("handles initial state", () => {
    expect(S3Reducer(undefined, {})).toEqual([]);
  });

  it("handles get accounts success state", () => {
    const s3Data = ["s3Data1", "s3Data2"];
    expect(S3Reducer(null, { type: Constants.AWS_GET_S3_DATA_SUCCESS, s3Data })).toEqual(s3Data);
  });

  it("handles get accounts fail state", () => {
    const s3Data = ["s3Data1", "s3Data2"];
    expect(S3Reducer(s3Data, { type: Constants.AWS_GET_S3_DATA_ERROR })).toEqual([]);
  });

  it("handles wrong type state", () => {
    const s3Data = ["s3Data1", "s3Data2"];
    expect(S3Reducer(s3Data, { type: "" })).toEqual(s3Data);
  });

});
