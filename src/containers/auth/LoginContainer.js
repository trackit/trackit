import React, {Component} from 'react';
import {connect} from 'react-redux';
import {Redirect} from "react-router-dom";
// import PropTypes from 'prop-types';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../common/forms';

import Actions from '../../actions/index';

import '../../styles/Login.css';
import logo from '../../assets/logo-coloured.png';

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
      return (<Redirect to="/"/>);
    return (
      <div className="login">
        <div className="row">
          <div
            className="col-lg-4 col-lg-offset-4 col-md-4 col-md-offset-4 col-sm-6 col-sm-offset-3 parent"
          >
            <div className="white-box vertCentered">

              <img src={logo} id="logo" alt="TrackIt logo" />

              <hr />
              <Form ref={form => {
                this.form = form;
              }} onSubmit={this.submit.bind(this)}>

                <div className="form-group">
                  <label htmlFor="email">Email address</label>
                  <Input
                    name="email"
                    type="email"
                    className="form-control"
                    validations={[Validation.required, Validation.email]}
                  />
                </div>
                <div className="form-group">
                  <label htmlFor="password">Password</label>
                  <Input
                    type="password"
                    name="password"
                    className="form-control"
                    validations={[Validation.required]}
                  />
                </div>

                <div>
                  <Button
                    className="btn btn-primary btn-block"
                    type="submit"
                  >
                    <i className="fa fa-sign-in" />
                    &nbsp;
                    Sign in
                  </Button>
                </div>
              </Form>

            </div>

          </div>

        </div>
      </div>
    );
  }

}

LoginContainer.propTypes = {};

const mapStateToProps = (state) => ({token: state.auth.token});

const mapDispatchToProps = (dispatch) => ({
  login: (email, password) => {
    dispatch(Actions.Auth.login(email, password))
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(LoginContainer);
