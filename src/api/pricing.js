import Config from '../config.js';

export default {
  getPricing: (provider, frequent, infrequent, archive) => {
    return fetch(`${Config.apiUrl}/${provider}/storage/cost?region=all&frequent=${frequent}&infrequent=${infrequent}&archive=${archive}`)
      .then(response => {
        return response.json()
      });
  }
};
