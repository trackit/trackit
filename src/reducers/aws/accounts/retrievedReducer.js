import Constants from '../../../constants';

export default (state=false, action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNTS_SUCCESS:
      return true;
    default:
      return state;
  }
};
