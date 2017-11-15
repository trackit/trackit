import React from 'react';

const required = (value) => {
  if (!value.toString().trim().length)
    return (<div className="alert alert-warning">This field is required</div>);
};

export default {
  required
};
