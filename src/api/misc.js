import Config from '../config.js';

export const call = (route, method, body=null, token=null) => {
  let headers = {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  };
  if (token !== null)
    headers['Authorization'] = token;
  return fetch(Config.apiUrl + route, {
    method,
    headers,
    body: (body !== null ? JSON.stringify(body) : null)
  }).then(data => (data.json()))
    .then(response => ({success: true, data: response }))
    .then(error => ({success: false, ...error}));
};
