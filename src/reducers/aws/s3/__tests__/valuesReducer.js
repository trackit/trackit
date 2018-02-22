import ValuesReducer from '../valuesReducer';
import Constants from '../../../../constants';

describe("ValuesReducer", () => {

  const id = "id";
  const values = {
    value1: 1,
    value2: 2
  };

  const defaultValue = {};
  const requestedValue = {status: false};
  const successValue = {status: true, values};
  const errorValue = {status: true, error: Error()};

  it("handles initial state", () => {
    expect(ValuesReducer(undefined, {})).toEqual(defaultValue);
  });

  it("handles get values requested state", () => {
    expect(ValuesReducer(defaultValue, { type: Constants.AWS_GET_S3_DATA, id })).toEqual(requestedValue);
  });

  it("handles set values state", () => {
    expect(ValuesReducer(defaultValue, { type: Constants.AWS_GET_S3_DATA_SUCCESS, id, data: values })).toEqual(successValue);
  });

  it("handles error with values state", () => {
    expect(ValuesReducer(successValue, { type: Constants.AWS_GET_S3_DATA_ERROR, id, error: Error() })).toEqual(errorValue);
  });

  it("handles wrong type state", () => {
    expect(ValuesReducer(successValue, { type: "" })).toEqual(successValue);
  });

});
