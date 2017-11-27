import Constants from '../../../constants';

export default (state=null, action) => {
  switch (action.type) {
    case Constants.AWS_NEW_EXTERNAL_SUCCESS:
      return action.external;
    default:
      return state;
  }
};
