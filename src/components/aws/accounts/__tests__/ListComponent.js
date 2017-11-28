import React from 'react';
import ListComponent, { ListItem } from '../ListComponent';
import { shallow } from 'enzyme';

const actionsProps = {
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

const accountWithoutBills = {
  id: 42,
  userId: 42,
  roleArn: "arn:aws:iam::000000000001:role/TEST_ROLE",
  bills: []
};

const accountWithBills = {
  id: 42,
  userId: 42,
  roleArn: "arn:aws:iam::000000000001:role/TEST_ROLE",
  pretty: "Name",
  bills: []
};

describe('<ListComponent />', () => {

  const props = {
    ...actionsProps,
    accounts: []
  };

  const propsWithAccounts = {
    ...props,
    accounts: [accountWithoutBills, accountWithoutBills]
  };

  beforeEach(() => {
    actionsProps.accountActions.new.mockReset();
    actionsProps.accountActions.edit.mockReset();
    actionsProps.accountActions.delete.mockReset();
    actionsProps.billActions.new.mockReset();
    actionsProps.billActions.edit.mockReset();
    actionsProps.billActions.delete.mockReset();
  });

  it('renders a <ListComponent /> component', () => {
    const wrapper = shallow(<ListComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div/> component when no account is available', () => {
    const wrapper = shallow(<ListComponent {...props}/>);
    const alert = wrapper.find('div');
    expect(alert.length).toBe(1);
  });

  it('renders a <ul/> component when accounts are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithAccounts}/>);
    const listWrapper = wrapper.find('ul');
    expect(listWrapper.length).toBe(1);
  });

  it('renders 2 <ListItem /> component when 2 accounts are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithAccounts}/>);
    const list = wrapper.find(ListItem);
    expect(list.length).toBe(2);
  });

});

describe('<ListItem />', () => {

  const props = {
    ...actionsProps,
    account: accountWithoutBills
  };

  beforeEach(() => {
    actionsProps.accountActions.new.mockReset();
    actionsProps.accountActions.edit.mockReset();
    actionsProps.accountActions.delete.mockReset();
    actionsProps.billActions.new.mockReset();
    actionsProps.billActions.edit.mockReset();
    actionsProps.billActions.delete.mockReset();
  });

  it('renders a <ListItem /> component', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <li/> component', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    const item = wrapper.find('li');
    expect(item.length).toBe(1);
  });

  it('renders 2 <button/> components', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    const buttons = wrapper.find('button');
    expect(buttons.length).toBe(2);
  });

  it('can expand edit form', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    expect(wrapper.state('editForm')).toBe(false);
    wrapper.find('button.btn.edit').prop('onClick')({ preventDefault() {} });
    expect(wrapper.state('editForm')).toBe(true);
  });

  it('can edit item', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    expect(props.accountActions.edit.mock.calls.length).toBe(0);
    wrapper.instance().editAccount(accountWithBills);
    expect(props.accountActions.edit.mock.calls.length).toBe(1);
  });

  it('can delete item', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    expect(props.accountActions.delete.mock.calls.length).toBe(0);
    wrapper.find('button.btn.delete').prop('onClick')({ preventDefault() {} });
    expect(props.accountActions.delete.mock.calls.length).toBe(1);
  });

});
