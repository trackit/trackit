import React from 'react';
import TableComponent from '../TableComponent';
import { shallow } from 'enzyme';
import ReactTable from 'react-table';

const props = {
  data: [{
    _id: "id",
    size: 42,
    storage_cost: 42,
    bw_cost: 42,
    total_cost: 42,
    transfer_in: 42,
    transfer_out: 42,
    chargify: 'not_synced'
  }]
};

describe('<TableComponent />', () => {

  it('renders a <TableComponent /> component', () => {
    const wrapper = shallow(<TableComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<TableComponent {...props}/>);
    const form = wrapper.find(ReactTable);
    expect(form.length).toBe(1);
  });

});
