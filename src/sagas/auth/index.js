import { takeLatest } from 'redux-saga/effects';
import LoginSaga from './loginSaga';
import GetUserTokenSaga from './getUserTokenSaga';
import GetUserMailSaga from './getUserMailSaga';
import LogoutSaga from './logoutSaga';
import CleanUserTokenSaga from './cleanUserTokenSaga';
import CleanUserMailSaga from './cleanUserMailSaga';
import RegistrationSaga from './registrationSaga';
import Constants from '../../constants';

export function* watchGetLogin() {
  yield takeLatest(Constants.LOGIN_REQUEST, LoginSaga);
}

export function* watchGetToken() {
  yield takeLatest(Constants.GET_USER_TOKEN, GetUserTokenSaga);
}

export function* watchGetUserMail() {
  yield takeLatest(Constants.GET_USER_MAIL, GetUserMailSaga);
}

export function* watchGetLogout() {
  yield takeLatest(Constants.LOGOUT_REQUEST, LogoutSaga);
  yield takeLatest(Constants.GET_USER_TOKEN_ERROR, LogoutSaga);
  yield takeLatest(Constants.GET_USER_MAIL_ERROR, LogoutSaga);
}

export function* watchCleanToken() {
  yield takeLatest(Constants.CLEAN_USER_TOKEN, CleanUserTokenSaga);
}

export function* watchCleanUserMail() {
  yield takeLatest(Constants.CLEAN_USER_MAIL, CleanUserMailSaga);
}

export function* watchRegistration() {
  yield takeLatest(Constants.REGISTRATION_REQUEST, RegistrationSaga);
}
