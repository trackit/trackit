import React from 'react';
import { ElasticSearchComponent } from '../ElasticSearchComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";
import ReactTable from 'react-table';
import Moment from 'moment';
import Misc from '../../../misc';

const Tooltip = Misc.Popover;

const props = {
  getData: jest.fn(),
  clear: jest.fn(),
  dates: {
    startDate: Moment().startOf("months"),
    endDate: Moment().endOf("months")
  }
};

const propsWithData = {
  ...props,
  data: {
    status: true,
    value: [
      {
        account: '420',
        reportDate: Moment().toISOString(),
        domain: {
          domainId: 'id',
          domainName: 'name',
          region: 'us-west-1',
          costs: {
            instance: 42
          },
          stats: {
            cpu: {
              average: 42,
              peak: 42
            },
            JVMMemoryPressure: {
              average: 42,
              peak: 42
            },
            freeSpace: 42
          },
          totalStorageSpace: 42,
          instanceType: 'type',
          instanceCount: 42,
          tags: {
            key: "value"
          }
        }
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

describe('<ElasticSearchComponent />', () => {

  it('renders a <ElasticSearchComponent /> component', () => {
    const wrapper = shallow(<ElasticSearchComponent {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Tooltip /> component', () => {
    const wrapper = shallow(<ElasticSearchComponent {...propsWithData}/>);
    const tooltip = wrapper.find(Tooltip);
    expect(tooltip.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<ElasticSearchComponent {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is loading', () => {
    const wrapper = shallow(<ElasticSearchComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an <div class="alert" /> component when data is not available', () => {
    const wrapper = shallow(<ElasticSearchComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

});
