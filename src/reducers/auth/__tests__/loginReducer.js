import LoginReducer from '../loginReducer';
import Constants from '../../../constants';

describe("LoginReducer", () => {

  it("handles initial state", () => {
    expect(LoginReducer(undefined, {})).toEqual(null);
  });

  it("handles get login success state", () => {
    const payload = {};
    expect(LoginReducer(null, { type: Constants.LOGIN_REQUEST_SUCCESS, ...payload })).toEqual({ status: true });
  });

  it("handles get login fail state", () => {
    const payload = { error: "error" };
    expect(LoginReducer(null, { type: Constants.LOGIN_REQUEST_ERROR, ...payload })).toEqual({ ...payload, status: false });
  });

  it("handles get login loading state", () => {
    expect(LoginReducer(null, { type: Constants.LOGIN_REQUEST_LOADING })).toEqual({});
  });

  it("handles clean login state", () => {
    expect(LoginReducer({status: true}, { type: Constants.LOGIN_REQUEST })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const payload = { status: true };
    expect(LoginReducer(payload, { type: "" })).toEqual(payload);
  });

});
