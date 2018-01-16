import React from 'react';
import Misc from "./MiscForm";

const roleArnFormat = (value) => {
  const regex = /arn:aws:iam::[\d]{12}:role\/(?:[a-zA-Z0-9+=,.@_-](?:\/[a-zA-Z0-9+=,.@_-])?)+/;
  if (!regex.exec(value))
    return (<div className="alert alert-warning">{value} is not a valid role ARN.</div>);
};

const getAccountIDFromRole = (value) => {
  const regex = /arn:aws:iam::(.*?):role\/(?:[a-zA-Z0-9+=,.@_-](?:\/[a-zA-Z0-9+=,.@_-])?)+/;
  return regex.exec(value)[1];
};

const s3BucketFormat = (value) => {
  const result = getS3BucketValues(value);
  if (!result || result.length !== 2 || !result[0] || !result[0].length || !result[1] || !result[1].length)
    return (<div className="alert alert-warning">{value} is not a valid S3 bucket.</div>);
};

const getS3BucketValues = (value) => {
  // Capturing groups :
  // 1. S3 Bucket name
  // 2. Path
  const regex = /^s3:\/\/((?![^/]{1,61}\.\.[^/]{1,61})[a-z.-]{3,63})(?:\/(.{0,1024}))?$/;
  let result = regex.exec(value);
  if (result && result.length && result.length === 3) {
    result.shift();
    return result;
  }
  return null;
};

export default {
  required: Misc.required,
  roleArnFormat,
  getAccountIDFromRole,
  s3BucketFormat,
  getS3BucketValues
};
