import React from 'react';
import Validator from 'validator';

const required = (value) => {
  if (!value.toString().trim().length)
    return (<span className="error">Value is required</span>);
};

const email = (value) => {
  if (!Validator.isEmail(value))
    return (<span className="error">{value} is not a valid email.</span>);
};

const lt = (value, props) => {
  if (value.toString().trim().length > props.maxLength)
    return (<span className="error">The value exceeded {props.maxLength} symbols.</span>);
};

const password = (value, props, components) => {
  if (value !== components['confirm'][0].value)
    return (<span className="error">Passwords are not equal.</span>);
};

export default {
  required,
  email,
  lt,
  password
};
