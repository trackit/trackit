import ValuesReducer from '../valuesReducer';
import Constants from '../../../../constants';

describe("ValuesReducer", () => {

  const id = "id";
  const costs = "costs";
  let state = {};
  state[id] = costs;

  it("handles initial state", () => {
    expect(ValuesReducer(undefined, {})).toEqual({});
  });

  it("handles set values state", () => {
    expect(ValuesReducer({}, { type: Constants.AWS_GET_COSTS_SUCCESS, id, costs })).toEqual(state);
  });

  it("handles error with values state", () => {
    expect(ValuesReducer(state, { type: Constants.AWS_GET_COSTS_ERROR })).toEqual({});
  });

  it("handles wrong type state", () => {
    expect(ValuesReducer(state, { type: "" })).toEqual(state);
  });

});
