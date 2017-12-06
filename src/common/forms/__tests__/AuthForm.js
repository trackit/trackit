import { render } from 'enzyme';
import Validation from '../AuthForm';

describe('Authentication Form Validation', () => {

  describe('Email Validation', () => {

    it('should render no error', () => {
      expect(Validation.email("test@test.test")).toBe(undefined);
    });

    it('should render an error', () => {
      const result = render(Validation.email("test"));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

  describe('Password Confirmation', () => {

    const validConfirmation = {
      passwordConfirmation: [{
        value: "password"
      }]
    };

    const invalidConfirmation = {
      passwordConfirmation: [{
        value: "passw0rd"
      }]
    };

    it('should render no error', () => {
      expect(Validation.passwordConfirmation("password", null, validConfirmation)).toBe(undefined);
    });

    it('should render an error', () => {
      const result = render(Validation.passwordConfirmation("password", null, invalidConfirmation));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

});
