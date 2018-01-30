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
    initCharts: jest.fn(),
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

  const propsWithInvalidCharts = {
    ...props,
    charts: ["id", "id2"]
  };

  const propsWithValidCharts = {
    ...propsWithInvalidCharts,
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

  const propsWithThreeCharts = {
    ...props,
    charts: ["id", "id2", "id3"],
    costsDates: {
      id: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      },
      id2: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      },
      id3: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      }
    },
    costsInterval: {
      id: "interval",
      id2: "interval",
      id3: "interval"
    },
    costsFilter: {
      id: "filter",
      id2: "filter",
      id3: "filter"
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
    const wrapperInvalid = shallow(<CostBreakdownContainer {...propsWithInvalidCharts}/>);
    let charts = wrapperInvalid.find(Chart);
    expect(propsWithInvalidCharts.charts.length).toBe(2);
    expect(Object.keys(propsWithInvalidCharts.costsDates).length).toBe(0);
    expect(charts.length).toBe(0);
    const wrapperValid = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    charts = wrapperValid.find(Chart);
    expect(propsWithValidCharts.charts.length).toBe(2);
    expect(Object.keys(propsWithValidCharts.costsDates).length).toBe(2);
    expect(charts.length).toBe(2);
    const wrapperThree = shallow(<CostBreakdownContainer {...propsWithThreeCharts}/>);
    charts = wrapperThree.find(Chart);
    expect(propsWithThreeCharts.charts.length).toBe(3);
    expect(Object.keys(propsWithThreeCharts.costsDates).length).toBe(3);
    expect(charts.length).toBe(3);
  });

  it('generates default <Chart/> component if no chart available', () => {
    expect(props.initCharts).not.toHaveBeenCalled();
    shallow(<CostBreakdownContainer {...props}/>);
    expect(props.initCharts).toHaveBeenCalled();
  });

  it('can add chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().addChart({ preventDefault() {} });
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can reset charts', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.removeChart).not.toHaveBeenCalled();
    wrapper.instance().resetCharts({ preventDefault() {} });
    expect(props.removeChart).toHaveBeenCalled();
  });

  it('adds a chart when there is no chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.initCharts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.initCharts).toHaveBeenCalled();
  });

  it('does not add a chart when there is available charts', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsWithValidCharts);
    expect(props.addChart).not.toHaveBeenCalled();
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
