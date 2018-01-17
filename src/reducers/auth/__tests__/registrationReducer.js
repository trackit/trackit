import RegistrationReducer from '../registrationReducer';
import Constants from '../../../constants';

describe("RegistrationReducer", () => {

  it("handles initial state", () => {
    expect(RegistrationReducer(undefined, {})).toEqual(null);
  });

  it("handles get registration success state", () => {
    const payload = { status: true };
    expect(RegistrationReducer(null, { type: Constants.REGISTRATION_SUCCESS, payload })).toEqual(payload);
  });

  it("handles get registration fail state", () => {
    const payload = { status: false };
    expect(RegistrationReducer(null, { type: Constants.REGISTRATION_ERROR, payload })).toEqual(payload);
  });

  it("handles get registration loading state", () => {
    expect(RegistrationReducer(null, { type: Constants.REGISTRATION_REQUEST_LOADING })).toEqual({});
  });

  it("handles clean registration state", () => {
    const payload = { status: true };
    expect(RegistrationReducer(payload, { type: Constants.REGISTRATION_CLEAR })).toEqual(null);
  });

  it("handles wrong type state", () => {
    const payload = { status: true };
    expect(RegistrationReducer(payload, { type: "" })).toEqual(payload);
  });

});
