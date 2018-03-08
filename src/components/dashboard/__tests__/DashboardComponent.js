import React from 'react';
import { DashboardComponent, Header } from '../DashboardComponent';
import AWS from '../aws';
import Moment from 'moment';
import { shallow } from 'enzyme';

const itemProps = {
  position: [0, 1],
  size: [3, 5],
  static: false,
  maxSize: [6, null],
  type: "cb_bar"
};

const props = {
  accounts: [],
  items: {},
  values: {},
  dates: {},
  intervals: {},
  filters: {},
  initDashboard: jest.fn(),
  updateDashboard: jest.fn(),
  addItem: jest.fn(),
  removeItem: jest.fn(),
  getData: jest.fn(),
  setItemDates: jest.fn(),
  setItemInterval: jest.fn(),
  setItemFilter: jest.fn(),
  resetItemDates: jest.fn(),
  resetItemInterval: jest.fn(),
  resetItemFilter: jest.fn(),
};

const propsWithInvalidItems = {
  ...props,
  items: {
    id: itemProps,
    id2: itemProps
  }
};

const propsWithValidItems = {
  ...propsWithInvalidItems,
  dates: {
    id: {
      startDate: Moment().startOf('month'),
      endDate: Moment().endOf('month'),
    },
    id2: {
      startDate: Moment().startOf('month'),
      endDate: Moment().endOf('month'),
    }
  },
  intervals: {
    id: "interval",
    id2: "interval"
  },
  filters: {
    id: "filter",
    id2: "filter"
  }
};

const propsWithThreeItems = {
  ...props,
  items: {
    id: itemProps,
    id2: itemProps,
    id3: itemProps
  },
  dates: {
    id: {
      startDate: Moment().startOf('month'),
      endDate: Moment().endOf('month'),
    },
    id2: {
      startDate: Moment().startOf('month'),
      endDate: Moment().endOf('month'),
    },
    id3: {
      startDate: Moment().startOf('month'),
      endDate: Moment().endOf('month'),
    }
  },
  intervals: {
    id: "interval",
    id2: "interval",
    id3: "interval"
  },
  filters: {
    id: "filter",
    id2: "filter",
    id3: "filter"
  }
};

const propsUpdateEmpty = {
  ...props,
  initDashboard: jest.fn()
};

const propsWithThreeItemsUpdate = {
  ...propsWithThreeItems,
  initDashboard: jest.fn()
};

describe('<DashboardComponent />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <DashboardComponent /> component', () => {
    const wrapper = shallow(<DashboardComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <ResponsiveReactGridLayout/> component', () => {
    const wrapper = shallow(<DashboardComponent {...props}/>);
    const grid = wrapper.find("WidthProvider.layout");
    expect(grid.length).toBe(1);
  });

  it('renders <Header/> component', () => {
    const wrapper = shallow(<DashboardComponent {...props}/>);
    const header = wrapper.find(Header);
    expect(header.length).toBe(1);
  });

  it('renders Dashboard Items components when available', () => {
    const wrapper = shallow(<DashboardComponent {...propsWithValidItems}/>);
    const items = wrapper.find("div.dashboard-item");
    expect(items.length).toBe(Object.keys(propsWithValidItems.items).length + 1);
  });

  it('can init Dashboard when empty', () => {
    expect(props.initDashboard).not.toHaveBeenCalled();
    expect(propsUpdateEmpty.initDashboard).not.toHaveBeenCalled();
    expect(propsWithThreeItemsUpdate.initDashboard).not.toHaveBeenCalled();
    const wrapper = shallow(<DashboardComponent {...props}/>);
    expect(props.initDashboard).toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsWithThreeItemsUpdate);
    expect(propsWithThreeItemsUpdate.initDashboard).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsUpdateEmpty);
    expect(propsUpdateEmpty.initDashboard).toHaveBeenCalled();
  });

  it('can add Item to Dashboard', () => {
    const wrapper = shallow(<DashboardComponent {...props}/>);
    expect(props.addItem).not.toHaveBeenCalled();
    wrapper.instance().addItem("mode");
    expect(props.addItem).toHaveBeenCalled();
  });

  it('can remove Item from Dashboard', () => {
    const wrapper = shallow(<DashboardComponent {...props}/>);
    expect(props.removeItem).not.toHaveBeenCalled();
    wrapper.instance().removeItem("key");
    expect(props.removeItem).toHaveBeenCalled();
  });

  it('can reset Dashboard', () => {
    const wrapper = shallow(<DashboardComponent {...propsWithValidItems}/>);
    expect(props.removeItem).not.toHaveBeenCalled();
    wrapper.instance().resetDashboard({ preventDefault() {} });
    expect(props.removeItem).toHaveBeenCalledTimes(2);
  });

  it('can update Dashboard layout', () => {
    const wrapper = shallow(<DashboardComponent {...propsWithValidItems}/>);
    const updatedLayout = [{
      ...itemProps, i: "id", x: 0, y: 0, w: 0, h: 0
    }, {
      ...itemProps, i: "wrongId", x: 0, y: 0, w: 0, h: 0
    }];
    expect(propsWithValidItems.updateDashboard).not.toHaveBeenCalled();
    wrapper.instance().updateLayout(updatedLayout);
    expect(propsWithValidItems.updateDashboard).toHaveBeenCalled();
  });

  describe('Items rendering', () => {
    const wrapper = shallow(<DashboardComponent {...propsWithValidItems}/>);
    const instance = wrapper.instance();

    it('can render a S3 Info component', () => {
      const output = instance.renderItem("id", {type: "s3_infos"});
      const item = shallow(output);
      const res = item.find(AWS.S3AnalyticsInfos);
      expect(res.length).toBe(1);
    });

    it('can render a S3 Info component', () => {
      const output = instance.renderItem("id", {type: "s3_chart"});
      const item = shallow(output);
      const res = item.find(AWS.S3AnalyticsCharts);
      expect(res.length).toBe(1);
    });

    it('can render a Cost Breakdown Info component', () => {
      const output = instance.renderItem("id", {type: "cb_infos"});
      const item = shallow(output);
      const res = item.find(AWS.CostBreakdownInfos);
      expect(res.length).toBe(1);
    });

    it('can render a Cost Breakdown Pie Chart component', () => {
      const output = instance.renderItem("id", {type: "cb_pie"});
      const item = shallow(output);
      const res = item.find(AWS.CostBreakdownPieChart);
      expect(res.length).toBe(1);
    });

    it('can render a Cost Breakdown Bar Chart component', () => {
      const output = instance.renderItem("id", {type: "cb_bar"});
      const item = shallow(output);
      const res = item.find(AWS.CostBreakdownBarChart);
      expect(res.length).toBe(1);
    });

    it('can render empty item if data is not available', () => {
      const output = instance.renderItem("wrongId", {type: "s3_info"});
      const item = shallow(output);
      const res = item.find(AWS.S3AnalyticsInfos);
      expect(res.length).toBe(0);
    });

    it('can render empty item if invalid type', () => {
      const output = instance.renderItem("id", {type: "invalid_type"});
      const item = shallow(output);
      const res = item.find(".dashboard-item-content");
      expect(res.length).toBe(1);
    });

  });

});

/*
describe('<Header />', () => {

  it('renders <TimerangeSelector/> component when BarChart is selected', () => {
    const wrapper = shallow(<Header {...props}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(1);
    const interval = wrapper.find(IntervalNavigator);
    expect(interval.length).toBe(0);
  });

  it('renders <IntervalNavigator/> component when PieChart is selected', () => {
    const wrapper = shallow(<Header {...propsPie}/>);
    const timerange = wrapper.find(TimerangeSelector);
    expect(timerange.length).toBe(0);
    const interval = wrapper.find(IntervalNavigator);
    expect(interval.length).toBe(1);
  });

  it('renders <Selector/> component', () => {
    const wrapper = shallow(<Header {...props}/>);
    const selector = wrapper.find(Selector);
    expect(selector.length).toBe(1);
  });

  it('renders a <Spinner/> component when data is unavailable', () => {
    const wrapper = shallow(<Header {...propsWithoutData}/>);
    const spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders a <button/> component when can be closed', () => {
    const wrapper = shallow(<Header {...propsWithClose}/>);
    const button = wrapper.find("button");
    expect(button.length).toBe(1);
  });

  it('renders an alert component when there is an error', () => {
    const wrapper = shallow(<Header {...propsWithError}/>);
    const alert = wrapper.find("div.alert");
    expect(alert.length).toBe(1);
  });

  it('can set dates', () => {
    const wrapper = shallow(<Header {...props}/>);
    expect(props.setDates).not.toHaveBeenCalled();
    wrapper.instance().setDates(Moment().startOf('month'), Moment().endOf('month'));
    expect(props.setDates).toHaveBeenCalledTimes(1);
  });

  it('can set filter', () => {
    const wrapper = shallow(<Header {...props}/>);
    expect(props.setFilter).not.toHaveBeenCalled();
    wrapper.instance().setFilter("filter");
    expect(props.setFilter).toHaveBeenCalledTimes(1);
  });

  it('can set interval', () => {
    const wrapper = shallow(<Header {...props}/>);
    expect(props.setInterval).not.toHaveBeenCalled();
    wrapper.instance().setInterval("interval");
    expect(props.setInterval).toHaveBeenCalledTimes(1);
  });

  it('can close', () => {
    const wrapper = shallow(<Header {...propsWithClose}/>);
    expect(propsWithClose.close).not.toHaveBeenCalled();
    wrapper.instance().close({ preventDefault() {} });
    expect(propsWithClose.close).toHaveBeenCalledTimes(1);
  });

});
*/