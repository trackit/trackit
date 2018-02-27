import React from 'react';
import InfosComponent from '../InfosComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";

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
        StorageCost: 84,
        RequestsCost: 126
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

describe('<InfosComponent />', () => {

  it('renders a <InfosComponent /> component', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is not available', () => {
    const wrapper = shallow(<InfosComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<InfosComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('calculates totals based on data', () => {
    const wrapper = shallow(<InfosComponent {...propsWithData}/>);
    const totals = wrapper.instance().extractTotals();
    expect(totals.buckets).toBe(Object.keys(propsWithData.data.values).length);
    expect(totals.size).toBe(propsWithData.data.values.bucket.GbMonth);
    expect(totals.bandwidth_cost).toBe(propsWithData.data.values.bucket.BandwidthCost);
    expect(totals.storage_cost).toBe(propsWithData.data.values.bucket.StorageCost);
    expect(totals.requests_cost).toBe(propsWithData.data.values.bucket.RequestsCost);
  });

});
