import React from 'react';
import TimerangeSelector from '../TimerangeSelector';
import Moment from 'moment';
import DateRangePicker from 'react-bootstrap-daterangepicker';
import { shallow } from "enzyme";

const range = {
  startDate: Moment().startOf('week'),
  endDate: Moment(),
};

const props = {
  ...range,
  setDatesFunc: jest.fn()
};

describe('<TimerangeSelector />', () => {

  it('renders a <TimerangeSelector /> component', () => {
    const wrapper = shallow(<TimerangeSelector {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <DateRangePicker /> component inside', () => {
    const wrapper = shallow(<TimerangeSelector {...props}/>);
    const picker = wrapper.find(DateRangePicker);
    expect(picker.length).toBe(1);
  });

  it('can select range', () => {
    const wrapper = shallow(<TimerangeSelector {...props}/>);
    expect(props.setDatesFunc.mock.calls.length).toBe(0);
    wrapper.instance().handleApply({ preventDefault(){} }, range);
    expect(props.setDatesFunc.mock.calls.length).toBe(1);
  });

});
