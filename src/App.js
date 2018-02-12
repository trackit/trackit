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
    const retrieved = this.props.retrieved;
    const hasAccounts = (retrieved ? this.props.accounts.length > 0 : true);

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
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string,
    })
  ),
  retrieved: PropTypes.bool,
  getAccounts: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  retrieved: aws.accounts.retrieved,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(App);
