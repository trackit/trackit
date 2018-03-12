import ItemsReducer from '../itemsReducer';
import Constants from '../../../constants';

describe("ItemsReducer", () => {

  const id = "id";
  const props = "type";
  const state = {id: props};
  const insert = state;

  it("handles initial state", () => {
    expect(ItemsReducer(undefined, {})).toEqual({});
  });

  it("handles insert items state", () => {
    expect(ItemsReducer({}, { type: Constants.DASHBOARD_INSERT_ITEMS, items: insert })).toEqual(insert);
  });

  it("handles update items state", () => {
    expect(ItemsReducer({}, { type: Constants.DASHBOARD_UPDATE_ITEMS, items: insert })).toEqual(insert);
  });

  it("handles add item state", () => {
    expect(ItemsReducer({}, { type: Constants.DASHBOARD_ADD_ITEM, id, props})).toEqual(state);
  });

  it("handles remove chart state", () => {
    expect(ItemsReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id })).toEqual({});
    expect(ItemsReducer(state, { type: Constants.DASHBOARD_REMOVE_ITEM, id: 42 })).toEqual(state);
  });

  it("handles wrong type state", () => {
    expect(ItemsReducer(state, { type: "" })).toEqual(state);
  });

});
