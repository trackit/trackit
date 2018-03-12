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

const propsWithoutMargin = {
  ...props,
  margin: false
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

  it('renders <h4/> component when values are empty', () => {
    const wrapper = shallow(<BarChart {...propsEmptyCosts}/>);
    const error = wrapper.find("h4.no-data");
    expect(error.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<BarChart {...propsWithoutCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('renders <NVD3Chart/> component without margin', () => {
    const wrapper = shallow(<BarChart {...propsWithoutMargin}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

});
