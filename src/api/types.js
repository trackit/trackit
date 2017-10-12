import Config from '../config.js';

export default {
  types: () => {
    return fetch(`${Config.apiUrl}/all/storage/types`)
      .then(response => {
        return response.json()
      })
      .then(json => {
        return json.types;
      });
  }
};
