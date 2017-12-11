import React, {Component} from 'react';
import {connect} from 'react-redux';
import {Redirect} from "react-router-dom";
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions/index';

const Form = Components.Auth.Form;

// RegisterContainer Component
export class RegisterContainer extends Component {

  componentWillUnmount() {
    this.props.clear();
  }

  render() {
    if (this.props.registration && this.props.registration.status)
      return (<Redirect to="/login"/>);
    return (<Form submit={this.props.register} registration/>);
  }

}

RegisterContainer.propTypes = {
  register: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({registration: state.auth.registration});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  register: (email, password) => {
    dispatch(Actions.Auth.register(email, password))
  },
  clear: () => {
    dispatch(Actions.Auth.clearRegister());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(RegisterContainer);
