import React from 'react';
import ConnectecAccountsContainer, { AccountsContainer } from '../AccountsContainer';
import Components from '../../../../components';
import { shallow } from 'enzyme';
import { createMockStore } from 'redux-test-utils';

const List = Components.AWS.Accounts.List;
const Form = Components.AWS.Accounts.Form;
const Panel = Components.Misc.Panel;

const defaultActions = {
  accountActions: {
    new: jest.fn(),
    edit: jest.fn(),
    delete: jest.fn(),
  },
  billActions: {
    new: jest.fn(),
    edit: jest.fn(),
    delete: jest.fn(),
  }
};

const props = {
  ...defaultActions,
  accounts: [],
  external: "external",
  getAccounts: jest.fn(),
  newExternal: jest.fn()
};

const state = {
  aws: {
    accounts: {
      all: [{
        id: 42,
        roleArn: "role",
        pretty: "pretty",
        bills: []
      }],
      external: "external"
    }
  }
};

describe('<AccountsContainer />', () => {

  beforeEach(() => {
    defaultActions.accountActions.new.mockReset();
    defaultActions.accountActions.edit.mockReset();
    defaultActions.accountActions.delete.mockReset();
    defaultActions.billActions.new.mockReset();
    defaultActions.billActions.edit.mockReset();
    defaultActions.billActions.delete.mockReset();
  });

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
