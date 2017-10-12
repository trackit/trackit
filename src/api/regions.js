import Config from '../config.js';

export default {
  regions: (provider) => {
    return fetch(`${Config.apiUrl}/${provider}/regions`)
      .then(response => {
        return response.json()
      })
      .then(json => {
        return json.regions;
      });
  }
};
