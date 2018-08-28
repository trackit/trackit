import React from 'react';
import { VMsComponent } from '../VMsComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";
import ReactTable from 'react-table';
import Moment from 'moment';
import Misc from '../../../misc';

const Tooltip = Misc.Popover;

const props = {
  getData: jest.fn(),
  clear: jest.fn()
};

const propsWithData = {
  ...props,
  data: {
    status: true,
    value: [
      {
        account: '420',
        reportDate: Moment().toISOString(),
        instances: [
          {
            id: '42',
            state: 'running',
            region: 'us-west-1',
            cpuAverage: 42,
            cpuPeak: 42,
            ioRead: {
              internal: 42
            },
            ioWrite: {
              internal: 42
            },
            networkIn: 42,
            networkOut: 42,
            keyPair: 'key',
            type: 'type',
            tags: {
              Name: 'name'
            }
          }
        ]
      }
    ]
  }
};

const propsLoading = {
  ...props,
  data: {
    status: false,
    value: null
  }
};

const propsWithError = {
  ...props,
  data: {
    status: true,
    error: Error()
  }
};

describe('<VMsComponent />', () => {

  it('renders a <VMsComponent /> component', () => {
    const wrapper = shallow(<VMsComponent {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Tooltip /> component', () => {
    const wrapper = shallow(<VMsComponent {...propsWithData}/>);
    const tooltip = wrapper.find(Tooltip);
    expect(tooltip.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<VMsComponent {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is loading', () => {
    const wrapper = shallow(<VMsComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an <div class="alert" /> component when data is not available', () => {
    const wrapper = shallow(<VMsComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

});
