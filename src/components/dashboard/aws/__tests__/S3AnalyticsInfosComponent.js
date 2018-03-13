import React from 'react';
import S3AnalyticsInfos  from '../S3AnalyticsInfosComponent';
import { shallow } from 'enzyme';
import moment from 'moment';
import AWS from '../../../aws';

const Infos = AWS.S3Analytics.Infos;

const props = {
  id: "id",
  accounts: [],
  getValues: jest.fn(),
  setDates: jest.fn(),
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

const updatedDateProps = {
  ...validProps,
  dates: {
    startDate: moment().startOf('year'),
    endDate: moment(),
  },
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

describe('<S3AnalyticsInfos />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <S3AnalyticsInfos /> component', () => {
    const wrapper = shallow(<S3AnalyticsInfos {...validProps}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Infos/> component', () => {
    const wrapper = shallow(<S3AnalyticsInfos {...validProps}/>);
    const chart = wrapper.find(Infos);
    expect(chart.length).toBe(1);
  });

  it('can get values while mounting if dates are available', () => {
    expect(props.getValues).not.toHaveBeenCalled();
    shallow(<S3AnalyticsInfos {...validProps}/>);
    expect(props.getValues).toHaveBeenCalled();
  });

  it('can not get values while mounting if dates are not available', () => {
    expect(props.getValues).not.toHaveBeenCalled();
    shallow(<S3AnalyticsInfos {...props}/>);
    expect(props.getValues).not.toHaveBeenCalled();
  });

  it('reloads values when dates are updated', () => {
    const wrapper = shallow(<S3AnalyticsInfos {...validProps}/>);
    expect(updatedDateProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedDateProps);
    expect(updatedDateProps.getValues).toHaveBeenCalled();
  });

  it('reloads values when accounts are updated', () => {
    const wrapper = shallow(<S3AnalyticsInfos {...validProps}/>);
    expect(updatedAccountsProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedAccountsProps);
    expect(updatedAccountsProps.getValues).toHaveBeenCalled();
  });

  it('will not reloads values when props are not updated', () => {
    const wrapper = shallow(<S3AnalyticsInfos {...validProps}/>);
    expect(notUpdatedProps.getValues).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(notUpdatedProps);
    expect(notUpdatedProps.getValues).not.toHaveBeenCalled();
  });

});
