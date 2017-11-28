import React from 'react';
import InfosComponent from '../InfosComponent';
import { shallow } from 'enzyme';

const props = {
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

describe('<InfosComponent />', () => {

  it('renders a <InfosComponent /> component', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('calculates totals based on data', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    const totals = wrapper.instance().extractTotals();
    expect(totals.buckets).toBe(props.data.length);
    expect(totals.size).toBe(props.data[0].size);
    expect(totals.storage_cost).toBe(props.data[0].storage_cost);
  });

});
