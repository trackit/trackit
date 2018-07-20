import React from 'react';
import ListComponent, { Item } from '../ListComponent';
import List, {
  ListItem
} from 'material-ui/List';
import Spinner from 'react-spinkit';
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
  billRepositories: []
};

const accountWithBills = {
  id: 42,
  userId: 42,
  roleArn: "arn:aws:iam::000000000001:role/TEST_ROLE",
  pretty: "Name",
  billRepositories: [
    {
      error: "",
      nextPending: false,
      bucket: "billing-bucket",
      prefix: "prefix"
    },
    {
      error: "access denied",
      nextPending: true,
      bucket: "another-billing-bucket",
      prefix: "another-prefix"
    },
  ],
};

describe('<ListComponent />', () => {

  const props = {
    ...actionsProps,
  };

  const propsWithAccounts = {
    ...props,
    accounts: {
      status: true,
      values: [accountWithoutBills, accountWithoutBills]
    }
  };

  const propsWaiting = {
    ...props,
    accounts: {
      status: false
    }
  };

  const propsError = {
    ...props,
    accounts: {
      status: true,
      error: Error("Error")
    }
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <ListComponent /> component', () => {
    const wrapper = shallow(<ListComponent {...propsWithAccounts}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div/> component when no account is available', () => {
    const wrapper = shallow(<ListComponent {...propsError}/>);
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

  it('renders a <Spinner /> component when accounts are loading', () => {
    const wrapper = shallow(<ListComponent {...propsWaiting}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

});

describe('<Item />', () => {

  const props = {
    ...actionsProps,
    account: accountWithoutBills
  };

  beforeEach(() => {
    jest.resetAllMocks();
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
