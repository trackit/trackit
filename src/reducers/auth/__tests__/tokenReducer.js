import TokenReducer from '../tokenReducer';
import Constants from '../../../constants';

describe("TokenReducer", () => {

  it("handle initial state", () => {
    expect(TokenReducer(undefined, {})).toEqual(null);
  });

  it("handle get token success state", () => {
    const token = "token";
    expect(TokenReducer(null, { type: Constants.GET_USER_TOKEN_SUCCESS, token })).toEqual(token);
  });

  it("handle get token fail state", () => {
    expect(TokenReducer("token", { type: Constants.GET_USER_TOKEN_ERROR })).toEqual(null);
  });

  it("handle clean token success state", () => {
    expect(TokenReducer("token", { type: Constants.CLEAN_USER_TOKEN_SUCCESS })).toEqual(null);
  });

  it("handle clean token fail state", () => {
    expect(TokenReducer("token", { type: Constants.CLEAN_USER_TOKEN_ERROR })).toEqual(null);
  });

  it("handle wrong type state", () => {
    expect(TokenReducer("token", { type: "" })).toEqual("token");
  });

});