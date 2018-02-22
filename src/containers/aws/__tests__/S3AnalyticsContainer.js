import React from 'react';
import { S3AnalyticsContainer } from '../S3AnalyticsContainer';
import Components from '../../../components';
import Moment from 'moment';
import { shallow } from "enzyme";

const TimerangeSelector = Components.Misc.TimerangeSelector;
const S3Analytics = Components.AWS.S3Analytics;

const props = {
  values: {},
  accounts: [{
    name: "account1"
  }],
  dates: {
    startDate: Moment(),
    endDate: Moment(),
  },
  getData: jest.fn(),
  setDates: jest.fn()
};

const propsNoDates = {
  ...props,
  dates: null
};

const propsUpdatedAccounts = {
  ...props,
  accounts: [{
    name: "account1"
  },{
    name: "account2"
  }]
};

describe('<S3AnalyticsContainer />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <S3AnalyticsContainer /> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <TimerangeSelector/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(TimerangeSelector);
    expect(navigation.length).toBe(1);
  });

  it('renders no <TimerangeSelector/> component if no dates are available', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...propsNoDates}/>);
    const navigation = wrapper.find(TimerangeSelector);
    expect(navigation.length).toBe(0);
  });

  it('renders <S3Analytics.Infos/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.Infos);
    expect(navigation.length).toBe(1);
  });

  it('renders <S3Analytics.BandwidthCostChart/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.BandwidthCostChart);
    expect(navigation.length).toBe(1);
  });

  it('renders <S3Analytics.StorageCostChart/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.StorageCostChart);
    expect(navigation.length).toBe(1);
  });

  it('renders <S3Analytics.Table/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.Table);
    expect(navigation.length).toBe(1);
  });

  it('loads data when mounting', () => {
    expect(props.getData).not.toHaveBeenCalled();
    shallow(<S3AnalyticsContainer {...props}/>);
    expect(props.getData).toHaveBeenCalled();
  });

  it('set dates if not available when mounting', () => {
    expect(props.setDates).not.toHaveBeenCalled();
    shallow(<S3AnalyticsContainer {...propsNoDates}/>);
    expect(props.setDates).toHaveBeenCalled();
  });

  it('does not reload data when dates are not available', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    expect(props.getData).toHaveBeenCalledTimes(1);
    wrapper.instance().componentWillReceiveProps(propsNoDates);
    expect(props.getData).toHaveBeenCalledTimes(1);
  });

  it('reloads data when dates are updated', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...propsNoDates}/>);
    expect(props.getData).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.getData).toHaveBeenCalled();
  });

  it('reloads data when selected accounts are updated', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    expect(propsUpdatedAccounts.getData).toHaveBeenCalledTimes(1);
    wrapper.instance().componentWillReceiveProps(propsUpdatedAccounts);
    expect(propsUpdatedAccounts.getData).toHaveBeenCalledTimes(2);
  });

});
