import React from 'react';
import { ResourcesContainer } from '../ResourcesContainer';
import { shallow } from 'enzyme';
import Components from "../../../components";
import Moment from "moment";

const Panel = Components.Misc.Panel;
const Ec2 = Components.AWS.Resources.Ec2;
const Rds = Components.AWS.Resources.Rds;

const props = {
  account: '42',
  selectAccount: jest.fn(),
  dates: {
    startDate: Moment().startOf("months"),
    endDate: Moment().endOf("months")
  }
};

describe('<ResourcesContainer />', () => {

  it('renders a <ResourcesContainer /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders an <Panel /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    const panel = wrapper.find(Panel);
    expect(panel.length).toBe(1);
  });

  it('renders an <VMs /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    const block = wrapper.find(Ec2);
    expect(block.length).toBe(1);
  });

  it('renders an <Databases /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    const block = wrapper.find(Rds);
    expect(block.length).toBe(1);
  });

});
