import React from 'react';
import Components from '../../../components';
import moment from 'moment';
import { shallow } from "enzyme";

const UnusedComponent = Components.HighLevel.TopUnused;

const props = {
    date: moment().startOf('month'),
    unused: {
        ec2: {
            status: true,
            values: [
                {
                    id: "1234",
                    cost: 123,
                    cpuAverage: 2,
                    tags: {},
                },
                {
                    id: "2234",
                    cost: 223,
                    cpuAverage: 3,
                    tags: {},
                }
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



describe('<S3AnalyticsContainer />', () => {

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
