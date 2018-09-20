import React, {Component} from 'react';
import {connect} from 'react-redux';
import {Redirect} from "react-router-dom";
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions/index';

const Form = Components.Auth.Form;

// LoginContainer Component
export class LoginContainer extends Component {

  componentWillUnmount() {
    this.props.clear();
  }

  componentWillMount() {
    this.props.clear();
  }

  render() {
    const awstoken = (this.props.match.params.hasOwnProperty("awstoken") ? this.props.match.params.awstoken : "");
    if (this.props.token)
      return (<Redirect to="/"/>);
    return (<Form
      awsToken={decodeURIComponent(awstoken)}
      submit={this.props.login}
      loginStatus={this.props.loginStatus}
      registrationStatus={this.props.registrationStatus}
      timeout={(this.props.match.params && this.props.match.params.prefill && this.props.match.params.prefill === "timeout")}
      />);
  }

}

LoginContainer.propTypes = {
  login: PropTypes.func.isRequired,
  token: PropTypes.string,
  loginStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  token: state.auth.token,
  loginStatus: state.auth.loginStatus,
  registrationStatus: state.auth.registration
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  login: (email, password, awsToken) => {
    dispatch(Actions.Auth.login(email, password, awsToken))
  },
  clear: () => {
    dispatch(Actions.Auth.clearRegister());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(LoginContainer);
