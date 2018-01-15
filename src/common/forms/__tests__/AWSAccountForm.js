import { render } from 'enzyme';
import Validation from '../AWSAccountForm';

describe('Authentication Form Validation', () => {

  const s3BucketName = "bucket.name";
  const s3BucketPath = "my/path";
  const s3Bucket = `s3://${s3BucketName}/${s3BucketPath}`;

  describe('Role ARN Validation', () => {

    it('should render no error', () => {
      expect(Validation.roleArnFormat("arn:aws:iam::000000000000:role/path/role")).toBe(undefined);
    });

    it('should render an error', () => {
      const result = render(Validation.roleArnFormat("invalid:arn"));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

  describe('S3 Bucket Validation', () => {

    it('should render no error', () => {
      expect(Validation.s3BucketFormat(s3Bucket)).toBe(undefined);
    });

    it('should render an error', () => {
      const invalidName = render(Validation.s3BucketFormat("s3://test..test"));
      expect(invalidName.hasClass("alert alert-warning")).toBe(true);
      const invalidPrefix = render(Validation.s3BucketFormat("s4://test.test"));
      expect(invalidPrefix.hasClass("alert alert-warning")).toBe(true);
    });

  });

  describe('Get S3 Bucket Values', () => {

    it('should return values', () => {
      const result = Validation.getS3BucketValues(s3Bucket);
      expect(result.length).toBe(2);
      expect(result[0]).toBe(s3BucketName);
      expect(result[1]).toBe(s3BucketPath);
    });

    it('should return nothing', () => {
      const result = Validation.getS3BucketValues("");
      expect(result).toBe(null);
    });

  });

});
