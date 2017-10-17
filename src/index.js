import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import registerServiceWorker from './registerServiceWorker';
import { MuiThemeProvider, createMuiTheme } from 'material-ui/styles';
import { BrowserRouter, Route } from 'react-router-dom'

import Theme from './common/Theme';

// CSS inclusion for whole app
import 'bootstrap/dist/css/bootstrap.min.css';
import 'font-awesome/css/font-awesome.min.css';
import 'normalize.css/normalize.css';
import 'animate.css/animate.min.css';
import 'react-table/react-table.css'
import './styles/style.css';

// Components
import App from './App';
import Containers from './containers';


// Setup
import configureStore from './store';

// Initialize store
const store = configureStore();

// Creating theme
const theme = createMuiTheme(Theme.theme);

ReactDOM.render((
    <MuiThemeProvider theme={theme}>
      <Provider store={store}>
        <BrowserRouter>
          <div>
            <Route path="/app" component={App}/>
            <Route path="/login" component={Containers.Login}/>
          </div>
        </BrowserRouter>
      </Provider>
    </MuiThemeProvider>
    ), document.getElementById('root')
);
registerServiceWorker();
