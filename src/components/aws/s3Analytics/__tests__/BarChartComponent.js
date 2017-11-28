import React from 'react';
import BarChartComponent from '../BarChartComponent';
import { shallow } from 'enzyme';

const data = {
  _id: "id",
  size: 42,
  storage_cost: 42,
  bw_cost: 42,
  total_cost: 42,
  transfer_in: 42,
  transfer_out: 42
};

const otherData = {
  _id: "id2",
  size: 84,
  storage_cost: 84,
  bw_cost: 84,
  total_cost: 84,
  transfer_in: 84,
  transfer_out: 84
};

const defaultProps = {
  elementId: "barchart",
  data: []
};

const props = {
  ...defaultProps,
  data: [data]
};

const propsWithMultipleData = {
  ...defaultProps,
  data: [data, otherData]
};

const propsWithMultipleDataInversed = {
  ...defaultProps,
  data: [otherData, data]
};

const bandwidth = {
  x: [],
  y: [],
  name: 'Bandwidth',
  type: 'bar',
  opacity: 0.8,
  marker: {
    color: '#1e88e5',
  }
};
const storage = {
  x: [],
  y: [],
  name: 'Storage',
  type: 'bar',
  opacity: 0.8,
  marker: {
    color: '#ff9800',
  },
  hoverlabel: {
    bordercolor: '#ffffff',
  }
}

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

  it('formats data for chart', () => {
    const wrapper = shallow(<BarChartComponent {...props}/>);
    const data = wrapper.instance().formatDataForChart();
    const bandwidthFormatted = {
      ...bandwidth,
      x: [props.data[0]._id],
      y: [props.data[0].bw_cost.toFixed(2)],
    };
    const storageFormatted = {
      ...storage,
      x: [props.data[0]._id],
      y: [props.data[0].storage_cost.toFixed(2)],
    };
    expect(data[0]).toEqual(bandwidthFormatted);
    expect(data[1]).toEqual(storageFormatted);
  });

  it('formats multiple data for chart', () => {
    const wrapper = shallow(<BarChartComponent {...propsWithMultipleData}/>);
    const data = wrapper.instance().formatDataForChart();
    const bandwidthFormatted = {
      ...bandwidth,
      x: [
        propsWithMultipleData.data[1]._id,
        propsWithMultipleData.data[0]._id
      ],
      y: [
        propsWithMultipleData.data[1].bw_cost.toFixed(2),
        propsWithMultipleData.data[0].bw_cost.toFixed(2)
      ],
    };
    const storageFormatted = {
      ...storage,
      x: [
        propsWithMultipleData.data[1]._id,
        propsWithMultipleData.data[0]._id
      ],
      y: [
        propsWithMultipleData.data[1].storage_cost.toFixed(2),
        propsWithMultipleData.data[0].storage_cost.toFixed(2)
      ],
    };
    expect(data[0]).toEqual(bandwidthFormatted);
    expect(data[1]).toEqual(storageFormatted);
  });

  it('orders data for chart', () => {
    const wrapper = shallow(<BarChartComponent {...propsWithMultipleDataInversed}/>);
    const data = wrapper.instance().formatDataForChart();
    const bandwidthFormatted = {
      ...bandwidth,
      x: [
        propsWithMultipleDataInversed.data[0]._id,
        propsWithMultipleDataInversed.data[1]._id
      ],
      y: [
        propsWithMultipleDataInversed.data[0].bw_cost.toFixed(2),
        propsWithMultipleDataInversed.data[1].bw_cost.toFixed(2)
      ],
    };
    const storageFormatted = {
      ...storage,
      x: [
        propsWithMultipleDataInversed.data[0]._id,
        propsWithMultipleDataInversed.data[1]._id
      ],
      y: [
        propsWithMultipleDataInversed.data[0].storage_cost.toFixed(2),
        propsWithMultipleDataInversed.data[1].storage_cost.toFixed(2)
      ],
    };
    expect(data[0]).toEqual(bandwidthFormatted);
    expect(data[1]).toEqual(storageFormatted);
  });

});
