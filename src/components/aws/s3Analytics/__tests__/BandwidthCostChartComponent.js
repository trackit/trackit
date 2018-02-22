import React from 'react';
import BandwidthCostChartComponent from '../BandwidthCostChartComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";
import NVD3Chart from "react-nvd3";

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

const propsWithEmptyData = {
  data: {
    status: true,
    values: {}
  }
};

const propsWithError = {
  data: {
    status: true,
    error: Error()
  }
};

describe('<BandwidthCostChartComponent />', () => {

  it('renders a <BandwidthCostChartComponent /> component', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is not available', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are available', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...propsWithData}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are empty', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...propsWithEmptyData}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...propsWithError}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

});
