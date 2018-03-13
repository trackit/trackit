import React from 'react';
import S3AnalyticsCharts  from '../S3AnalyticsChartsComponent';
import { shallow } from 'enzyme';
import Spinner from 'react-spinkit';
import moment from 'moment';
import AWS from '../../../aws';
import Misc from '../../../misc';

const Selector = Misc.Selector;
const Charts = AWS.S3Analytics;

const props = {
  id: "id",
  accounts: [],
  getValues: jest.fn(),
  setDates: jest.fn(),
  filter: "filter",
  setFilter: jest.fn(),
};

const validProps = {
  ...props,
  values: {
    status: true,
    data: {
      value: 1,
      otherValue: 2
    }
  },
  dates: {
    startDate: moment().startOf("month"),
    endDate: moment().endOf("month")
  },
};

const propsNoData = {
  ...validProps,
  values: null
};

const propsStorage = {
  ...validProps,
  filter: "storage"
};

const propsBandwidth = {
  ...validProps,
  filter: "bandwidth"
};

const propsRequests = {
  ...validProps,
  filter: "requests"
};

const updatedDateProps = {
  ...validProps,
  dates: {
    startDate: moment().startOf('year'),
    endDate: moment(),
  },
  getValues: jest.fn()
};

const updatedFilterProps = {
  ...validProps,
  filter: "requests",
  getValues: jest.fn()
};

const updatedAccountsProps = {
  ...validProps,
  accounts: ["account"],
  getValues: jest.fn()
};

const notUpdatedProps = {
  ...validProps,
  getValues: jest.fn()
};

describe('<S3AnalyticsCharts />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <S3AnalyticsCharts /> component', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <StorageCostChart/> component when values are available', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...propsStorage}/>);
    const chart = wrapper.find(Charts.StorageCostChart);
    expect(chart.length).toBe(1);
  });

  it('renders <BandwidthCostChart/> component when values are available', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...propsBandwidth}/>);
    const chart = wrapper.find(Charts.BandwidthCostChart);
    expect(chart.length).toBe(1);
  });

  it('renders <RequestsCostChart/> component when values are available', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...propsRequests}/>);
    const chart = wrapper.find(Charts.RequestsCostChart);
    expect(chart.length).toBe(1);
  });

  it('renders <Spinner/> component when values are not available', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...propsNoData}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders <Selector/> component', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...props}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

  it('can get values while mounting if dates are available', () => {
    expect(props.getValues).not.toHaveBeenCalled();
    shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(props.getValues).toHaveBeenCalled();
  });

  it('can not get values while mounting if dates are not available', () => {
    expect(props.getValues).not.toHaveBeenCalled();
    shallow(<S3AnalyticsCharts {...props}/>);
    expect(props.getValues).not.toHaveBeenCalled();
  });

  it('can set filter', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(props.setFilter).not.toHaveBeenCalled();
    wrapper.instance().setFilter("filter");
    expect(props.setFilter).toHaveBeenCalled();
  });

  it('reloads values when dates are updated', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(updatedDateProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedDateProps);
    expect(updatedDateProps.getValues).toHaveBeenCalled();
  });

  it('reloads values when filters are updated', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(updatedFilterProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedFilterProps);
    expect(updatedFilterProps.getValues).toHaveBeenCalled();
  });

  it('reloads values when accounts are updated', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(updatedAccountsProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedAccountsProps);
    expect(updatedAccountsProps.getValues).toHaveBeenCalled();
  });

  it('will not reloads values when props are not updated', () => {
    const wrapper = shallow(<S3AnalyticsCharts {...validProps}/>);
    expect(notUpdatedProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(notUpdatedProps);
    expect(notUpdatedProps.getValues).not.toHaveBeenCalled();
  });

});
