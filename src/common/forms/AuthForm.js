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

export default {
  required,
  email
};
