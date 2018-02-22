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

    const redirectToSetup = () => <Redirect to={`${this.props.match.url}/setup/false`}/>;
    const hasAccounts = (this.props.accounts.status && this.props.accounts.hasOwnProperty("values") ? this.props.accounts.values.length > 0 : true);

    return (
      <div>
        <Containers.Main>
          <div className="app-container">
            <Route
              path={this.props.match.url} exact
              component={hasAccounts ? Containers.Home : redirectToSetup}
            />
            <Route
              path={this.props.match.url + '/s3'}
              component={hasAccounts ? Containers.AWS.S3Analytics : redirectToSetup}
            />
            <Route
              path={this.props.match.url + '/costbreakdown'}
              component={hasAccounts ? Containers.AWS.CostBreakdown : redirectToSetup}
            />
            <Route
              path={this.props.match.url + "/setup/:hasAccounts*"}
              component={Containers.Setup.Main}
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
  getAccounts: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(App);
