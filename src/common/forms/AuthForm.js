import React from 'react';
import Validator from 'validator';
import Misc from './MiscForm';

const email = (value) => {
  if (!Validator.isEmail(value))
    return (<div className="alert alert-warning">{value} is not a valid email.</div>);
};

const passwordConfirmation = (value, props, components) => {
  if (value !== components.passwordConfirmation[0].value)
    return <span className="alert alert-warning">Passwords are not equal.</span>
};

export default {
  required: Misc.required,
  email,
  passwordConfirmation
};
