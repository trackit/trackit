import React from 'react';
import { CostBreakdownContainer } from '../CostBreakdownContainer';
import Components from '../../../components';
import Moment from 'moment';
import { shallow } from "enzyme";

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;
const Chart = Components.AWS.CostBreakdown.Chart;

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

  it('renders <Chart/> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    const chart = wrapper.find(Chart);
    expect(chart.length).toBe(1);
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

});
