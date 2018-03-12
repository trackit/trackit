import React from 'react';
import Charts, { PieChartComponent } from '../PieChartComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";
import NVD3Chart from "react-nvd3";

const StorageCostChartComponent = Charts.StorageCostChartComponent;
const BandwidthCostChartComponent = Charts.BandwidthCostChartComponent;
const RequestsCostChartComponent = Charts.RequestsCostChartComponent;

const props = {
  data: {},
  mode: "storage"
};

const propsLoading = {
  ...props,
  data: {
    status: false
  }
};

const propsWithData = {
  ...props,
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
  ...props,
  data: {
    status: true,
    values: {}
  }
};

const propsWithError = {
  ...props,
  data: {
    status: true,
    error: Error()
  }
};

const propsWithDataBandwidth = {
  ...propsWithData,
  mode: "bandwidth"
};

const propsWithDataRequests = {
  ...propsWithData,
  mode: "requests"
};

const propsWithDataWrongMode = {
  ...propsWithData,
  mode: "mode"
};

describe('<PieChartComponent />', () => {

  it('renders a <PieChartComponent /> component', () => {
    const wrapper = shallow(<PieChartComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is not available', () => {
    const wrapper = shallow(<PieChartComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('renders <NVD3Chart/> component when values are available', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithData}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(1);
  });

  it('renders no <NVD3Chart/> component when values are empty', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithEmptyData}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('renders no <NVD3Chart/> component when values are unavailable', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithError}/>);
    const chart = wrapper.find(NVD3Chart);
    expect(chart.length).toBe(0);
  });

  it('can generate data for "storage" mode', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithData}/>);
    expect(wrapper.instance().generateDatum()).not.toBe(null);
  });

  it('can generate data for "bandwidth" mode', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithDataBandwidth}/>);
    expect(wrapper.instance().generateDatum()).not.toBe(null);
  });

  it('can generate data for "requests" mode', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithDataRequests}/>);
    expect(wrapper.instance().generateDatum()).not.toBe(null);
  });

  it('can not generate data for an invalid mode', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithDataWrongMode}/>);
    expect(wrapper.instance().generateDatum()).toBe(null);
  });

  it('can not generate data if data is not available', () => {
    const wrapper = shallow(<PieChartComponent {...propsWithEmptyData}/>);
    expect(wrapper.instance().generateDatum()).toBe(null);
  });

});

describe('<StorageCostChartComponent />', () => {

  it('renders a <StorageCostChartComponent /> component', () => {
    const wrapper = shallow(<StorageCostChartComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <PieChartComponent /> component', () => {
    const wrapper = shallow(<StorageCostChartComponent {...props}/>);
    const chart = wrapper.find(PieChartComponent);
    expect(chart.length).toBe(1);
  });

});

describe('<BandwidthCostChartComponent />', () => {

  it('renders a <BandwidthCostChartComponent /> component', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <PieChartComponent /> component', () => {
    const wrapper = shallow(<BandwidthCostChartComponent {...props}/>);
    const chart = wrapper.find(PieChartComponent);
    expect(chart.length).toBe(1);
  });

});

describe('<RequestsCostChartComponent />', () => {

  it('renders a <RequestsCostChartComponent /> component', () => {
    const wrapper = shallow(<RequestsCostChartComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <PieChartComponent /> component', () => {
    const wrapper = shallow(<RequestsCostChartComponent {...props}/>);
    const chart = wrapper.find(PieChartComponent);
    expect(chart.length).toBe(1);
  });

});
