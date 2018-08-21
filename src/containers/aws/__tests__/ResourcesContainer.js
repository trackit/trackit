import React from 'react';
import { ResourcesContainer } from '../ResourcesContainer';
import { shallow } from 'enzyme';
import Components from "../../../components";

const AccountSelector = Components.AWS.Accounts.AccountSelector;
const Panel = Components.Misc.Panel;
const VMs = Components.AWS.Resources.VMs;
const Databases = Components.AWS.Resources.Databases;

const props = {
  account: '42',
  selectAccount: jest.fn(),
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

  it('renders an <AccountSelector /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    const selector = wrapper.find(AccountSelector);
    expect(selector.length).toBe(1);
  });

  it('renders an <VMs /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    const block = wrapper.find(VMs);
    expect(block.length).toBe(1);
  });

  it('renders an <Databases /> component', () => {
    const wrapper = shallow(<ResourcesContainer {...props}/>);
    const block = wrapper.find(Databases);
    expect(block.length).toBe(1);
  });

});
