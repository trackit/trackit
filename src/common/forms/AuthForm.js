import React from 'react';
import Validator from 'validator';

const required = (value) => {
  if (!value.toString().trim().length)
    return (<div className="alert alert-warning">This field is required</div>);
};

const email = (value) => {
  if (!Validator.isEmail(value))
    return (<div className="alert alert-warning">{value} is not a valid email.</div>);
};

const lt = (value, props) => {
  if (value.toString().trim().length > props.maxLength)
    return (<div className="alert alert-warning">The value exceeded {props.maxLength} symbols.</div>);
};

const password = (value, props, components) => {
  if (value !== components['confirm'][0].value)
    return (<div className="alert alert-warning">Passwords are not equal.</div>);
};

export default {
  required,
  email,
  lt,
  password
};
