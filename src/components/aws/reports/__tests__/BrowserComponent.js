import React from 'react';
import { BrowserComponent } from '../BrowserComponent';
import { shallow, mount } from 'enzyme';
import Spinner from "react-spinkit";
import ReactTable from 'react-table';

const props = {
  startDownload: jest.fn(),
  account: '420'
};

const propsWithReports = {
  ...props,
  reportList: {
    status: true,
    values: ['mytest/test.xlsx',]
  }
}

const propsLoading = {
  ...props,
  reportList: {
    status: false,
    values: []
  }
};

const propsWithError = {
  ...props,
  reportList: {
    status: true,
    error: Error()
  }
};

describe('<BrowserComponent />', () => {
  it('renders a <BrowserComponent /> component', () => {
    const wrapper = shallow(<BrowserComponent {...propsWithReports}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <BrowserComponent /> component', () => {
    const wrapper = shallow(<BrowserComponent {...propsWithReports}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is not available', () => {
    const wrapper = shallow(<BrowserComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<BrowserComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('calls startDownload when clicking on a report', () => {
    const wrapper = mount(<BrowserComponent {...propsWithReports}/>);
    wrapper.find("button").at(0).simulate('click');
    expect(props.startDownload).toHaveBeenCalled();
  });
});
