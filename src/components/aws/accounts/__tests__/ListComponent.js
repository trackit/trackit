import React from 'react';
import ListComponent, { ListItem } from '../ListComponent';
import { shallow } from 'enzyme';

describe('<ListComponent />', () => {

  const props = {
    delete: jest.fn(),
    accounts: []
  };

  const propsWithAccounts = {
    ...props,
    accounts: [{
      id: 42,
      userId: 42,
      roleArn: "arn:aws:iam::000000000001:role/TEST_ROLE",
      pretty: "Name"
    }, {
      id: 84,
      userId: 84,
      roleArn: "arn:aws:iam::000000000002:role2/TEST_ROLE_2",
      pretty: "Name2"
    }]
  };

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
    delete: jest.fn(),
    account: {
      id: 42,
      userId: 42,
      roleArn: "arn:aws:iam::000000000001:role/TEST_ROLE",
      pretty: "Name"
    }
  };

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

});
