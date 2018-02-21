import React from 'react';
import PieChart  from '../PieChartComponent';
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

const propsEmptyCosts = {
  ...props,
  values: {}
};

describe('<PieChart />', () => {

  it('renders a <PieChart /> component', () => {
    const wrapper = shallow(<PieChart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are available', () => {
    const wrapper = shallow(<PieChart {...props}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are empty', () => {
    const wrapper = shallow(<PieChart {...propsEmptyCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<PieChart {...propsWithoutCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

});
