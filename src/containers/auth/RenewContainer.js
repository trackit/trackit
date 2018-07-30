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

  componentWillMount() {
    this.props.clear();
  }

  componentWillUnmount() {
    this.props.clear();
  }

  submit = (email, password) => {
    const token = this.props.match.params.token;
    const id = parseInt(this.props.match.params.id, 10);
    this.props.renew(id, password, token);
  };

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
  renewStatus: state.auth.renewStatus
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  renew: (id, password, token) => {
    dispatch(Actions.Auth.renew(id, password, token))
  },
  clear: () => {
    dispatch(Actions.Auth.clearRenew());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(RenewContainer);
