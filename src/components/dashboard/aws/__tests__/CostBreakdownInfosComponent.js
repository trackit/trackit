import React from 'react';
import CostBreakdownInfos  from '../CostBreakdownInfosComponent';
import { shallow } from 'enzyme';
import moment from 'moment';
import AWS from '../../../aws';

const Infos = AWS.CostBreakdown.Infos;

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

describe('<CostBreakdownInfos />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <CostBreakdownInfos /> component', () => {
    const wrapper = shallow(<CostBreakdownInfos {...validProps}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Infos/> component', () => {
    const wrapper = shallow(<CostBreakdownInfos {...validProps}/>);
    const chart = wrapper.find(Infos);
    expect(chart.length).toBe(1);
  });

  it('can get values', () => {
    const wrapper = shallow(<CostBreakdownInfos {...validProps}/>);
    expect(props.getValues).not.toHaveBeenCalled();
    wrapper.instance().getValues("id", validProps.dates.startDate, validProps.dates.endDate, []);
    expect(props.getValues).toHaveBeenCalled();
  });

});
