import { all, call, put } from "redux-saga/effects";
import Constants from "../../constants/index";
import { getToken } from "../../common/localStorage";

export default function* getUserTokenSaga() {
  try {
    const token = yield call(getToken);
    if (token)
      yield all([
        put({ type: Constants.GET_USER_TOKEN_SUCCESS, token }),
      ]);
    else
      throw Error("No token available");
  } catch (error) {
    yield put({ type: Constants.GET_USER_TOKEN_ERROR, error });
  }
}
