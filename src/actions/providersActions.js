import * as actionTypes from '../constants/actionTypes';


export const setStorageTypes = (types) => {
  return ({
    type: actionTypes.SET_STORAGE_TYPES,
    types: types
  });
};

export const getPricingGcp = () => ({
  type: actionTypes.GET_PRICING_GCP
});

export const getPricingAws = () => ({
  type: actionTypes.GET_PRICING_AWS
});
