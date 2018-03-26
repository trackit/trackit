import React from 'react';
import IntervalSelector from '../IntervalSelector';
import Selector from '../Selector';
import { shallow } from "enzyme";

const props = {
  interval: "interval",
  setInterval: jest.fn()
};

const propsWithAvailableIntervals = {
  ...props,
  availableIntervals: ["month", "week"]
};

describe('<IntervalSelector />', () => {

  it('renders a <IntervalSelector /> component', () => {
    const wrapper = shallow(<IntervalSelector {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Selector /> component inside', () => {
    const wrapper = shallow(<IntervalSelector {...props}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

  it('renders a <Selector /> component inside with available intervals', () => {
    const wrapper = shallow(<IntervalSelector {...propsWithAvailableIntervals}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

});
