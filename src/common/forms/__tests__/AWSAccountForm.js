import { render } from 'enzyme';
import Validation from '../AWSAccountForm';

describe('Authentication Form Validation', () => {

  describe('Required Validation', () => {

    it('should render no error', () => {
      expect(Validation.required("test")).toBe(undefined);
    });

    it('should render an error', () => {
      const result = render(Validation.required(""));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

  describe('Role ARN Validation', () => {

    it('should render no error', () => {
      expect(Validation.roleArnFormat("arn:aws:iam::000000000000:role/path/role")).toBe(undefined);
    });

    it('should render an error', () => {
      const result = render(Validation.roleArnFormat("invalid:arn"));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

});
