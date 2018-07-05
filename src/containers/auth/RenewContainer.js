import React, {Component} from 'react';
import {connect} from 'react-redux';
import {Redirect} from "react-router-dom";
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions/index';

const Form = Components.Auth.RenewPassword;

// RenewContainer Component
export class RenewContainer extends Component {

  constructor(props) {
    super(props);
    this.submit = this.submit.bind(this);
  }

  submit = (email, password) => {
    const token = this.props.match.params.token;
    this.props.renew(email, password, token);
  };

  componentWillUnmount() {
    this.props.clear();
  }

  render() {
    if (this.props.token)
      return (<Redirect to="/"/>);
    return (<Form
      submit={this.submit}
      renewStatus={this.props.renewStatus}
    />);
  }

}

RenewContainer.propTypes = {
  renew: PropTypes.func.isRequired,
  token: PropTypes.string,
  renewStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  token: state.auth.token,
  renewStatus: state.auth.recoverStatus
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  renew: (email, password, token) => {
    dispatch(Actions.Auth.renew(email, password, token))
  },
  clear: () => {
    dispatch(Actions.Auth.clearRenew());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(RenewContainer);
