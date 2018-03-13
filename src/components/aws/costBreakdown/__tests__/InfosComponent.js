import React from 'react';
import InfosComponent from '../InfosComponent';
import { shallow } from 'enzyme';
import Spinner from "react-spinkit";
import Moment from "moment/moment";
import Misc from '../../../misc';

const TimerangeSelector = Misc.TimerangeSelector;

const props = {
  id: "42",
  type: "bar",
  values: {},
  dates: {
    startDate: Moment().startOf('month'),
    endDate: Moment(),
  },
  filter: "product",
  getCosts: jest.fn(),
  setDates: jest.fn(),
  setFilter: jest.fn(),
};

const propsNoDates = {
  ...props,
  setDates: undefined
};

const propsLoading = {
  ...props,
  values: {
    status: false
  }
};

const propsWithData = {
  ...props,
  values: {
    status: true,
    values: {
      region: {
        region1: {product: {product1: 42}},
        region2: {product: {product2: 84, product1: 0}}
      }
    }
  }
};

const propsWithWrongData = {
  ...props,
  values: {
    status: true,
    values: {
      test: 42
    }
  }
};

const propsWithError = {
  ...props,
  values: {
    status: true,
    error: Error()
  }
};

const updatedDateProps = {
  ...props,
  dates: {
    startDate: Moment().startOf('year'),
    endDate: Moment(),
  },
  getCosts: jest.fn()
};

const propsWithClose = {
  ...props,
  close: jest.fn()
};

describe('<InfosComponent />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <InfosComponent /> component', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Spinner /> component when data is not available', () => {
    const wrapper = shallow(<InfosComponent {...propsLoading}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders a <TimerangeSelector /> component when setDates is available', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    const selector = wrapper.find(TimerangeSelector);
    expect(selector.length).toBe(1);
  });

  it('renders no <TimerangeSelector /> component when setDates is not available', () => {
    const wrapper = shallow(<InfosComponent {...propsNoDates}/>);
    const selector = wrapper.find(TimerangeSelector);
    expect(selector.length).toBe(0);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<InfosComponent {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('calculates totals based on data', () => {
    const wrapper = shallow(<InfosComponent {...propsWithData}/>);
    const totals = wrapper.instance().extractTotals();
    expect(totals.cost).toBe(propsWithData.values.values.region.region1.product.product1 + propsWithData.values.values.region.region2.product.product2);
    expect(totals.services).toBe(2);
    expect(totals.regions).toBe(2);
  });

  it('calculates nothing when there is no data', () => {
    const wrapper = shallow(<InfosComponent {...propsWithWrongData}/>);
    const totals = wrapper.instance().extractTotals();
    expect(totals).toBe(null);
  });

  it('loads costs when mounting', () => {
    expect(props.getCosts).not.toHaveBeenCalled();
    shallow(<InfosComponent {...props}/>);
    expect(props.getCosts).toHaveBeenCalled();
  });

  it('reloads costs when dates are updated', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    expect(updatedDateProps.getCosts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(updatedDateProps);
    expect(updatedDateProps.getCosts).toHaveBeenCalled();
  });

  it('does not reload when dates nor filters are updated', () => {
    const wrapper = shallow(<InfosComponent {...props}/>);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.getCosts).toHaveBeenCalledTimes(1);
  });

  it('can close', () => {
    const wrapper = shallow(<InfosComponent {...propsWithClose}/>);
    expect(propsWithClose.close).not.toHaveBeenCalled();
    wrapper.instance().close({ preventDefault() {} });
    expect(propsWithClose.close).toHaveBeenCalledTimes(1);
  });

});
