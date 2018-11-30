import React from 'react';
import { AccountsContainer } from '../AccountsContainer';
import Components from '../../../../components';
import { shallow } from 'enzyme';

const List = Components.AWS.Accounts.List;
const Wizard = Components.AWS.Accounts.Wizard;

const defaultActions = {
  accountActions: {
    new: jest.fn(),
    edit: jest.fn(),
    delete: jest.fn(),
    clearNew: jest.fn()
  }
};

const props = {
  ...defaultActions,
  accounts: {
    status: true,
    values: []
  },
  match: { params : [] },
  external: {
    external: "external",
    accountId: "accountId"
  },
  getAccounts: jest.fn(),
  newExternal: jest.fn(),
  addBill: jest.fn(),
  clearBill: jest.fn()
};

const propsWithParam = {
  ...props,
  match: {
    params: {
      hasAccounts: "false"
    }
  }
};

describe('<AccountsContainer />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <AccountsContainer /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <List /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    const list = wrapper.find(List);
    expect(list.length).toBe(1);
  });

  it('renders a <Wizard /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    const wizard = wrapper.find(Wizard);
    expect(wizard.length).toBe(1);
  });

  it('renders a welcome message', () => {
    const wrapper = shallow(<AccountsContainer {...propsWithParam} />);
    const message = wrapper.find("div#welcome");
    expect(message.length).toBe(1);
  });

});
