import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Redirect } from "react-router-dom";
// import PropTypes from 'prop-types';

import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../common/forms';

import Actions from '../../actions/index';

const Validation = Validations.Auth;

// LoginContainer Component
class LoginContainer extends Component {

  submit(e) {
    e.preventDefault();
    let values = this.form.getValues();
    this.props.login(values.email, values.password);
  }

  render() {
    if (this.props.token !== null)
      return (
        <Redirect to="/"/>
      );
    return (
      <div className="container-fluid">
        <Form
          ref={form => { this.form = form; }}
          onSubmit={this.submit.bind(this)}
        >
          <h3>Login</h3>
          <div>
            <label>
              Email
              <Input
                name='email'
                validations={[Validation.required, Validation.email]}
              />
            </label>
          </div>
          <div>
            <label>
              Password
              <Input
                type='password'
                name='password'
                validations={[Validation.required]}
              />
            </label>
          </div>
          <div>
            <Button>Login</Button>
          </div>
        </Form>
      </div>
    );
  }

}

LoginContainer.propTypes = {};

const mapStateToProps = (state) => ({
  token: state.auth.token
});

const mapDispatchToProps = (dispatch) => ({
  login: (email, password) => {
    dispatch(Actions.Auth.login(email, password))
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(LoginContainer);
