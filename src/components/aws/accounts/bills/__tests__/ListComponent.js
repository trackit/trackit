import React from 'react';
import ListComponent, { Item } from '../ListComponent';
import List, {
  ListItem
} from 'material-ui/List';
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

  beforeEach(() => {
    defaultProps.new.mockReset();
    defaultProps.edit.mockReset();
    defaultProps.delete.mockReset();
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
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    const listWrapper = wrapper.find(List);
    expect(listWrapper.length).toBe(1);
  });

  it('renders 2 <Item /> component when 2 bills are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    const list = wrapper.find(Item);
    expect(list.length).toBe(2);
  });

});

describe('<Item />', () => {

  const props = {
    ...defaultProps,
    bill
  };

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
    expect(props.edit.mock.calls.length).toBe(0);
    wrapper.instance().editBill(bill);
    expect(props.edit.mock.calls.length).toBe(1);
  });

  it('can delete item', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(props.delete.mock.calls.length).toBe(0);
    wrapper.instance().deleteBill();
    expect(props.delete.mock.calls.length).toBe(1);
  });

});
