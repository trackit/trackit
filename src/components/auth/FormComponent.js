import React, {Component} from 'react';
import { NavLink, withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../common/forms';

import '../../styles/Login.css';
import logo from '../../assets/logo-coloured.png';

const Validation = Validations.Auth;

// Login Form Component
export class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.submit = this.submit.bind(this);
  }

  submit = (e) => {
    e.preventDefault();
    let values = this.form.getValues();
    this.props.submit(values.email, values.password);
  };

  render() {
    const password = (this.props.registration ? (
      <div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <Input
            type="password"
            name="password"
            className="form-control"
            validations={[Validation.required]}
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Confirm your password</label>
          <Input
            type="password"
            name="passwordConfirmation"
            className="form-control"
            validations={[Validation.required, Validation.passwordConfirmation]}
          />
        </div>
      </div>
    ) : (
      <div className="form-group">
        <label htmlFor="password">Password</label>
        <Input
          type="password"
          name="password"
          className="form-control"
          validations={[Validation.required]}
        />
      </div>
    ));

    const buttons = (this.props.registration ? (
      <div>

        <NavLink
          exact to='/login'
          className="btn btn-default col-md-5 btn-left"
        >
          <i className="fa fa-sign-in" />
          &nbsp;
          Sign in
        </NavLink>

        <Button
          className="btn btn-primary col-md-5 btn-right"
          type="submit"
        >
          <i className="fa fa-user-plus" />
          &nbsp;
          Sign up
        </Button>

      </div>
    ) : (
      <div>

        <NavLink
          exact to='/register'
          className="btn btn-default col-md-5 btn-left"
        >
          <i className="fa fa-user-plus" />
          &nbsp;
          Sign up
        </NavLink>

        <Button
          className="btn btn-primary col-md-5 btn-right"
          type="submit"
        >
          <i className="fa fa-sign-in" />
          &nbsp;
          Sign in
        </Button>

      </div>
    ));

    return (
      <div className="login">
        <div className="row">
          <div
            className="col-lg-4 col-lg-offset-4 col-md-4 col-md-offset-4 col-sm-6 col-sm-offset-3 parent"
          >
            <div className="white-box vertCentered">

              <img src={logo} id="logo" alt="TrackIt logo" />

              <hr />
              <Form
                ref={
                  /* istanbul ignore next */
                  (form) => {this.form = form;}
                }
                onSubmit={this.submit}>

                <div className="form-group">
                  <label htmlFor="email">Email address</label>
                  <Input
                    name="email"
                    type="email"
                    className="form-control"
                    validations={[Validation.required, Validation.email]}
                  />
                </div>

                {password}

                {buttons}

              </Form>

            </div>

          </div>

        </div>
      </div>
    );
  }

}

FormComponent.propTypes = {
  submit: PropTypes.func.isRequired,
  registration: PropTypes.bool
};

FormComponent.defaultProps = {
  registration: false
};

export default withRouter(FormComponent);
