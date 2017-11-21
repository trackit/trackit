import { shallow } from 'enzyme';
import { noNeg, capitalizeFirstLetter, formatBytes, formatPrice } from '../formatters';

describe('Formatters', () => {

  describe('NoNeg', () => {

    const validInput = 42;
    const negativeInput = -42;
    const stringInput = "aa";
    const stringNegativeInput = "-42";

    it('should return input value', () => {
      expect(noNeg(validInput)).toBe(validInput);
      expect(noNeg(stringInput)).toBe(stringInput);
    });

    it('should return zero', () => {
      expect(noNeg(negativeInput)).toBe(0);
      expect(noNeg(stringNegativeInput)).toBe(0);
    });

  });

  describe('CapitalizeFirstLetter', () => {

    const validInput = "Test";
    const invalidInput = "test";

    it('should return input value', () => {
      expect(capitalizeFirstLetter(validInput)).toBe(validInput);
    });

    it('should return capitalized value', () => {
      expect(capitalizeFirstLetter(invalidInput)).toBe(validInput);
    });

  });

  describe('FormatBytes', () => {

    const validInput = 42;
    const zeroValue = "0 Bytes";

    it('should return formatted value', () => {
      expect(formatBytes(validInput)).toBe(validInput + " Bytes");
    });

    it('should return zero value', () => {
      expect(formatBytes(0)).toBe(zeroValue);
    });

  });

  describe('FormatPrice', () => {

    const validInput = 42.042;
    const formattedValue = "42.04";

    it('should return formatted value', () => {
      const output = shallow(formatPrice(validInput));
      expect(output.length).toBe(1);
      const spans = output.find("span");
      expect(spans.length).toBe(2);
      expect(spans.first().props().children[1]).toBe(formattedValue);
    });

  });

});
