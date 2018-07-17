import React from 'react';
import DifferentiatorChart  from '../DifferentiatorChartComponent';
import Moment from 'moment';
import ReactTable from 'react-table';
import { shallow } from "enzyme";

const props = {
  values: {
    product: [{
      Date: Moment().format(),
      Cost: 42,
      PercentVariation: 4.2
    }, {
      Date: Moment().subtract(1, 'months').format(),
      Cost: 84,
      PercentVariation: 8.4
    }]
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

describe('<DifferentiatorChart />', () => {

  it('renders a <DifferentiatorChart /> component', () => {
    const wrapper = shallow(<DifferentiatorChart {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <ReactTable/> component when values are available', () => {
    const wrapper = shallow(<DifferentiatorChart {...props}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

  it('renders no <ReactTable/> component when values are unavailable', () => {
    const wrapper = shallow(<DifferentiatorChart {...propsWithoutCosts}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(0);
  });

  it('renders <h4/> component when values are empty', () => {
    const wrapper = shallow(<DifferentiatorChart {...propsEmptyCosts}/>);
    const error = wrapper.find("h4.no-data");
    expect(error.length).toBe(1);
  });

});
