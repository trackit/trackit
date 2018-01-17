import { put } from "redux-saga/effects";
import Constants from "../../constants/index";
import { unsetToken } from "../../common/localStorage";

export default function* cleanUserTokenSaga() {
  unsetToken();
  yield put({type: Constants.CLEAN_USER_TOKEN_SUCCESS});
}
