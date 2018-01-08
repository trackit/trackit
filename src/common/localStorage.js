export const setValue = (key, value) => {
  window.localStorage.setItem(key, value);
};

export const unsetValue = (key) => {
  window.localStorage.removeItem(key);
};

export const getValue = (key) => {
  return window.localStorage.getItem(key);
};

export const setToken = (token) => {
  setValue('userToken', token);
};

export const unsetToken = () => {
  unsetValue('userToken');
};

export const getToken = () => {
  return getValue('userToken');
};

export const setUserMail = (mail) => {
  setValue('userMail', mail);
};

export const unsetUserMail = () => {
  unsetValue('userMail');
};

export const getUserMail = () => {
  return getValue('userMail');
};
