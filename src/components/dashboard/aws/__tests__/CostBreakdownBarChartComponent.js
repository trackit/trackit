import React from 'react';
import CostBreakdownBarChart  from '../CostBreakdownBarChartComponent';
import { shallow } from 'enzyme';
import moment from 'moment';
import AWS from '../../../aws';

const Chart = AWS.CostBreakdown.Chart;

const props = {
  id: "id",
  accounts: [],
  values: {
    value: 1,
    otherValue: 2
  },
  getValues: jest.fn(),
  dates: {
    startDate: moment().startOf("month"),
    endDate: moment().endOf("month")
  },
  setDates: jest.fn(),
  filter: "product",
  setFilter: jest.fn(),
  interval: "day",
  setInterval: jest.fn()
};

describe('<CostBreakdownBarChart />', () => {

  it('renders a <CostBreakdownBarChart /> component', () => {
    const wrapper = shallow(<CostBreakdownBarChart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Chart/> component when values are available', () => {
    const wrapper = shallow(<CostBreakdownBarChart {...props}/>);
    const chart = wrapper.find(Chart);
    expect(chart.length).toBe(1);
  });

  it('can get values', () => {
    const wrapper = shallow(<CostBreakdownBarChart {...props}/>);
    expect(props.getValues).not.toHaveBeenCalled();
    wrapper.instance().getValues("id", moment().startOf("month"), moment().endOf("month"), []);
    expect(props.getValues).toHaveBeenCalled();
  });

});
