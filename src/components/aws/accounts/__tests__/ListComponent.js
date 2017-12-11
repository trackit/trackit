import React from 'react';
import ListComponent, { Item } from '../ListComponent';
import List, {
  ListItem
} from 'material-ui/List';
import { shallow } from 'enzyme';

const actionsProps = {
  accountActions: {
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

  it('renders a <List/> component when accounts are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithAccounts}/>);
    const listWrapper = wrapper.find(List);
    expect(listWrapper.length).toBe(1);
  });

  it('renders 2 <Item /> component when 2 accounts are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithAccounts}/>);
    const list = wrapper.find(Item);
    expect(list.length).toBe(2);
  });

});

describe('<Item />', () => {

  const props = {
    ...actionsProps,
    account: accountWithoutBills
  };

  beforeEach(() => {
    actionsProps.accountActions.new.mockReset();
    actionsProps.accountActions.edit.mockReset();
    actionsProps.accountActions.delete.mockReset();
  });

  it('renders a <Item /> component', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <ListItem/> component', () => {
    const wrapper = shallow(<Item {...props}/>);
    const item = wrapper.find(ListItem);
    expect(item.length).toBe(1);
  });

  it('can edit item', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(props.accountActions.edit).not.toHaveBeenCalled();
    wrapper.instance().editAccount(accountWithBills);
//    expect(props.accountActions.edit).toHaveBeenCalled();
  });

  it('can delete item', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(props.accountActions.delete).not.toHaveBeenCalled();
    wrapper.instance().deleteAccount();
//    expect(props.accountActions.delete).toHaveBeenCalled();
  });

});
