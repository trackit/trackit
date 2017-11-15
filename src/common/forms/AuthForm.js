import React from 'react';
import Validator from 'validator';
import Misc from './MiscForm';

const email = (value) => {
  if (!Validator.isEmail(value))
    return (<div className="alert alert-warning">{value} is not a valid email.</div>);
};

export default {
  required: Misc.required,
  email
};
