export const setToken = (token) => {
  window.localStorage.setItem('userToken', token);
};

export const unsetToken = () => {
  window.localStorage.removeItem('userToken');
};

export const getToken = () => {
  return window.localStorage.getItem('userToken');
};
