import ValuesReducer from '../valuesReducer';
import Constants from '../../../../constants';

describe("ValuesReducer", () => {

  it("handles initial state", () => {
    expect(ValuesReducer(undefined, {})).toEqual(null);
  });

  it("handles set values state", () => {
    const costs = { cost: 0, anotherCost: 1 };
    expect(ValuesReducer(null, { type: Constants.AWS_GET_COSTS_SUCCESS, costs })).toEqual(costs);
  });

  it("handles request values state", () => {
    const costs = { cost: 0, anotherCost: 1 };
    expect(ValuesReducer(costs, { type: Constants.AWS_GET_COSTS })).toEqual(null);
  });

  it("handles error with values state", () => {
    const costs = { cost: 0, anotherCost: 1 };
    expect(ValuesReducer(costs, { type: Constants.AWS_GET_COSTS_ERROR })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const costs = { cost: 0, anotherCost: 1 };
    expect(ValuesReducer(costs, { type: "" })).toEqual(costs);
  });

});
