import React from 'react';
import { SelectorComponent, Item } from '../SelectorComponent';
import { shallow } from 'enzyme';
import List, {
  ListItem,
  ListItemText,
} from 'material-ui/List';
import Spinner from 'react-spinkit';
import Checkbox from 'material-ui/Checkbox';

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

describe('<SelectorComponent />', () => {

  const props = {
    select: jest.fn(),
    clear: jest.fn(),
    getAccounts: jest.fn(),
    accounts: [],
    selected: []
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
      values: [account1, account2]
    }
  };

  const propsWithSelectedAccounts = {
    ...propsWithAccounts,
    selected: [account1]
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <SelectorComponent /> component', () => {
    const wrapper = shallow(<SelectorComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <List /> component', () => {
    const wrapper = shallow(<SelectorComponent {...props}/>);
    const form = wrapper.find(List);
    expect(form.length).toBe(1);
  });

  it('renders multiple <Item /> components', () => {
    const wrapper = shallow(<SelectorComponent {...propsWithAccounts}/>);
    const form = wrapper.find(Item);
    expect(form.length).toBe(propsWithAccounts.accounts.values.length);
  });

  it('renders multiple <Item /> components with some selected', () => {
    const wrapper = shallow(<SelectorComponent {...propsWithSelectedAccounts}/>);
    const form = wrapper.find(Item);
    expect(form.length).toBe(propsWithAccounts.accounts.values.length);
  });

  it('renders an alert components when accounts are missing', () => {
    const wrapper = shallow(<SelectorComponent {...propsError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('renders a <Spinner /> component when accounts are loading', () => {
    const wrapper = shallow(<SelectorComponent {...propsWaiting}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

});

describe('<Item />', () => {

  const props = {
    account: account1,
    select: jest.fn(),
    isSelected: false
  };

  const propWithARN = {
    ...props,
    account: {
      id: 42,
      roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE",
    }
  };

  const propsWithSelected = {
    ...props,
    isSelected: true
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <Item /> component', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <ListItem /> component', () => {
    const wrapper = shallow(<Item {...props}/>);
    const form = wrapper.find(ListItem);
    expect(form.length).toBe(1);
  });

  it('renders a <ListItemText /> component with pretty name', () => {
    const wrapper = shallow(<Item {...props}/>);
    const form = wrapper.find(ListItemText);
    expect(form.length).toBe(1);
  });

  it('renders a <ListItemText /> component with roleARN', () => {
    const wrapper = shallow(<Item {...propWithARN}/>);
    const form = wrapper.find(ListItemText);
    expect(form.length).toBe(1);
  });

  it('renders a <Checkbox /> component when not selected', () => {
    const wrapper = shallow(<Item {...props}/>);
    const form = wrapper.find(Checkbox);
    expect(form.length).toBe(1);
  });

  it('renders a <Checkbox /> component when selected', () => {
    const wrapper = shallow(<Item {...propsWithSelected}/>);
    const form = wrapper.find(Checkbox);
    expect(form.length).toBe(1);
  });

  it('can select an account', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(props.select).not.toHaveBeenCalled();
    wrapper.instance().selectAccount({ preventDefault(){} });
    expect(props.select).toHaveBeenCalled();
  });

});
