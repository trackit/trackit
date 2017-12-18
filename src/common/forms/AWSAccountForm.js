import React from 'react';
import Misc from "./MiscForm";

const roleArnFormat = (value) => {
const regex = /arn:aws:iam::[\d]{12}:role\/(?:[a-zA-Z0-9+=,.@_-](?:\/[a-zA-Z0-9+=,.@_-])?)+/;
if (!regex.exec(value))
  return (<div className="alert alert-warning">{value} is not a valid role ARN.</div>);
};

const s3BucketFormat = (value) => {
  //s3:\/\/
  const regex = /^((?![^/]{1,61}\.\.[^/]{1,61})[a-z.-]{3,63})$/;
  if (!regex.exec(value))
    return (<div className="alert alert-warning">{value} is not a valid S3 bucket.</div>);
};

const pathFormat = (value) => {
  const regex = /^(?:\/(.{0,1024}))?$/;
  if (!regex.exec(value))
    return (<div className="alert alert-warning">{value} is not a valid path.</div>);
};

export default {
  required: Misc.required,
  roleArnFormat,
  s3BucketFormat,
  pathFormat
};
