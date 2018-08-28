import { shallow } from 'enzyme';
import Moment from 'moment';
import {
  noNeg,
  capitalizeFirstLetter,
  formatBytes, formatGigaBytes, formatPrice, formatPercent, formatDate,
  costBreakdown, s3Analytics
} from '../formatters';

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
      expect(formatBytes(validInput)).toBe(validInput + "Bytes");
    });

    it('should return zero value', () => {
      expect(formatBytes(0)).toBe(zeroValue);
    });

  });

  describe('FormatGigaBytes', () => {

    const validInput = 42;
    const zeroValue = "0 Bytes";

    it('should return formatted value', () => {
      expect(formatGigaBytes(validInput)).toBe(validInput + "GB");
    });

    it('should return zero value', () => {
      expect(formatGigaBytes(0)).toBe(zeroValue);
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

  describe('FormatPercent', () => {

    const validPositiveInput = 42.042;
    const formattedPositiveValue = "+42.04";
    const validNegativeInput = -42.042;
    const formattedNegativeValue = "-42.04";

    const redValue = 1;
    const nullValue = 0;
    const greenValue = -1;
    const boldValue = 100;
    const notBoldValue = 99;

    it('should return formatted positive value', () => {
      const output = shallow(formatPercent(validPositiveInput));
      expect(output.length).toBe(1);
      const spans = output.find("span");
      expect(spans.length).toBe(2);
      expect(spans.first().props().children[0]).toBe(formattedPositiveValue);
    });

    it('should return formatted negative value', () => {
      const output = shallow(formatPercent(validNegativeInput));
      expect(output.length).toBe(1);
      const spans = output.find("span");
      expect(spans.length).toBe(2);
      expect(spans.first().props().children[0]).toBe(formattedNegativeValue);
    });

    it('should have red class when value is above 0', () => {
      const output = shallow(formatPercent(redValue));
      expect(output.length).toBe(1);
      let color = output.find(".red-color");
      expect(color.length).toBe(1);
      color = output.find(".green-color");
      expect(color.length).toBe(0);
    });

    it('should have green class when value is below 0', () => {
      const output = shallow(formatPercent(greenValue));
      expect(output.length).toBe(1);
      let color = output.find(".green-color");
      expect(color.length).toBe(1);
      color = output.find(".red-color");
      expect(color.length).toBe(0);
    });

    it('should have no color class when value is 0', () => {
      const output = shallow(formatPercent(nullValue));
      expect(output.length).toBe(1);
      let color = output.find(".green-color");
      expect(color.length).toBe(0);
      color = output.find(".red-color");
      expect(color.length).toBe(0);
    });

    it('should have bold class when value is higher or equal to 100', () => {
      const output = shallow(formatPercent(boldValue));
      expect(output.length).toBe(1);
      let color = output.find(".percent-bold");
      expect(color.length).toBe(1);
    });

    it('should not have bold class when value is lower than 100', () => {
      const output = shallow(formatPercent(notBoldValue));
      expect(output.length).toBe(1);
      let color = output.find(".percent-bold");
      expect(color.length).toBe(0);
    });

  });

  describe('FormatDate', () => {

    const validInput = Moment();
    const formattedValues = {
      year: validInput.format('Y'),
      month: validInput.format('MMM Y'),
      day: validInput.format('MMM Do Y')
    };

    it('should return formatted value', () => {
      Object.keys(formattedValues).forEach((precision) => {
        expect(formatDate(validInput, precision)).toBe(formattedValues[precision]);
      })
    });

  });

  describe('Cost Breakdown', () => {

    describe('TansformProductsBarChart', () => {

      const transformProductsBarChart = costBreakdown.transformProductsBarChart;

      const days = {
        day: {
          day1: 42,
          day2: 21
        }
      };

      const costsByProductPerDay = {
        product: {
          product1: {...days},
          product2: {...days}
        }
      };

      const costsAll = {...days};

      const costsMissingDays = {
        product: {
          product1: {...days},
          product2: {
            day: {
              ...days.day,
              day3: 84
            }
          },
        }
      };

      const costsMissingKeys = {
        product: {
          ...costsByProductPerDay.product,
          "": days
        }
      };

      it('returns an empty array when invalid filter', () => {
        expect(transformProductsBarChart(costsByProductPerDay, "region", "day")).toEqual([]);
      });

      it('returns an empty array when valid filter and invalid interval', () => {
        expect(transformProductsBarChart(costsByProductPerDay, "product", "month")).toEqual([]);
      });

      it('returns an empty array when filter is "all" and invalid interval', () => {
        expect(transformProductsBarChart(costsAll, "all", "month")).toEqual([]);
      });

      it('returns formatted array when valid filter and valid interval', () => {
        const output = [{
          key: "product1",
          values: [["day1", days.day.day1], ["day2", days.day.day2]]
        },{
          key: "product2",
          values: [["day1", days.day.day1], ["day2", days.day.day2]]
        }];
        expect(transformProductsBarChart(costsByProductPerDay, "product", "day")).toEqual(output);
      });

      it('returns formatted array when filter is "all" and valid interval', () => {
        const output = [{
          key: "Total",
          values: [["day1", days.day.day1], ["day2", days.day.day2]]
        }];
        expect(transformProductsBarChart(costsAll, "all", "day")).toEqual(output);
      });

      it('fills missing days', () => {
        const output = [{
          key: "product1",
          values: [["day1", days.day.day1], ["day2", days.day.day2], ["day3", 0]]
        },{
          key: "product2",
          values: [["day1", days.day.day1], ["day2", days.day.day2], ["day3", costsMissingDays.product.product2.day.day3]]
        }];
        expect(transformProductsBarChart(costsMissingDays, "product", "day")).toEqual(output);
      });

      it('fills missing keys', () => {
        const output = [{
          key: "product1",
          values: [["day1", days.day.day1], ["day2", days.day.day2]]
        },{
          key: "product2",
          values: [["day1", days.day.day1], ["day2", days.day.day2]]
        },{
          key: "No product",
          values: [["day1", days.day.day1], ["day2", days.day.day2]]
        }];
        expect(transformProductsBarChart(costsMissingKeys, "product", "day")).toEqual(output);
      });

    });

    describe('TransformProductsPieChart', () => {

      const transformProductsPieChart = costBreakdown.transformProductsPieChart;

      const costsByProduct = {
        product: {
          product1: 42,
          product2: 84
        }
      };

      it('returns an empty array when invalid filter', () => {
        expect(transformProductsPieChart(costsByProduct, "region")).toEqual([]);
      });

      it('returns formatted array when valid filter and valid interval', () => {
        const output = [{
          key: "product1",
          value: costsByProduct.product.product1
        },{
          key: "product2",
          value: costsByProduct.product.product2
        }];
        expect(transformProductsPieChart(costsByProduct, "product")).toEqual(output);
      });

    });

    describe('GetTotalPieChart', () => {

      const getTotalPieChart = costBreakdown.getTotalPieChart;

      const data = [{
        key: "product1",
        value: 42
      },{
        key: "product2",
        value: 84
      }];

      it('returns an empty array when invalid data', () => {
        expect(getTotalPieChart(42)).toEqual(0);
      });

      it('returns total when valid data', () => {
        expect(getTotalPieChart(data)).toEqual((42 +  84));
      });

    });

    describe('TransformCostDifferentiator', () => {

      const transformCostDifferentiator = costBreakdown.transformCostDifferentiator;

      it('returns formatted data', () => {
        const data = {
          product: [{
            Date: "date1",
            Cost: 42,
            PercentVariation: 4.2
          }, {
            Date: "date2",
            Cost: 84,
            PercentVariation: 8.4
          }],
          product2: [{
            Date: "date1",
            Cost: 21,
            PercentVariation: 2.1
          }]
        };

        const output = {
          dates: ["date1", "date2"],
          values: [{
            key: "product",
            date1: {cost: 42, variation: 4.2},
            date2: {cost: 84, variation: 8.4}
          }, {
            key: "product2",
            date1: {cost: 21, variation: 2.1}
          }],
          total: {
            key: "Total",
            date1: {
              cost: 42 + 21,
              variation: 0
            },
            date2: {
              cost: 84,
              variation: (84 - 42 - 21) / (42 + 21) * 100
            }
          }
        };
        expect(transformCostDifferentiator(data)).toEqual(output);
      });

      it('formats variation when cost are almost null', () => {
        const data = {
          product: [{
            Date: "date1",
            Cost: 0,
            PercentVariation: 0
          }, {
            Date: "date2",
            Cost: 0,
            PercentVariation: 8.4
          }],
          product2: [{
            Date: "date1",
            Cost: 21,
            PercentVariation: 2.1
          }]
        };

        const output = {
          dates: ["date1", "date2"],
          values: [{
            key: "product",
            date1: {cost: 0, variation: 0},
            date2: {cost: 0, variation: 0}
          }, {
            key: "product2",
            date1: {cost: 21, variation: 2.1}
          }],
          total: {
            key: "Total",
            date1: {
              cost: 21,
              variation: 0
            },
            date2: {
              cost: 0,
              variation: (0 - 21) / 21 * 100
            }
          }
        };
        expect(transformCostDifferentiator(data)).toEqual(output);
      });

    });

  });

  describe('S3 Analytics', () => {

    describe('TransformBuckets', () => {

      const transformBuckets = s3Analytics.transformBuckets;

      const buckets = {
        bucket1: {
          BandwidthCost: 42,
          StorageCost: 84
        },
        bucket2: {
          BandwidthCost: 21,
          StorageCost: 42
        }
      };

      it('returns formatted array when valid data', () => {
        const output = [{
          key: "bucket1",
          values: [
            ["Bandwidth", 42],
            ["Storage", 84],
          ]
        },{
          key: "bucket2",
          values: [
            ["Bandwidth", 21],
            ["Storage", 42],
          ]
        }];
        expect(transformBuckets(buckets)).toEqual(output);
      });

    });

    describe('TransformBandwidthPieChart', () => {

      const transformBandwidthPieChart = s3Analytics.transformBandwidthPieChart;

      const buckets = {
        bucket1: {
          BandwidthCost: 42,
          StorageCost: 84
        },
        bucket2: {
          BandwidthCost: 21,
          StorageCost: 42
        }
      };

      it('returns formatted array when valid data', () => {
        const output = [{
          key: "bucket1",
          value: 42
        },{
          key: "bucket2",
          value: 21
        }];
        expect(transformBandwidthPieChart(buckets)).toEqual(output);
      });

    });

    describe('TransformStoragePieChart', () => {

      const transformStoragePieChart = s3Analytics.transformStoragePieChart;

      const buckets = {
        bucket1: {
          BandwidthCost: 42,
          StorageCost: 84
        },
        bucket2: {
          BandwidthCost: 21,
          StorageCost: 42
        }
      };

      it('returns formatted array when valid data', () => {
        const output = [{
          key: "bucket1",
          value: 84
        },{
          key: "bucket2",
          value: 42
        }];
        expect(transformStoragePieChart(buckets)).toEqual(output);
      });

    });

    describe('GetTotalPieChart', () => {

      const getTotalPieChart = s3Analytics.getTotalPieChart;

      const data = [{
        key: "bucket1",
        value: 42
      },{
        key: "bucket2",
        value: 84
      }];

      it('returns an empty array when invalid data', () => {
        expect(getTotalPieChart(42)).toEqual(0);
      });

      it('returns total when valid data', () => {
        expect(getTotalPieChart(data)).toEqual((42 +  84));
      });

    });

  });

});
