import React from 'react';
import { AccountsContainer } from '../AccountsContainer';
import Components from '../../../../components';
import { shallow } from 'enzyme';

const List = Components.AWS.Accounts.List;
const Wizard = Components.AWS.Accounts.Wizard;
const Panel = Components.Misc.Panel;

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
  accounts: [],
  match: { params : [] },
  external: {
    external: "external",
    accountId: "accountId"
  },
  getAccounts: jest.fn(),
  newExternal: jest.fn(),
  newBill: jest.fn()
};

describe('<AccountsContainer />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
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

  it('renders a <Wizard /> component', () => {
    const wrapper = shallow(<AccountsContainer {...props} />);
    const form = wrapper.find(Wizard);
    expect(form.length).toBe(1);
  });

});
