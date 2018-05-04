import React from 'react';
import { ReportsContainer } from '../ReportsContainer';
import { shallow } from 'enzyme';

const baseProps = {
  account: '42',
  reportList: ['myreports/myfile.xlsx',],
  getAccounts: jest.fn(),
  selectAccount: jest.fn(),
  requestGetReports: jest.fn()
};

const validProps = {
  ...baseProps,
  accounts: {
    status: true,
    values: [{
      id: 42,
      roleArn: 'rolearn',
      pretty: 'pretty',
      bills: []
    }],
  },
  downloadStatus: {
    failed: false,
  },
};

const accountsMissingProps = {
  ...baseProps,
  accounts: {
    status: false,
  },
  downloadStatus: {
    failed: false,
  },
}

const reportsMissingProps = {
  ...baseProps,
  accounts: {
    status: true,
  },
  reportList: {
    status: false,
  },
  downloadStatus: {
    failed: false,
  },
}

const downloadErrorProps = {
  ...baseProps,
  accounts: {
    status: true,
    values: [{
      id: 42,
      roleArn: 'rolearn',
      pretty: 'pretty',
      bills: []
    }],
  },
  downloadStatus: {
    failed: true,
    error: Error()
  },
}

describe('<ReportsContainer />', () => {
  it('renders a <ReportsContainer /> component', () => {
    const wrapper = shallow(<ReportsContainer {...validProps}/>);
    expect(wrapper.length).toBe(1);
  });

  it('Calls getAccounts when accounts are missing', () => {
    shallow(<ReportsContainer {...accountsMissingProps}/>);
    expect(baseProps.getAccounts).toHaveBeenCalled();
  });

  it('Calls requestGetReports when accounts are missing', () => {
    const wrapper = shallow(<ReportsContainer {...reportsMissingProps}/>);
    wrapper.instance().componentDidUpdate(reportsMissingProps);
    expect(baseProps.requestGetReports).toHaveBeenCalled();
  });

  it('displays an error when the download fails', () => {
    const wrapper = shallow(<ReportsContainer {...downloadErrorProps}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });
});
