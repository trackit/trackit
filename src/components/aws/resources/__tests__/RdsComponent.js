import React from 'react';
import { RdsComponent } from '../RdsComponent';
import { shallow } from 'enzyme';
import ReactTable from 'react-table';
import Moment from 'moment';
import Misc from '../../../misc';

const Collapsible = Misc.Collapsible;

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

describe('<RdsComponent />', () => {

  it('renders a <RdsComponent /> component', () => {
    const wrapper = shallow(<RdsComponent {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Collapsible /> component', () => {
    const wrapper = shallow(<RdsComponent {...propsWithData}/>);
    const collapsible = wrapper.find(Collapsible);
    expect(collapsible.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<RdsComponent {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

});
