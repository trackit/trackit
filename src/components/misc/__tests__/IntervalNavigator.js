import React from 'react';
import IntervalNavigator from '../IntervalNavigator';
import Moment from 'moment';
import IntervalSelector from '../IntervalSelector';
import { shallow } from "enzyme";

const range = {
  startDate: Moment().startOf('week'),
  endDate: Moment(),
};

const props = {
  ...range,
  setDatesFunc: jest.fn(),
  interval: "interval",
  setIntervalFunc: jest.fn()
};

const propsDay = {
  ...props,
  interval: "day"
};

const propsWeek = {
  ...props,
  interval: "week"
};

const propsMonth = {
  ...props,
  interval: "month"
};

const propsYear = {
  ...props,
  interval: "year"
};

describe('<IntervalNavigator />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <IntervalNavigator /> component', () => {
    const wrapper = shallow(<IntervalNavigator {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <IntervalSelector /> component inside', () => {
    const wrapper = shallow(<IntervalNavigator {...props}/>);
    const interval = wrapper.find(IntervalSelector);
    expect(interval.length).toBe(1);
  });

  it('can update interval', () => {
    const wrapper = shallow(<IntervalNavigator {...props}/>);
    expect(props.setIntervalFunc).not.toHaveBeenCalled();
    expect(props.setDatesFunc).not.toHaveBeenCalled();
    wrapper.instance().updateInterval("day");
    expect(props.setIntervalFunc).toHaveBeenCalledTimes(1);
    expect(props.setDatesFunc).toHaveBeenCalledTimes(1);
    wrapper.instance().updateInterval("week");
    expect(props.setIntervalFunc).toHaveBeenCalledTimes(2);
    expect(props.setDatesFunc).toHaveBeenCalledTimes(2);
    wrapper.instance().updateInterval("month");
    expect(props.setIntervalFunc).toHaveBeenCalledTimes(3);
    expect(props.setDatesFunc).toHaveBeenCalledTimes(3);
    wrapper.instance().updateInterval("year");
    expect(props.setIntervalFunc).toHaveBeenCalledTimes(4);
    expect(props.setDatesFunc).toHaveBeenCalledTimes(4);
  });

  it('can go to previous date', () => {
    const wrapperDay = shallow(<IntervalNavigator {...propsDay}/>);
    expect(props.setDatesFunc).not.toHaveBeenCalled();
    wrapperDay.instance().previousDate({ preventDefault() {}});
    expect(propsDay.setDatesFunc).toHaveBeenCalledTimes(1);
    const wrapperWeek = shallow(<IntervalNavigator {...propsWeek}/>);
    wrapperWeek.instance().previousDate({ preventDefault() {}});
    expect(propsWeek.setDatesFunc).toHaveBeenCalledTimes(2);
    const wrapperMonth = shallow(<IntervalNavigator {...propsMonth}/>);
    wrapperMonth.instance().previousDate({ preventDefault() {}});
    expect(propsMonth.setDatesFunc).toHaveBeenCalledTimes(3);
    const wrapperYear = shallow(<IntervalNavigator {...propsYear}/>);
    wrapperYear.instance().previousDate({ preventDefault() {}});
    expect(propsYear.setDatesFunc).toHaveBeenCalledTimes(4);
  });;

  it('can go to next date', () => {
    const wrapperDay = shallow(<IntervalNavigator {...propsDay}/>);
    expect(props.setDatesFunc).not.toHaveBeenCalled();
    wrapperDay.instance().nextDate({ preventDefault() {}});
    expect(propsDay.setDatesFunc).toHaveBeenCalledTimes(1);
    const wrapperWeek = shallow(<IntervalNavigator {...propsWeek}/>);
    wrapperWeek.instance().nextDate({ preventDefault() {}});
    expect(propsWeek.setDatesFunc).toHaveBeenCalledTimes(2);
    const wrapperMonth = shallow(<IntervalNavigator {...propsMonth}/>);
    wrapperMonth.instance().nextDate({ preventDefault() {}});
    expect(propsMonth.setDatesFunc).toHaveBeenCalledTimes(3);
    const wrapperYear = shallow(<IntervalNavigator {...propsYear}/>);
    wrapperYear.instance().nextDate({ preventDefault() {}});
    expect(propsYear.setDatesFunc).toHaveBeenCalledTimes(4);
  });

});
