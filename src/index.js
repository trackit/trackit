import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import registerServiceWorker from './registerServiceWorker';
import { MuiThemeProvider, createMuiTheme } from 'material-ui/styles';
import { BrowserRouter, Route } from 'react-router-dom'

import Theme from './common/Theme';
import { PrivateRoute, IndexRedirect } from './common/routing';
import Constants from './constants';

// CSS inclusion for whole app
import 'bootstrap/dist/css/bootstrap.min.css';
import 'font-awesome/css/font-awesome.min.css';
import 'normalize.css/normalize.css';
import 'animate.css/animate.min.css';
import 'react-table/react-table.css'
import 'react-bootstrap-daterangepicker/css/daterangepicker.css';
import './styles/style.css';
import './styles/style-sm.css';
import './styles/style-xs.css';
import './styles/style-md.css';
import './styles/helpers.css';

// Components
import App from './App';
import Containers from './containers';

// Setup
import configureStore from './store';

// Initialize store
const store = configureStore();

// Creating theme
const theme = createMuiTheme(Theme.theme);

// Retrieving Token from localStorage if any
store.dispatch({ type: Constants.GET_USER_TOKEN });

ReactDOM.render((
    <MuiThemeProvider theme={theme}>
      <Provider store={store}>
        <BrowserRouter>
          <div>
            <Route exact path="/" component={IndexRedirect}/>
            <Route path="/login" component={Containers.Auth.Login}/>
            <Route path="/register" component={Containers.Auth.Register}/>
            <PrivateRoute path="/app" component={App} store={store}/>
          </div>
        </BrowserRouter>
      </Provider>
    </MuiThemeProvider>
    ), document.getElementById('root')
);
registerServiceWorker();
