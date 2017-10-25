import React from 'react';

const required = (value) => {
  if (!value.toString().trim().length)
    return (<div className="alert alert-warning">This field is required</div>);
};

const roleArnFormat = (value) => {
  const regex = /arn:aws:iam::[\d]{12}:role\/(?:[a-zA-Z0-9+=,.@-_](?:\/[a-zA-Z0-9+=,.@-_])?)+/;
  if (!regex.exec(value))
    return (<div className="alert alert-warning">{value} is not a valid role ARN.</div>);
};
export default {
  required,
  roleArnFormat
};
