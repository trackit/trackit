export const setToken = (token) => {
  window.localStorage.setItem('userToken', token);
}

export const getToken = () => {
  return window.localStorage.getItem('userToken');
}
