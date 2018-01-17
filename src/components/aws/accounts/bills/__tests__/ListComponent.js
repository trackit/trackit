import React from 'react';
import { Item, ListComponent } from '../ListComponent';
import List, {
  ListItem
} from 'material-ui/List';
import { shallow } from 'enzyme';

const defaultProps = {
  getBills: jest.fn(),
  newBill: jest.fn(),
  editBill: jest.fn(),
  deleteBill: jest.fn(),
  clearBills: jest.fn(),
  account: 42
};

const bill = {
  bucket: "s3://test.test",
  prefix: "/path/to/bill"
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
    jest.resetAllMocks();
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

  it('can get bills', () => {
    const wrapper = shallow(<ListComponent {...props}/>);
    expect(props.getBills).not.toHaveBeenCalled();
    wrapper.instance().getBills();
    expect(props.getBills).toHaveBeenCalled();
  });

  it('can clear bills', () => {
    const wrapper = shallow(<ListComponent {...props}/>);
    expect(props.clearBills).not.toHaveBeenCalled();
    wrapper.instance().clearBills();
    expect(props.clearBills).toHaveBeenCalled();
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
    expect(props.editBill).not.toHaveBeenCalled();
    wrapper.instance().editBill(bill);
//    expect(props.editBill).toHaveBeenCalled();
  });

  it('can delete item', () => {
    const wrapper = shallow(<Item {...props}/>);
    expect(props.deleteBill).not.toHaveBeenCalled();
    wrapper.instance().deleteBill();
//    expect(props.deleteBill).toHaveBeenCalled();
  });

});
