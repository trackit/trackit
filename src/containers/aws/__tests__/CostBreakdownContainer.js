import React from 'react';
import { CostBreakdownContainer } from '../CostBreakdownContainer';
import Components from '../../../components';
import Moment from 'moment';
import NVD3Chart from 'react-nvd3';
import { shallow } from "enzyme";

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;

const props = {
  costsValues: {
    value: 1,
    otherValue: 2
  },
  costsDates: {
    startDate: Moment().startOf('month'),
    endDate: Moment(),
  },
  costsInterval: "day",
  costsFilter: "product",
  getCosts: jest.fn(),
  setCostsDates: jest.fn(),
  setCostsInterval: jest.fn(),
  setCostsFilter: jest.fn(),
};

const updatedDateProps = {
  ...props,
  costsDates: {
    startDate: Moment().startOf('year'),
    endDate: Moment(),
  },
  getCosts: jest.fn()
};

const updatedIntervalProps = {
  ...props,
  costsInterval: "month",
  getCosts: jest.fn()
};

const updatedFilterProps = {
  ...props,
  costsFilter: "region",
  getCosts: jest.fn()
};

const propsWithoutCosts = {
  ...props,
  costsValues: null
};

describe('<CostBreakdownContainer />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <CostBreakdownContainer /> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <TimerangeSelector/> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(1);
  });

  it('renders <Selector/> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

  it('renders <NVD3Chart/>  when costs are available', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders no <NVD3Chart/>  when costs are unavailable', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithoutCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('loads costs when mounting', () => {
    expect(props.getCosts).not.toHaveBeenCalled();
    shallow(<CostBreakdownContainer {...props}/>);
    expect(props.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when dates are updated', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(updatedDateProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedDateProps);
    expect(updatedDateProps.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when interval is updated', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(updatedIntervalProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedIntervalProps);
    expect(updatedIntervalProps.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when filter is updated', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(updatedFilterProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedFilterProps);
    expect(updatedFilterProps.getCosts).toHaveBeenCalled();
  });

  it('does not reload when dates, interval nor filters are updated', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
  });

  describe('Costs scheme transformation', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    const transformProducts = wrapper.instance().transformProducts;

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
      expect(transformProducts(costsByProductPerDay, "region", "day")).toEqual([]);
    });

    it('returns an empty array when valid filter and invalid interval', () => {
      expect(transformProducts(costsByProductPerDay, "product", "month")).toEqual([]);
    });

    it('returns an empty array when filter is "all" and invalid interval', () => {
      expect(transformProducts(costsAll, "all", "month")).toEqual([]);
    });

    it('returns formatted array when valid filter and valid interval', () => {
      const output = [{
        key: "product1",
        values: [["day1", days.day.day1], ["day2", days.day.day2]]
      },{
        key: "product2",
        values: [["day1", days.day.day1], ["day2", days.day.day2]]
      }];
      expect(transformProducts(costsByProductPerDay, "product", "day")).toEqual(output);
    });

    it('returns formatted array when filter is "all" and valid interval', () => {
      const output = [{
        key: "Total",
        values: [["day1", days.day.day1], ["day2", days.day.day2]]
      }];
      expect(transformProducts(costsAll, "all", "day")).toEqual(output);
    });

    it('fills missing days', () => {
      const output = [{
        key: "product1",
        values: [["day1", days.day.day1], ["day2", days.day.day2], ["day3", 0]]
      },{
        key: "product2",
        values: [["day1", days.day.day1], ["day2", days.day.day2], ["day3", costsMissingDays.product.product2.day.day3]]
      }];
      expect(transformProducts(costsMissingDays, "product", "day")).toEqual(output);
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
      expect(transformProducts(costsMissingKeys, "product", "day")).toEqual(output);
    });

  });


});
