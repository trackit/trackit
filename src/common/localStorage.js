export const setValue = (key, value) => {
  window.localStorage.setItem(key, value);
};


export const unsetValue = (key) => {
  window.localStorage.removeItem(key);
};

export const getValue = (key) => {
  return window.localStorage.getItem(key);
};

/* Token */

export const setToken = (token) => {
  setValue('userToken', token);
};

export const unsetToken = () => {
  unsetValue('userToken');
};

export const getToken = () => {
  return getValue('userToken');
};

/* User mail */

export const setUserMail = (mail) => {
  setValue('userMail', mail);
};

export const unsetUserMail = () => {
  unsetValue('userMail');
};

export const getUserMail = () => {
  return getValue('userMail');
};

/* Selected accounts */

export const setSelectedAccounts = (accounts) => {
  setValue('selectedAccounts', JSON.stringify(accounts));
};

export const getSelectedAccounts = () => {
  return JSON.parse(getValue('selectedAccounts'));
};

/* Cost Breakdown charts */

export const setCostBreakdownCharts = (charts) => {
  setValue('cb_charts', JSON.stringify(charts));
};

export const getCostBreakdownCharts = () => {
  return JSON.parse(getValue('cb_charts'));
};

/* S3 Analytics dates */

export const setS3Dates = (dates) => {
  setValue('s3_dates', JSON.stringify(dates));
};

export const getS3Dates = () => {
  return JSON.parse(getValue('s3_dates'));
};
