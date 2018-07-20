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

// Renew Password Form Component
export class RenewPasswordComponent extends Component {

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
    const buttons = (this.props.renewStatus && this.props.renewStatus.status && this.props.renewStatus.value ? null : (
      <div className="clearfix">
        <div>
          <Button
            className="btn btn-primary col-md-5 btn-right"
            type="submit"
          >
            {(this.props.renewStatus && !Object.keys(this.props.renewStatus).length ? (
              (<Spinner className="spinner" name='circle' color='white'/>)
            ) : (
              <div>
                <i className="fa fa-key" />
                &nbsp;
                Reset
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
    ));

    const error = (this.props.renewStatus && this.props.renewStatus.hasOwnProperty("error") ? (
      <div className="alert alert-warning">{this.props.renewStatus.error}</div>
      ): "");

    const success = (this.props.renewStatus && this.props.renewStatus.status && this.props.renewStatus.value ? (
      <div className="alert alert-success">
        <strong>Success : </strong>
        Your new password has been set. You may now <Link to="/login/">Sign in</Link>.
      </div>
    ) : null);

    const form = (
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
    );

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

              <Form
                ref={
                  /* istanbul ignore next */
                  (form) => {this.form = form;}
                }
                onSubmit={this.submit}>
                {success || form}
                {buttons}
              </Form>

            </div>

          </div>

        </div>
      </div>
    );
  }

}

RenewPasswordComponent.propTypes = {
  submit: PropTypes.func.isRequired,
  renewStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

export default withRouter(RenewPasswordComponent);
