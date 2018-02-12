import React from 'react';
import { SelectedIndicator } from '../SelectedIndicatorComponent';
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

describe('<SelectedIndicator />', () => {

  const props = {
    select: jest.fn(),
    clear: jest.fn(),
    getAccounts: jest.fn(),
    accounts: {
      status: true,
      values: []
    },
    selection: []
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

  const propsWithAccounts = {
    ...props,
    accounts: {
      status: true,
      values: [account1, account2, account3]
    }
  };

  const propsWithSelectedAccount = {
    ...propsWithAccounts,
    selection: [account1]
  };

  const propsWithSelectedAccountLong = {
    ...propsWithAccounts,
    selection: [account1],
    longVersion: true
  };

  const propsWithSelectedAccounts = {
    ...propsWithAccounts,
    selection: [account1, account3]
  };

  const propsWithSelectedAccountsLong = {
    ...propsWithAccounts,
    selection: [account1, account3],
    longVersion: true
  };

  const propsWithSelectedAllAccounts = {
    ...propsWithAccounts,
    selection: [account1, account2, account3]
  };

  const propsWithSelectedAllAccountsLong = {
    ...propsWithAccounts,
    selection: [account1, account2, account3],
    longVersion: true
  };

  const propsWithIcon = {
    ...props,
    icon: true
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <SelectedIndicator /> component', () => {
    const wrapper = shallow(<SelectedIndicator {...propsWithAccounts}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a badge component', () => {
    const wrapper = shallow(<SelectedIndicator {...propsWithAccounts}/>);
    const badge = wrapper.find("span.badge");
    expect(badge.length).toBe(1);
  });

  it('renders an icon component', () => {
    const wrapper = shallow(<SelectedIndicator {...propsWithIcon}/>);
    const icon = wrapper.find("i");
    expect(icon.length).toBe(1);
  });

  describe('Get Text', () => {

    const noAccount = 'No AWS account available';
    const allAccounts = 'All accounts';
    const multipleAccounts = (count) => (`${count} accounts`);
    const longVersion = (text) => (`Displaying ${text}`);

    it('returns nothing when data is loading', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWaiting}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(null);
    });
    it('returns "No account" message when no account available', () => {
      const wrapper = shallow(<SelectedIndicator {...props}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(noAccount);
    });

    it('returns "No account" message with error message when there is an error', () => {
      const wrapper = shallow(<SelectedIndicator {...propsError}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(`${noAccount} (${propsError.accounts.error.message})`);
    });

    it('returns "All accounts" message when selection is empty', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithAccounts}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(allAccounts);
    });

    it('returns "All accounts" message when all accounts are selected', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithSelectedAllAccounts}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(allAccounts);
    });

    it('returns "All accounts (Long version)" message when all accounts are selected and long version is enabled', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithSelectedAllAccountsLong}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(longVersion(allAccounts));
    });

    it('returns "Single account" message when only one account is selected', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithSelectedAccount}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(propsWithSelectedAccount.selection[0].pretty);
    });

    it('returns "Single account (Long version)" message when only one account is selected and long version is enabled', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithSelectedAccountLong}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(longVersion(propsWithSelectedAccountLong.selection[0].pretty));
    });

    it('returns "Multiple accounts" message when accounts are selected', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithSelectedAccounts}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(multipleAccounts(propsWithSelectedAccounts.selection.length));
    });

    it('returns "Multiple accounts (Long version)" message when accounts are selected and long version is enabled', () => {
      const wrapper = shallow(<SelectedIndicator {...propsWithSelectedAccountsLong}/>);
      const instance = wrapper.instance();
      expect(instance.getText()).toBe(longVersion(multipleAccounts(propsWithSelectedAccountsLong.selection.length)));
    });

  });

});
