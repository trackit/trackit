import React, {Component} from 'react';
import {connect} from 'react-redux';
import {Redirect} from "react-router-dom";
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions/index';

const Form = Components.Auth.Form;

// RegisterContainer Component
export class RegisterContainer extends Component {

  render() {
    if (this.props.registrationStatus && this.props.registrationStatus.status)
      return (<Redirect to="/login"/>);
    return (<Form
      submit={this.props.register}
      registration
      registrationStatus={this.props.registrationStatus}
    />);
  }

}

RegisterContainer.propTypes = {
  register: PropTypes.func.isRequired,
  registrationStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({registrationStatus: state.auth.registration});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  register: (email, password) => {
    dispatch(Actions.Auth.register(email, password))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(RegisterContainer);
