import React from 'react';
import { CostBreakdownContainer, Chart } from '../CostBreakdownContainer';
import Components from '../../../components';
import Moment from 'moment';
import { shallow } from "enzyme";

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;
const CostBreakdownChart = Components.AWS.CostBreakdown.Chart;

describe('<CostBreakdownContainer />', () => {

  const props = {
    costsValues: {},
    costsDates: {},
    costsInterval: {},
    costsFilter: {},
    getCosts: jest.fn(),
    setCostsDates: jest.fn(),
    setCostsInterval: jest.fn(),
    setCostsFilter: jest.fn(),
  };

  const propsAfterMounting = (id) => {
    let costsDates = {};
    costsDates[id] = {
      startDate: Moment().startOf('month'),
      endDate: Moment().endOf('month'),
    };
    let costsInterval = {};
    costsInterval[id] = "interval";
    let costsFilter = {};
    costsFilter[id] = "filter";
    return {
      ...props,
      costsDates,
      costsInterval,
      costsFilter
    }
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <CostBreakdownContainer /> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Chart/> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    let chart = wrapper.find(Chart);
    expect(chart.length).toBe(0);
    const id = wrapper.state("charts")[0];
    wrapper.setProps(propsAfterMounting(id));
    chart = wrapper.find(Chart);
    expect(chart.length).toBe(1);
  });

  it('can init dates, interval and filter for initial chart when mounting', () => {
    expect(props.setCostsDates).not.toHaveBeenCalled();
    expect(props.setCostsInterval).not.toHaveBeenCalled();
    expect(props.setCostsFilter).not.toHaveBeenCalled();
    shallow(<CostBreakdownContainer {...props}/>);
    expect(props.setCostsDates).toHaveBeenCalled();
    expect(props.setCostsInterval).toHaveBeenCalled();
    expect(props.setCostsFilter).toHaveBeenCalled();
  });

  it('can add and remove chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    const id = wrapper.state("charts")[0];
    wrapper.setProps(propsAfterMounting(id));
    expect(props.setCostsDates).toHaveBeenCalledTimes(1);
    expect(props.setCostsInterval).toHaveBeenCalledTimes(1);
    expect(props.setCostsFilter).toHaveBeenCalledTimes(1);
    expect(wrapper.state("charts").length).toBe(1);
    wrapper.instance().addChart({ preventDefault() {} });
    expect(props.setCostsDates).toHaveBeenCalledTimes(2);
    expect(props.setCostsInterval).toHaveBeenCalledTimes(2);
    expect(props.setCostsFilter).toHaveBeenCalledTimes(2);
    expect(wrapper.state("charts").length).toBe(2);
    wrapper.instance().removeChart(id);
    expect(wrapper.state("charts").length).toBe(1);
  });

});

describe('<Chart />', () => {

  const props = {
    id: "42",
    values: {
      value: 1,
      otherValue: 2
    },
    dates: {
      startDate: Moment().startOf('month'),
      endDate: Moment(),
    },
    interval: "day",
    filter: "product",
    getCosts: jest.fn(),
    setDates: jest.fn(),
    setInterval: jest.fn(),
    setFilter: jest.fn(),
  };

  const propsWithClose = {
    ...props,
    close: jest.fn()
  };

  const updatedDateProps = {
    ...props,
    dates: {
      startDate: Moment().startOf('year'),
      endDate: Moment(),
    },
    getCosts: jest.fn()
  };

  const updatedIntervalProps = {
    ...props,
    interval: "month",
    getCosts: jest.fn()
  };

  const updatedFilterProps = {
    ...props,
    filter: "region",
    getCosts: jest.fn()
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <Chart /> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <TimerangeSelector/> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(1);
  });

  it('renders <Selector/> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

  it('renders <CostBreakdownChart/> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    const chart = wrapper.find(CostBreakdownChart);
    expect(chart.length).toBe(1);
  });

  it('renders a <button/> component when can be closed', () => {
    const wrapper = shallow(<Chart {...propsWithClose}/>);
    const button = wrapper.find("button");
    expect(button.length).toBe(1);
  });

  it('loads costs when mounting', () => {
    expect(props.getCosts).not.toHaveBeenCalled();
    shallow(<Chart {...props}/>);
    expect(props.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when dates are updated', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(updatedDateProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedDateProps);
    expect(updatedDateProps.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when interval is updated', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(updatedIntervalProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedIntervalProps);
    expect(updatedIntervalProps.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when filter is updated', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(updatedFilterProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedFilterProps);
    expect(updatedFilterProps.getCosts).toHaveBeenCalled();
  });

  it('does not reload when dates, interval nor filters are updated', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
  });

  it('cat set dates', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.setDates).not.toHaveBeenCalled();
    wrapper.instance().setDates(Moment().startOf('month'), Moment().endOf('month'));
    expect(props.setDates).toHaveBeenCalledTimes(1);
  });

  it('cat set filter', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.setFilter).not.toHaveBeenCalled();
    wrapper.instance().setFilter("filter");
    expect(props.setFilter).toHaveBeenCalledTimes(1);
  });

  it('cat set interval', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.setInterval).not.toHaveBeenCalled();
    wrapper.instance().setInterval("interval");
    expect(props.setInterval).toHaveBeenCalledTimes(1);
  });

  it('cat close', () => {
    const wrapper = shallow(<Chart {...propsWithClose}/>);
    expect(propsWithClose.close).not.toHaveBeenCalled();
    wrapper.instance().close({ preventDefault() {} });
    expect(propsWithClose.close).toHaveBeenCalledTimes(1);
  });

});
