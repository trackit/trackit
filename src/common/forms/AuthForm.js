import React from 'react';
import Validator from 'validator';
import Misc from './MiscForm';

const email = (value) => {
  if (!Validator.isEmail(value))
    return (<div className="alert alert-warning">{value} is not a valid email.</div>);
};

const passwordConfirmation = (value, props, components) => {
  if (value !== components.password[0].value)
    return (<div className="alert alert-warning">Passwords are not equal.</div>);
};

export default {
  required: Misc.required,
  email,
  passwordConfirmation
};
