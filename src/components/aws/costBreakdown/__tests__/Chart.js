import React from 'react';
import Chart  from '../Chart';
import NVD3Chart from 'react-nvd3';
import { shallow } from "enzyme";

const props = {
  values: {
    value: 1,
    otherValue: 2
  },
  interval: "day",
  filter: "product"
};

const propsWithoutCosts = {
  ...props,
  values: null
};

describe('<Chart />', () => {
  it('renders a <Chart /> component', () => {
    const wrapper = shallow(<Chart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are available', () => {
    const wrapper = shallow(<Chart {...props}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<Chart {...propsWithoutCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

});
