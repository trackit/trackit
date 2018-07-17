import React from 'react';
import Chart, { Header } from '../ChartComponent';
import BarChart from '../BarChartComponent';
import PieChart from '../PieChartComponent';
import DiffChart from '../DifferentiatorChartComponent';
import Misc from '../../../misc';
import Moment from 'moment';
import Spinner from 'react-spinkit';
import { shallow } from 'enzyme';

const TimerangeSelector = Misc.TimerangeSelector;
const IntervalNavigator = Misc.IntervalNavigator;
const Selector = Misc.Selector;

const props = {
  id: "42",
  type: "bar",
  values: {
    status: true,
    values: {
      value: 1,
      otherValue: 2
    }
  },
  dates: {
    startDate: Moment().startOf('month'),
    endDate: Moment(),
  },
  accounts: [],
  interval: "day",
  filter: "product",
  getCosts: jest.fn(),
  setDates: jest.fn(),
  setInterval: jest.fn(),
  setFilter: jest.fn(),
};

const propsPie = {
  ...props,
  type: "pie"
};

const propsDiff = {
  ...props,
  type: "diff"
};

const propsWithError = {
  ...props,
  values: {
    status: true,
    error: Error()
  }
};

const propsWithoutData = {
  ...props,
  values: {
    status: false
  }
};

const propsPieWithoutData = {
  ...propsPie,
  values: {
    status: false
  }
};

const propsDiffWithoutData = {
  ...propsDiff,
  values: {
    status: false
  }
};

const propsWithClose = {
  ...props,
  close: jest.fn()
};

const propsWithTable = {
  ...propsPie,
  table: true,
};

const propsWithTableShow = {
  ...propsWithTable,
  tableStatus: false
};

const propsWithTableHide = {
  ...propsWithTable,
  tableStatus: true
};

const propsWithTableAndClose = {
  ...propsWithTable,
  close: jest.fn()
};

const updatedAccountsProps = {
  ...propsDiff,
  accounts: ["account"],
  getCosts: jest.fn()
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

const updatedFilterPropsPie = {
  ...propsPie,
  filter: "region",
  getCosts: jest.fn()
};

const propsNoDates = {
  ...props,
  setDates: undefined
};

const propsWithoutIcon = {
  ...props,
  icon: false
};

describe('<Chart />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <Chart /> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Header/> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    const header = wrapper.find(Header);
    expect(header.length).toBe(1);
  });

  it('renders <BarChart/> component when BarChart is selected and data is available', () => {
    const wrapper = shallow(<Chart {...props}/>);
    const chart = wrapper.find(BarChart);
    expect(chart.length).toBe(1);
  });

  it('renders no <BarChart/> component when BarChart is selected and data is unavailable', () => {
    const wrapper = shallow(<Chart {...propsWithoutData}/>);
    const chart = wrapper.find(BarChart);
    expect(chart.length).toBe(0);
  });

  it('renders <PieChart/> component when PieChart is selected and data is available', () => {
    const wrapper = shallow(<Chart {...propsPie}/>);
    const chart = wrapper.find(PieChart);
    expect(chart.length).toBe(1);
  });

  it('renders no <PieChart/> component when PieChart is selected and data is unavailable', () => {
    const wrapper = shallow(<Chart {...propsPieWithoutData}/>);
    const chart = wrapper.find(PieChart);
    expect(chart.length).toBe(0);
  });

  it('renders <DiffChart/> component when DiffChart is selected and data is available', () => {
    const wrapper = shallow(<Chart {...propsDiff}/>);
    const chart = wrapper.find(DiffChart);
    expect(chart.length).toBe(1);
  });

  it('renders no <DiffChart/> component when DiffChart is selected and data is unavailable', () => {
    const wrapper = shallow(<Chart {...propsDiffWithoutData}/>);
    const chart = wrapper.find(DiffChart);
    expect(chart.length).toBe(0);
  });

  it('loads costs when mounting', () => {
    expect(props.getCosts).not.toHaveBeenCalled();
    shallow(<Chart {...props}/>);
    expect(props.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when accounts are updated', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(updatedAccountsProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedAccountsProps);
    expect(updatedAccountsProps.getCosts).toHaveBeenCalled();
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
    expect(updatedFilterPropsPie.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedFilterPropsPie);
    expect(updatedFilterPropsPie.getCosts).toHaveBeenCalled();
  });

  it('does not reload when dates, interval nor filters are updated', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
  });

  it('can show and hide table', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(wrapper.state('table')).toBe(false);
    wrapper.instance().toggleTable({ preventDefault(){} });
    expect(wrapper.state('table')).toBe(true);
    wrapper.instance().toggleTable({ preventDefault(){} });
    expect(wrapper.state('table')).toBe(false);
  });

});

describe('<Header />', () => {

  it('renders <TimerangeSelector/> component when BarChart is selected', () => {
    const wrapper = shallow(<Header {...props}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(1);
    const interval = wrapper.find(IntervalNavigator);
    expect(interval.length).toBe(0);
  });

  it('renders <IntervalNavigator/> component when PieChart is selected', () => {
    const wrapper = shallow(<Header {...propsPie}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(0);
    const interval = wrapper.find(IntervalNavigator);
    expect(interval.length).toBe(1);
  });

  it('renders <TimerangeSelector/> component when DiffChart is selected', () => {
    const wrapper = shallow(<Header {...propsDiff}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(1);
    const interval = wrapper.find(IntervalNavigator);
    expect(interval.length).toBe(0);
  });

  it('renders <i.fa-pie-chart/> component when PieChart is selected', () => {
    const wrapper = shallow(<Header {...propsPie}/>);
    const icon = wrapper.find("i.fa-pie-chart");
    expect(icon.length).toBe(1);
  });

  it('renders <i.fa-table/> component when DiffChart is selected', () => {
    const wrapper = shallow(<Header {...propsDiff}/>);
    const icon = wrapper.find("i.fa-table");
    expect(icon.length).toBe(1);
  });

  it('renders <i.fa-bar-chart/> component when BarChart is selected', () => {
    const wrapper = shallow(<Header {...props}/>);
    const icon = wrapper.find("i.fa-bar-chart");
    expect(icon.length).toBe(1);
  });

  it('renders no <i/> component when icon is not asked', () => {
    const wrapper = shallow(<Header {...propsWithoutIcon}/>);
    const icon = wrapper.find("i.fa");
    expect(icon.length).toBe(0);
  });

  it('renders nothing when setDates is not set', () => {
    const wrapper = shallow(<Header {...propsNoDates}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(0);
    const interval = wrapper.find(IntervalNavigator);
    expect(interval.length).toBe(0);
  });

  it('renders <Selector/> component', () => {
    const wrapper = shallow(<Header {...props}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

  it('renders a <Spinner/> component when data is unavailable', () => {
    const wrapper = shallow(<Header {...propsWithoutData}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders a <button/> component when can be closed', () => {
    const wrapper = shallow(<Header {...propsWithClose}/>);
    const button = wrapper.find("button");
    expect(button.length).toBe(1);
  });

  it('renders a <button/> component when can toggle table', () => {
    let wrapper = shallow(<Header {...propsWithTableShow}/>);
    let button = wrapper.find("button");
    expect(button.length).toBe(1);
    wrapper = shallow(<Header {...propsWithTableHide}/>);
    button = wrapper.find("button");
    expect(button.length).toBe(1);
  });

  it('renders two <button/> components when can toggle table and can be closed', () => {
    const wrapper = shallow(<Header {...propsWithTableAndClose}/>);
    const button = wrapper.find("button");
    expect(button.length).toBe(2);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<Header {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('can set dates', () => {
    const wrapper = shallow(<Header {...props}/>);
    expect(props.setDates).not.toHaveBeenCalled();
    wrapper.instance().setDates(Moment().startOf('month'), Moment().endOf('month'));
    expect(props.setDates).toHaveBeenCalledTimes(1);
  });

  it('can set filter', () => {
    const wrapper = shallow(<Header {...props}/>);
    expect(props.setFilter).not.toHaveBeenCalled();
    wrapper.instance().setFilter("filter");
    expect(props.setFilter).toHaveBeenCalledTimes(1);
  });

  it('can set interval', () => {
    const wrapper = shallow(<Header {...props}/>);
    expect(props.setInterval).not.toHaveBeenCalled();
    wrapper.instance().setInterval("interval");
    expect(props.setInterval).toHaveBeenCalledTimes(1);
  });

  it('can close', () => {
    const wrapper = shallow(<Header {...propsWithClose}/>);
    expect(propsWithClose.close).not.toHaveBeenCalled();
    wrapper.instance().close({ preventDefault() {} });
    expect(propsWithClose.close).toHaveBeenCalledTimes(1);
  });

});
