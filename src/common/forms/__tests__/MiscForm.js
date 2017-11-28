import { render } from 'enzyme';
import Validation from '../MiscForm';

describe('Misc Form Validation', () => {

  describe('Required Validation', () => {

    it('should render no error', () => {
      expect(Validation.required("test")).toBe(undefined);
    });

    it('should render an error', () => {
      const result = render(Validation.required(""));
      expect(result.hasClass("alert alert-warning")).toBe(true);
    });

  });

});
