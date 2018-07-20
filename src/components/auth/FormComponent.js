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

// Login Form Component
export class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      email: (this.props.match && this.props.match.params && this.props.match.params.prefill) ? decodeURIComponent(this.props.match.params.prefill) : '',
      password: '',
      passwordConfirmation: '',
    }
    this.submit = this.submit.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
  }

  submit = (e) => {
    e.preventDefault();
    this.props.submit(this.state.email, this.state.password);
  };

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({
      [name]: value
    });
  }


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
            value={this.state.password}
            onChange={this.handleInputChange}
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Confirm your password</label>
          <Input
            type="password"
            name="passwordConfirmation"
            className="form-control"
            validations={[Validation.required, Validation.passwordConfirmation]}
            value={this.state.passwordConfirmation}
            onChange={this.handleInputChange}
          />
        </div>
      </div>
    ) : (
      <div className="form-group">
        <label htmlFor="password">Password</label>
        <Link
          to="/forgot"
          className="pull-right"
        >
          Forgot password ?
        </Link>
        <Input
          type="password"
          name="password"
          className="form-control"
          validations={[Validation.required]}
          value={this.state.password}
          onChange={this.handleInputChange}
        />
      </div>
    ));

    const buttons = (this.props.registration ? (
      <div className="clearfix">
        <div>
          <Button
            className="btn btn-primary btn-block"
            type="submit"
          >
            {(this.props.registrationStatus && !Object.keys(this.props.registrationStatus).length ? (
              (<Spinner className="spinner" name='circle' color='white'/>)
            ) : (
              <div>
                <i className="fa fa-user-plus" />
                &nbsp;
                Sign up
              </div>
            ))}
          </Button>
        </div>
        <br />
        <Link
          to="/login"
        >
          Already have an account ? Sign in here.
        </Link>


      </div>
    ) : (
      <div className="clearfix">
        <div>
          <Button
            className="btn btn-primary btn-block"
            type="submit"
          >
            {(this.props.loginStatus && !Object.keys(this.props.loginStatus).length ? (
              (<Spinner className="spinner" name='circle' color='white'/>)
            ) : (
              <div>
                <i className="fa fa-sign-in" />
                &nbsp;
                Sign in
              </div>
            ))}
          </Button>
        </div>
        <br />
        <Link
          to="/register"
        >
          Don't have an account ? Register here.
        </Link>
      </div>
    ));

    let error;
    if (!this.props.registration) {
      error = (this.props.loginStatus && this.props.loginStatus.hasOwnProperty("error") ? (
        <div className="alert alert-warning">{this.props.loginStatus.error}</div>
      ): "");
    } else {
      error = (this.props.registrationStatus && this.props.registrationStatus.hasOwnProperty("error") ? (
        <div className="alert alert-warning">{this.props.registrationStatus.error}</div>
      ): "");
    }

    let success;
    if (this.props.registrationStatus && this.props.registrationStatus.status) {
      success = <div className="alert alert-success">
        <strong>Success : </strong>
        Your registration was successful. You may now <Link to={`/login/${encodeURIComponent(this.state.email)}`}>Sign in</Link>.
      </div>;
    }

    return (
      <div className="login">
        <div className="row">
          <div
            className="col-lg-4 col-lg-offset-4 col-md-4 col-md-offset-4 col-sm-6 col-sm-offset-3 parent"
          >
            <div className="white-box vertCentered">

              <img src={logo} id="logo" alt="TrackIt logo" />

              <hr />
              <h3 style={{textAlign: 'center'}}>
                {this.props.registration ? 'Register' : 'Sign in'}
              </h3>

              {error}
              {success}



              <Form
                onSubmit={this.submit}>

                {
                  !(this.props.registrationStatus && this.props.registrationStatus.status) &&
                  (
                    <div className="form-group">
                      <label htmlFor="email">Email address</label>
                      <Input
                        name="email"
                        type="email"
                        className="form-control"
                        validations={[Validation.required, Validation.email]}
                        value={this.state.email}
                        onChange={this.handleInputChange}
                      />
                    </div>
                  )
                }

                {!(this.props.registrationStatus && this.props.registrationStatus.status) && password}

                {!(this.props.registrationStatus && this.props.registrationStatus.status) && buttons}

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
  registration: PropTypes.bool,
  loginStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  }),
  registrationStatus: PropTypes.shape({
    status: PropTypes.bool,
    error: PropTypes.string
  })
};

FormComponent.defaultProps = {
  registration: false
};

export default withRouter(FormComponent);
