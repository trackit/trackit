import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { Route } from 'react-router-dom';
import Containers from './containers';
import Actions from "./actions";

// Use named export for unconnected component (for tests)
export class App extends Component {

  componentWillMount() {
    this.props.getAccounts();
  }

  render() {

    const redirectToLogin = () => <Redirect to="/login/timeout"/>;
    const redirectToSetup = () => <Redirect to={`${this.props.match.url}/setup/false`}/>;
    const hasAccounts = (this.props.accounts.status ? (this.props.accounts.hasOwnProperty("values") && this.props.accounts.values && this.props.accounts.values.length > 0 ): true);

    const checkRedirections = (container) => (!this.props.token ? redirectToLogin : (hasAccounts ? container : redirectToSetup));

    return (
      <div>
        <Containers.Main>
          <div className="app-container">
            <Route
              path={this.props.match.url} exact
              component={checkRedirections(Containers.Home)}
            />
            <Route
              path={this.props.match.url + '/s3'}
              component={checkRedirections(Containers.AWS.S3Analytics)}
            />
            <Route
              path={this.props.match.url + '/costbreakdown'}
              component={checkRedirections(Containers.AWS.CostBreakdown)}
            />
            <Route
              path={this.props.match.url + '/reports'}
              component={checkRedirections(Containers.AWS.Reports)}
            />
            <Route
              path={this.props.match.url + "/map"}
              component={checkRedirections(Containers.AWS.ResourcesMap)}
            />
            <Route
              path={this.props.match.url + "/setup/:hasAccounts*"}
              component={this.props.token ? Containers.Setup.Main : redirectToLogin}
            />
          </div>
        </Containers.Main>
      </div>
    );
  }
}

App.propTypes = {
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        roleArn: PropTypes.string.isRequired,
        pretty: PropTypes.string,
        bills: PropTypes.arrayOf(
          PropTypes.shape({
            bucket: PropTypes.string.isRequired,
            path: PropTypes.string.isRequired
          })
        ),
      })
    ),
  }),
  getAccounts: PropTypes.func.isRequired,
  token: PropTypes.string
};

/* istanbul ignore next */
const mapStateToProps = ({aws, auth}) => ({
  accounts: aws.accounts.all,
  token: auth.token
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(App);
