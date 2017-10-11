import Config from '../config.js';

export const types = () => {
  return fetch(`${Config.apiUrl}/all/storage/types`)
    .then(response => {
      return response.json()
    })
    .then(json => {
      return json.types;
    });
}

export const regions = (provider) => {
  return fetch(`${Config.apiUrl}/${provider}/regions`)
    .then(response => {
      return response.json()
    })
    .then(json => {
      return json.regions;
    });
}

export const getPricing = (provider, frequent, infrequent, archive) => {
  return fetch(`${Config.apiUrl}/${provider}/storage/cost?region=all&frequent=${frequent}&infrequent=${infrequent}&archive=${archive}`)
    .then(response => {
      return response.json()
    });
}
