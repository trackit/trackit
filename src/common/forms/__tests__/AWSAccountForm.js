import { render } from 'enzyme';
import Validation from '../AWSAccountForm';

describe('Authentication Form Validation', () => {

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
      expect(Validation.s3BucketFormat("test.test")).toBe(undefined);
    });

    it('should render an error', () => {
      const invalidName = render(Validation.s3BucketFormat("s3://test..test"));
      expect(invalidName.hasClass("alert alert-warning")).toBe(true);
      const invalidPrefix = render(Validation.s3BucketFormat("s4://test.test"));
      expect(invalidPrefix.hasClass("alert alert-warning")).toBe(true);
    });

  });

  describe('Path Validation', () => {

    it('should render no error', () => {
      expect(Validation.pathFormat("/test/for/path")).toBe(undefined);
    });

    it('should render an error', () => {
      const path = "/" + Array(1026).join();
      const result = render(Validation.pathFormat(path));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

});
