import React, {Component} from 'react';
import {connect} from 'react-redux';
import {Redirect} from "react-router-dom";
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions/index';

const Form = Components.Auth.ForgotPassword;

// ForgotContainer Component
export class ForgotContainer extends Component {

  componentWillMount() {
    this.props.clear();
  }

  componentWillUnmount() {
    this.props.clear();
  }

  render() {
    if (this.props.token)
      return (<Redirect to="/"/>);
    return (<Form
      submit={this.props.recover}
      recoverStatus={this.props.recoverStatus}
    />);
  }

}

ForgotContainer.propTypes = {
  recover: PropTypes.func.isRequired,
  token: PropTypes.string,
  recoverStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  token: state.auth.token,
  recoverStatus: state.auth.recoverStatus
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  recover: (email) => {
    dispatch(Actions.Auth.recover(email))
  },
  clear: () => {
    dispatch(Actions.Auth.clearRecover());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(ForgotContainer);
