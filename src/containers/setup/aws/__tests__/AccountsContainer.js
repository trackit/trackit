import React from 'react';
import { AccountsContainer } from '../AccountsContainer';
import Components from '../../../../components';
import { shallow } from 'enzyme';

const List = Components.AWS.Accounts.List;
const Form = Components.AWS.Accounts.Form;
const Panel = Components.Misc.Panel;

const accountActions = {
  new: jest.fn(),
  edit: jest.fn(),
  delete: jest.fn(),
};

const billActions = {
  new: jest.fn(),
  edit: jest.fn(),
  delete: jest.fn(),
};

const props = {
  accounts: [],
  external: "external",
  getAccounts: jest.fn(),
  accountActions,
  billActions,
  newExternal: jest.fn()
};

describe('<AccountsContainer />', () => {

  it('renders a <AccountsContainer /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Panel /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    const panel = wrapper.find(Panel);
    expect(panel.length).toBe(1);
  });

  it('renders a <List /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    const list = wrapper.find(List);
    expect(list.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

});
