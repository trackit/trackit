import Constants from '../../constants';

export default {
  getViewers: () => ({ type: Constants.USER_GET_VIEWERS }),
  createViewer: email => ({ type: Constants.USER_CREATE_VIEWER, email }),
};
