import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import registerServiceWorker from './registerServiceWorker';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
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

// Retrieving user Token, user Mail and CostBreakdown charts from localStorage if any
store.dispatch({ type: Constants.GET_USER_TOKEN });
store.dispatch({ type: Constants.GET_USER_MAIL });
store.dispatch({ type: Constants.AWS_LOAD_SELECTED_ACCOUNTS });
store.dispatch({ type: Constants.AWS_LOAD_CHARTS });
store.dispatch({ type: Constants.AWS_LOAD_S3_DATES });
store.dispatch({ type: Constants.DASHBOARD_LOAD_ITEMS });

ReactDOM.render((
    <MuiThemeProvider theme={theme}>
      <Provider store={store}>
        <BrowserRouter>
          <div>
            <Route path="/" component={IndexRedirect} exact/>
            <Route path="/login" component={Containers.Auth.Login} exact/>
            <Route path="/login/timeout" component={Containers.Auth.Login} exact/>
            <Route path="/login/:prefill" component={Containers.Auth.Login}/>
            <Route path="/register" component={Containers.Auth.Register}/>
            <Route path="/forgot" component={Containers.Auth.Forgot}/>
            <Route path="/reset/:id/:token" component={Containers.Auth.Renew}/>
            <PrivateRoute path="/app" component={App} store={store}/>
          </div>
        </BrowserRouter>
      </Provider>
    </MuiThemeProvider>
    ), document.getElementById('root')
);
registerServiceWorker();
