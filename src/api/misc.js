import Config from '../config';

const handleResponse = (response) => {
  if (response.status === 401)
    throw response.statusText;
  return response.json();
};

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
  }).then(handleResponse)
    .then(response => ({success: true, data: response }))
    .then(error => ({success: false, ...error}))
    .catch(error => ({success: null, error}));
};

export const download = (route, method='GET', body=null, token=null, contentType='text/csv') => {
  let headers = {
    'Content-Type': 'application/json',
    'Accept': contentType
  };
  if (token !== null)
    headers['Authorization'] = token;
  return fetch(Config.apiUrl + route, {
    method,
    headers,
    body: (body !== null ? JSON.stringify(body) : null)
  }).then(data => (data.blob()))
    .then(response => ({success: true, data: response }))
    .then(error => ({success: false, ...error}));
};
