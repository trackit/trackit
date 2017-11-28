import { put } from "redux-saga/effects";
import Constants from "../../constants/index";

export default function* cleanUserTokenSaga() {
  yield put({type: Constants.CLEAN_USER_TOKEN_SUCCESS});
}
