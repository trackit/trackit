import Constants from '../../constants';

export default {
  getViewers: () => ({ type: Constants.USER_GET_VIEWERS }),
  createViewer: email => ({ type: Constants.USER_NEW_VIEWER, email }),
};
