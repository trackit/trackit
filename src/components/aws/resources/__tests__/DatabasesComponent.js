import React from 'react';
import { DatabasesComponent } from '../DatabasesComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";
import ReactTable from 'react-table';
import Moment from 'moment';
import Misc from '../../../misc';
import PropTypes from "prop-types";

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
        instance: {
          id: 'name',
          type: 'type',
          availabilityZone: 'us-west-1',
          engine: 'engine',
          multiAZ: 'yes',
          allocatedStorage: 42,
          costs: {
            instance: 42
          },
          stats: {
            cpu: {
              average: 42,
              peak: 42
            },
            freeSpace: {
              minimum: 42,
              maximum: 42,
              average: 42
            }
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

describe('<DatabasesComponent />', () => {

  it('renders a <DatabasesComponent /> component', () => {
    const wrapper = shallow(<DatabasesComponent {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Tooltip /> component', () => {
    const wrapper = shallow(<DatabasesComponent {...propsWithData}/>);
    const tooltip = wrapper.find(Tooltip);
    expect(tooltip.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<DatabasesComponent {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is loading', () => {
    const wrapper = shallow(<DatabasesComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an <div class="alert" /> component when data is not available', () => {
    const wrapper = shallow(<DatabasesComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

});
