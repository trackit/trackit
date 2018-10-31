import React from 'react';
import Components from '../../../components';
import moment from 'moment';
import { shallow } from "enzyme";

const UnusedComponent = Components.HighLevel.TopUnused;

const instance = {
  account: '420',
  reportDate: moment().toISOString(),
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
};

const props = {
    date: moment().startOf('month'),
    unused: {
        ec2: {
            status: true,
            values: [
              instance,
              instance
            ],
        },      
    }
};

const propsNoData = {
    date: moment().startOf('month'),
    unused: {
        ec2: {
            status: true,
            values: [],
        },      
    }
};



describe('<UnusedComponent />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <UnusedComponent /> component', () => {
    const wrapper = shallow(<UnusedComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders panel title', () => {
    const wrapper = shallow(<UnusedComponent {...props}/>);
    const title = wrapper.find('.hl-panel-title');
    expect(title.length).toBe(1);
  });

  it('renders a message when no data are returned', () => {
    const wrapper = shallow(<UnusedComponent {...propsNoData}/>);
    const title = wrapper.find('.no-resource-message');
    expect(title.length).toBe(1);
  });

  it('renders a table displaying the EC2 data', () => {
    const wrapper = shallow(<UnusedComponent {...props}/>);
    const table = wrapper.find('tbody tr');
    expect(table.length).toBe(2);
  });

});
