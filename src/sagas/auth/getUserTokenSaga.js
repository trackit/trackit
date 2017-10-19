import { all, call, put } from "redux-saga/effects";
import Constants from "../../constants/index";
import { getToken } from "../../common/localStorage";

export default function* getUserTokenSaga() {
  try {
    const token = yield call(getToken);
    yield all([
      put({ type: Constants.GET_USER_TOKEN_SUCCESS, token }),
    ]);
  } catch (error) {
    yield put({ type: Constants.GET_USER_TOKEN_ERROR, error });
  }
}
