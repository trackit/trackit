import { call, put } from "redux-saga/effects";
import Constants from "../../constants/index";
import { getUserMail } from "../../common/localStorage";

export default function* getUserMailSaga() {
  try {
    const mail = yield call(getUserMail);
    if (mail)
      yield put({ type: Constants.GET_USER_MAIL_SUCCESS, mail });
    else
      throw Error("No user mail available");
  } catch (error) {
    yield put({ type: Constants.GET_USER_MAIL_ERROR, error });
  }
}
