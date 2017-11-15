import React from 'react';
import ListComponent, { ListItem } from '../ListComponent';
import { shallow } from 'enzyme';

const defaultProps = {
  new: jest.fn(),
  edit: jest.fn(),
  delete: jest.fn(),
  account: 42
};

const bill = {
  bucket: "s3://test.test",
  path: "/path/to/bill"
};

describe('<ListComponent />', () => {

  const props = {
    ...defaultProps,
    bills: []
  };

  const propsWithBills = {
    ...props,
    bills: [bill, bill]
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
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    const listWrapper = wrapper.find('ul');
    expect(listWrapper.length).toBe(1);
  });

  it('renders 2 <ListItem /> component when 2 bills are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    const list = wrapper.find(ListItem);
    expect(list.length).toBe(2);
  });

});

describe('<ListItem />', () => {

  const props = {
    ...defaultProps,
    bill
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

  it('can edit item', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    expect(props.edit.mock.calls.length).toBe(0);
    wrapper.find('button.btn.edit').prop('onClick')({ preventDefault() {} });
    expect(wrapper.state('editForm')).toBe(true);
    wrapper.instance().editBill(bill);
    expect(wrapper.state('editForm')).toBe(false);
//    expect(props.edit.mock.calls.length).toBe(1);
  });
/*
  it('can delete item', () => {
    const wrapper = shallow(<ListItem {...props}/>);
    expect(props.delete.mock.calls.length).toBe(0);
    wrapper.find('button.btn.delete').prop('onClick')({ preventDefault() {} });
    expect(props.delete.mock.calls.length).toBe(1);
  });
*/
});
