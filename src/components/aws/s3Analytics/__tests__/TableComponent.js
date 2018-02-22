import React from 'react';
import TableComponent from '../TableComponent';
import Spinner from 'react-spinkit';
import { shallow } from 'enzyme';
import ReactTable from 'react-table';

const props = {
  data: {}
};

const propsLoading = {
  data: {
    status: false
  }
};

const propsWithData = {
  data: {
    status: true,
    values: {
      bucket: {
        GbMonth: 21,
        BandwidthCost: 42,
        StorageCost: 84
      }
    }
  }
};

const propsWithError = {
  data: {
    status: true,
    error: Error()
  }
};

describe('<TableComponent />', () => {

  it('renders a <TableComponent /> component', () => {
    const wrapper = shallow(<TableComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is not available', () => {
    const wrapper = shallow(<TableComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<TableComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('renders a <ReactTable /> component when data is available', () => {
    const wrapper = shallow(<TableComponent {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

});
