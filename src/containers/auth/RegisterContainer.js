import React, {Component} from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions/index';

const Form = Components.Auth.Form;

// RegisterContainer Component
export class RegisterContainer extends Component {

  render() {
    const awstoken = (this.props.match.params.hasOwnProperty("awstoken") ? this.props.match.params.awstoken : "");
    return (
      <Form
        awsToken={decodeURIComponent(awstoken)}
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
  register: (email, password, awsToken) => {
    dispatch(Actions.Auth.register(email, password, awsToken))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(RegisterContainer);
