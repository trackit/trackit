import React from 'react';
import BarChartComponent from '../BarChartComponent';
import { shallow } from 'enzyme';

const props = {
  elementId: "barchart",
  data: [{
    _id: "id",
    size: 42,
    storage_cost: 42,
    bw_cost: 42,
    total_cost: 42,
    transfer_in: 42,
    transfer_out: 42
  }]
};

describe('<BarChartComponent />', () => {

  it('renders a <BarChartComponent /> component', () => {
    const wrapper = shallow(<BarChartComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div /> component with elementID as id', () => {
    const wrapper = shallow(<BarChartComponent {...props}/>);
    const div = wrapper.find(`div#${props.elementId}`);
    expect(div.length).toBe(1);
  });

});
