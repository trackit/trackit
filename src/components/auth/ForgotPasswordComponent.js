import React, {Component} from 'react';
import { Link, withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../common/forms';

import '../../styles/Login.css';
import logo from '../../assets/logo-coloured.png';

const Validation = Validations.Auth;

// Forgot Password Form Component
export class ForgotPasswordComponent extends Component {

  constructor(props) {
    super(props);
    this.submit = this.submit.bind(this);
  }

  submit = (e) => {
    e.preventDefault();
    let values = this.form.getValues();
    this.props.submit(values.email);
  };

  render() {
    const buttons = (
      <div className="clearfix">
        <div>
          <Button
            className="btn btn-primary col-md-5 btn-right"
            type="submit"
          >
            {(this.props.recoverStatus && !Object.keys(this.props.recoverStatus).length ? (
              (<Spinner className="spinner" name='circle' color='white'/>)
            ) : (
              <div>
                <i className="fa fa-envelope" />
                &nbsp;
                Send
              </div>
            ))}
          </Button>
        </div>
        <Link
          to="/login"
        >
          Return to Login
        </Link>


      </div>
    );

    const error = (this.props.recoverStatus && this.props.recoverStatus.hasOwnProperty("error") ? (
      <div className="alert alert-warning">{this.props.recoverStatus.error}</div>
      ): "");

    const success = (this.props.recoverStatus && this.props.recoverStatus.status && this.props.recoverStatus.value ? (
      <div className="alert alert-success">
        <strong>Success : </strong>
        An email has been sent to you with a link to setup a new password.
      </div>
    ) : null);

    return (
      <div className="login">
        <div className="row">
          <div
            className="col-lg-4 col-lg-offset-4 col-md-4 col-md-offset-4 col-sm-6 col-sm-offset-3 parent"
          >
            <div className="white-box vertCentered">

              <img src={logo} id="logo" alt="TrackIt logo" />

              <hr />

              {error}
              {success}



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

                {buttons}

              </Form>

            </div>

          </div>

        </div>
      </div>
    );
  }

}

ForgotPasswordComponent.propTypes = {
  submit: PropTypes.func.isRequired,
  recoverStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

export default withRouter(ForgotPasswordComponent);
