import React from 'react';
import Misc from "./MiscForm";

const roleArnFormat = (value) => {
const regex = /arn:aws:iam::[\d]{12}:role\/(?:[a-zA-Z0-9+=,.@_-](?:\/[a-zA-Z0-9+=,.@_-])?)+/;
if (!regex.exec(value))
  return (<div className="alert alert-warning">{value} is not a valid role ARN.</div>);
};

const s3BucketFormat = (value) => {
  const result = getS3BucketValues(value);
  if (!result || result.length !== 2)
    return (<div className="alert alert-warning">{value} is not a valid S3 bucket.</div>);
};

const getS3BucketValues = (value) => {
  // Capturing groups :
  // 1. S3 Bucket name
  // 2. Path
  const regex = /^s3:\/\/((?![^/]{1,61}\.\.[^/]{1,61})[a-z.-]{3,63})(?:\/(.{0,1024}))?$/;
  let result = regex.exec(value);
  if (result && result.length)
    result.shift();
  return result;
};

export default {
  required: Misc.required,
  roleArnFormat,
  s3BucketFormat,
  getS3BucketValues
};
