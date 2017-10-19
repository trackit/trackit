import { all, put } from "redux-saga/effects";
import Constants from "../../constants/index";

export default function* getUserTokenSaga() {
  try {
    yield all([
      put({type: Constants.CLEAN_USER_TOKEN_SUCCESS}),
    ]);
  } catch (error) {
    yield put({type: Constants.CLEAN_USER_TOKEN_ERROR, error});
  }
}
