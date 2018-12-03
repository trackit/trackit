import React from 'react';
import { Ec2Component } from '../Ec2Component';
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
          id: '42',
          state: 'running',
          region: 'us-west-1',
          keyPair: 'key',
          type: 'type',
          purchasing: 'value',
          tags: {
            Name: 'name'
          },
          costs: {
            instance: 42
          },
          stats: {
            cpu: {
              average: 42,
              peak: 42
            },
            network: {
              in: 42,
              out: 42
            },
            volumes: {
              read: {
                internal: 42
              },
              write: {
                internal: 42
              }
            }
          }
        }
      }
    ]
  }
};

describe('<Ec2Component />', () => {

  it('renders a <Ec2Component /> component', () => {
    const wrapper = shallow(<Ec2Component {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Collapsible /> component', () => {
    const wrapper = shallow(<Ec2Component {...propsWithData}/>);
    const collapsible = wrapper.find(Collapsible);
    expect(collapsible.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<Ec2Component {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

});
