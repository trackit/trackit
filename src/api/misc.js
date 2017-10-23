import Config from '../config.js';

export const call = (route, method, body=null) => {
  return fetch(Config.apiUrl + route, {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    body: (body !== null ? JSON.stringify(body) : null)
  }).then(data => (data.json()))
    .then(response => ({success: true, ...response }))
    .then(error => ({success: false, ...error}));
};
