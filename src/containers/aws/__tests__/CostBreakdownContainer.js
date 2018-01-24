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
    accounts: [],
    charts: [],
    costsValues: {},
    costsDates: {},
    costsInterval: {},
    costsFilter: {},
    addChart: jest.fn(),
    removeChart: jest.fn(),
    getCosts: jest.fn(),
    setCostsDates: jest.fn(),
    setCostsInterval: jest.fn(),
    setCostsFilter: jest.fn(),
    resetCostsDates: jest.fn(),
    resetCostsInterval: jest.fn(),
    resetCostsFilter: jest.fn(),
  };

  const propsWithUniqueChart = {
    ...props,
    charts: ["id"],
    costsDates: {
      id: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      }
    },
    costsInterval: {
      id: "interval"
    },
    costsFilter: {
      id: "filter"
    }
  };

  const propsWithCharts = {
    ...propsWithUniqueChart,
    charts: ["id", "id2"]
  };

  const propsWithValidCharts = {
    ...propsWithCharts,
    costsDates: {
      id: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      },
      id2: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      }
    },
    costsInterval: {
      id: "interval",
      id2: "interval"
    },
    costsFilter: {
      id: "filter",
      id2: "filter"
    }
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <CostBreakdownContainer /> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Chart/> component if data is available', () => {
    const wrapperUnique = shallow(<CostBreakdownContainer {...propsWithUniqueChart}/>);
    let charts = wrapperUnique.find(Chart);
    expect(propsWithUniqueChart.charts.length).toBe(1);
    expect(Object.keys(propsWithUniqueChart.costsDates).length).toBe(1);
    expect(charts.length).toBe(1);
    const wrapperOne = shallow(<CostBreakdownContainer {...propsWithCharts}/>);
    charts = wrapperOne.find(Chart);
    expect(propsWithCharts.charts.length).toBe(2);
    expect(Object.keys(propsWithCharts.costsDates).length).toBe(1);
    expect(charts.length).toBe(1);
    const wrapperTwo = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    charts = wrapperTwo.find(Chart);
    expect(propsWithValidCharts.charts.length).toBe(2);
    expect(Object.keys(propsWithValidCharts.costsDates).length).toBe(2);
    expect(charts.length).toBe(2);
  });

  it('generates default <Chart/> component if no chart available', () => {
    shallow(<CostBreakdownContainer {...props}/>);
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can add chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithUniqueChart}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().addChart({ preventDefault() {} });
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can reset charts', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithUniqueChart}/>);
    expect(props.resetCostsDates).not.toHaveBeenCalled();
    expect(props.resetCostsInterval).not.toHaveBeenCalled();
    expect(props.resetCostsFilter).not.toHaveBeenCalled();
    wrapper.instance().resetCharts({ preventDefault() {} });
    expect(props.resetCostsDates).toHaveBeenCalled();
    expect(props.resetCostsInterval).toHaveBeenCalled();
    expect(props.resetCostsFilter).toHaveBeenCalled();
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

  it('can set dates', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.setDates).not.toHaveBeenCalled();
    wrapper.instance().setDates(Moment().startOf('month'), Moment().endOf('month'));
    expect(props.setDates).toHaveBeenCalledTimes(1);
  });

  it('can set filter', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.setFilter).not.toHaveBeenCalled();
    wrapper.instance().setFilter("filter");
    expect(props.setFilter).toHaveBeenCalledTimes(1);
  });

  it('can set interval', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.setInterval).not.toHaveBeenCalled();
    wrapper.instance().setInterval("interval");
    expect(props.setInterval).toHaveBeenCalledTimes(1);
  });

  it('can close', () => {
    const wrapper = shallow(<Chart {...propsWithClose}/>);
    expect(propsWithClose.close).not.toHaveBeenCalled();
    wrapper.instance().close({ preventDefault() {} });
    expect(propsWithClose.close).toHaveBeenCalledTimes(1);
  });

});
