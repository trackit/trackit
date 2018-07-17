import React from 'react';
import PieChart  from '../PieChartComponent';
import NVD3Chart from 'react-nvd3';
import ReactTable from 'react-table';
import { shallow } from "enzyme";

const props = {
  values: {
    value: 1,
    otherValue: 2
  },
  interval: "day",
  filter: "product",
  legend: true,
  table: true
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

const propsWithoutTable = {
  ...props,
  table: false
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

  it('renders <h4/> component when values are empty', () => {
    const wrapper = shallow(<PieChart {...propsEmptyCosts}/>);
    const error = wrapper.find("h4.no-data");
    expect(error.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<PieChart {...propsWithoutCosts}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('renders <NVD3Chart/> component without margin', () => {
    const wrapper = shallow(<PieChart {...propsWithoutMargin}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders <ReactTable/> component if table is set', () => {
    const wrapper = shallow(<PieChart {...props}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

  it('renders no <ReactTable/> component if table is not set', () => {
    const wrapper = shallow(<PieChart {...propsWithoutTable}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(0);
  });

});
