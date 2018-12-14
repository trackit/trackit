import React from 'react';
import { ElastiCacheComponent } from '../ElastiCacheComponent';
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
          id: 'id',
          region: 'us-west-1',
          nodeType: 'm4.large',
          engine: 'redis',
          engineVersion: '3.1.4',
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
            }
          },
          tags: {
            key: "value"
          }
        }
      }
    ]
  }
};

describe('<ElastiCacheComponent />', () => {

  it('renders a <ElastiCacheComponent /> component', () => {
    const wrapper = shallow(<ElastiCacheComponent {...propsWithData}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Collapsible /> component', () => {
    const wrapper = shallow(<ElastiCacheComponent {...propsWithData}/>);
    const collapsible = wrapper.find(Collapsible);
    expect(collapsible.length).toBe(1);
  });

  it('renders a <ReactTable /> component', () => {
    const wrapper = shallow(<ElastiCacheComponent {...propsWithData}/>);
    const table = wrapper.find(ReactTable);
    expect(table.length).toBe(1);
  });

});
