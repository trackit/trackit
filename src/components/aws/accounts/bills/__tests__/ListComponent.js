import React from 'react';
import { Item, ListComponent } from '../ListComponent';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Spinner from 'react-spinkit';
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
  error: "",
  bucket: "s3://test.test",
  prefix: "/path/to/bill"
};

const billWithError = {
  error: "access denied",
  nextPending: true,
  bucket: "another-billing-bucket",
  prefix: "another-prefix"
};

describe('<ListComponent />', () => {

  const props = {
    ...defaultProps,
  };

  const propsWithBills = {
    ...props,
    bills: {
      status: true,
      values: [bill, bill]
    }
  };

  const propsWaiting = {
    ...props,
    bills: {
      status: false
    }
  };

  const propsError = {
    ...props,
    bills: {
      status: true,
      error: Error("Error")
    }
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <ListComponent /> component', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div/> component when no bills is available', () => {
    const wrapper = shallow(<ListComponent {...propsError}/>);
    const alert = wrapper.find('div');
    expect(alert.length).toBe(1);
  });

  it('renders a <List/> component when bills are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    const listWrapper = wrapper.find(List);
    expect(listWrapper.length).toBe(1);
  });

  it('renders 2 <Item /> component when 2 bills are available', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    const list = wrapper.find(Item);
    expect(list.length).toBe(2);
  });

  it('renders a <Spinner /> component when bills are loading', () => {
    const wrapper = shallow(<ListComponent {...propsWaiting}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('can get bills', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
    expect(props.getBills).not.toHaveBeenCalled();
    wrapper.instance().getBills();
    expect(props.getBills).toHaveBeenCalled();
  });

  it('can clear bills', () => {
    const wrapper = shallow(<ListComponent {...propsWithBills}/>);
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

  const propsWithErrorInBills = {
    ...defaultProps,
    bill: billWithError
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

  it('renders two <ListItem/> component with one for error message ', () => {
    const wrapper = shallow(<Item {...propsWithErrorInBills}/>);
    const item = wrapper.find(ListItem);
    expect(item.length).toBe(2);
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
