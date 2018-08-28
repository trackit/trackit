import React from 'react';
import { AccountSelectorComponent } from '../AccountSelectorComponent';
import { shallow } from 'enzyme';

const account1 = {
  id: 42,
  roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE",
  pretty: "pretty"
};

const account2 = {
  id: 84,
  roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE_BIS",
  pretty: "pretty_bis"
};

const account3 = {
  id: 21,
  roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE_BIS_AGAIN",
  pretty: "pretty_bis_again"
};

describe('<AccountSelectorComponent />', () => {

  const props = {
    selectAccount: jest.fn(),
    getAccounts: jest.fn(),
    account: '',
    accounts: {
      status: true,
      values: []
    }
  };

  const propsWithAccounts = {
    ...props,
    accounts: {
      status: true,
      values: [account1, account2, account3]
    }
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <AccountSelectorComponent /> component', () => {
    const wrapper = shallow(<AccountSelectorComponent {...propsWithAccounts}/>);
    expect(wrapper.length).toBe(1);
  });
  it('renders a <AccountSelectorComponent /> component without accounts', () => {
    const wrapper = shallow(<AccountSelectorComponent {...props}/>);
    expect(wrapper.get(0)).toBeNull();
  });
});
