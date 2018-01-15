import { put } from "redux-saga/effects";
import Constants from "../../constants/index";
import { unsetUserMail } from "../../common/localStorage";

export default function* cleanUserMailSaga() {
  unsetUserMail();
  yield put({type: Constants.CLEAN_USER_MAIL_SUCCESS});
}
