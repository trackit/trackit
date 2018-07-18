import { put } from "redux-saga/effects";
import Constants from "../../constants/index";
import { unsetSelectedAccounts } from "../../common/localStorage";

export default function* cleanUserSelectedAccountsSaga() {
  unsetSelectedAccounts();
  yield put({type: Constants.CLEAN_USER_SELECTED_ACCOUNTS_SUCCESS});
}
