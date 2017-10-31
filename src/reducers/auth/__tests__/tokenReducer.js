import TokenReducer from '../tokenReducer';
import Constants from '../../../constants';

describe("TokenReducer", () => {

  it("handles initial state", () => {
    expect(TokenReducer(undefined, {})).toEqual(null);
  });

  it("handles get token success state", () => {
    const token = "token";
    expect(TokenReducer(null, { type: Constants.GET_USER_TOKEN_SUCCESS, token })).toEqual(token);
  });

  it("handles get token fail state", () => {
    expect(TokenReducer("token", { type: Constants.GET_USER_TOKEN_ERROR })).toEqual(null);
  });

  it("handles clean token success state", () => {
    expect(TokenReducer("token", { type: Constants.CLEAN_USER_TOKEN_SUCCESS })).toEqual(null);
  });

  it("handles clean token fail state", () => {
    expect(TokenReducer("token", { type: Constants.CLEAN_USER_TOKEN_ERROR })).toEqual(null);
  });

  it("handles wrong type state", () => {
    expect(TokenReducer("token", { type: "" })).toEqual("token");
  });

});
