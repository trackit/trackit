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

});
