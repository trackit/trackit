import Constants from '../../../constants';

const defaultValue = {status: false, values: []};

export default (state=defaultValue, action) => {
  switch (action.type) {
    case Constants.AWS_GET_ACCOUNTS:
      return {status: false, values: state.values};
    case Constants.AWS_GET_ACCOUNTS_SUCCESS:
      return {status: true, values: action.accounts};
    case Constants.AWS_GET_ACCOUNTS_ERROR:
      return {status: true, error: action.error, values: state.values};
    case Constants.AWS_GET_ACCOUNTS_CLEAR:
      return {status: true, values: []};
    default:
      return state;
  }
};
