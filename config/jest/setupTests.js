import 'jest-enzyme';
import { configure } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

configure({ adapter: new Adapter() });

global.requestAnimationFrame = (callback) => {
  setTimeout(callback, 0);
};

Object.defineProperty(window, 'localStorage', {
  value: (function() {
    let store = {};

    return {
      getItem: function(key) {
        return store[key] || null;
      },
      setItem: function(key, value) {
        store[key] = value.toString();
      },
      removeItem: function(key) {
        delete store[key];
      },
      clear: function() {
        store = {};
      }
    };

  })()
});

Object.defineProperty(window, 'Plotly', {
  value: ({
    newPlot: (id, data, layout, params) => ({id, data, layout, params})
  })
})
