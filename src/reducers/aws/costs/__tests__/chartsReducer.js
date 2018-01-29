import ChartsReducer from '../chartsReducer';
import Constants from '../../../../constants';

describe("ChartsReducer", () => {

  const id = "id";
  const state = [id];
  const insert = state;

  it("handles initial state", () => {
    expect(ChartsReducer(undefined, {})).toEqual([]);
  });

  it("handles insert dates state", () => {
    expect(ChartsReducer([], { type: Constants.AWS_INSERT_CHARTS, charts: insert })).toEqual(insert);
  });

  it("handles add chart state", () => {
    expect(ChartsReducer([], { type: Constants.AWS_ADD_CHART, id })).toEqual(state);
  });

  it("handles remove chart state", () => {
    expect(ChartsReducer(state, { type: Constants.AWS_REMOVE_CHART, id })).toEqual([]);
    expect(ChartsReducer(state, { type: Constants.AWS_REMOVE_CHART, id: 42 })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(ChartsReducer(state, { type: "" })).toEqual(state);
  });

});
