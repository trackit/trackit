import React from 'react';
import BarChart  from '../BarChartComponent';
import NVD3Chart from 'react-nvd3';
import { shallow } from "enzyme";

const props = {
  values: {
    value: 1,
    otherValue: 2
  },
  interval: "day",
  filter: "product",
  legend: true,
  title: false
};

const propsWithoutCosts = {
  ...props,
  values: null
};

const propsEmptyCosts = {
  ...props,
  values: {}
};

const propsWithTitle = {
  ...props,
  title: true
};

describe('<BarChart />', () => {

  it('renders a <BarChart /> component', () => {
    const wrapper = shallow(<BarChart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are available', () => {
    const wrapper = shallow(<BarChart {...props}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are empty', () => {
    const wrapper = shallow(<BarChart {...propsEmptyCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<BarChart {...propsWithoutCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('renders no title in <h2 /> component when title is not asked', () => {
    const wrapper = shallow(<BarChart {...props}/>);
    const title = wrapper.find("h2");
    expect(title.length).toBe(0);
  });

  it('renders a title in <h2 /> component when title is asked', () => {
    const wrapper = shallow(<BarChart {...propsWithTitle}/>);
    const title = wrapper.find("h2");
    expect(title.length).toBe(1);
  });

});
