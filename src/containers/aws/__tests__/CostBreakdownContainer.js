import React from 'react';
import { CostBreakdownContainer } from '../CostBreakdownContainer';
import Components from '../../../components';
import Moment from 'moment';
import { shallow } from 'enzyme';

const Chart = Components.AWS.CostBreakdown.Chart;
const Infos = Components.AWS.CostBreakdown.Infos;

describe('<CostBreakdownContainer />', () => {

  const props = {
    accounts: [],
    charts: {},
    costsValues: {},
    costsDates: {},
    costsInterval: {},
    costsFilter: {},
    initCharts: jest.fn(),
    addChart: jest.fn(),
    removeChart: jest.fn(),
    getCosts: jest.fn(),
    setCostsDates: jest.fn(),
    setCostsInterval: jest.fn(),
    setCostsFilter: jest.fn(),
    resetCostsDates: jest.fn(),
    resetCostsInterval: jest.fn(),
    resetCostsFilter: jest.fn(),
  };

  const propsWithInvalidCharts = {
    ...props,
    charts: {
      id: "bar",
      id2: "pie"
    }
  };

  const propsWithValidCharts = {
    ...propsWithInvalidCharts,
    costsDates: {
      id: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      },
      id2: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      }
    },
    costsInterval: {
      id: "interval",
      id2: "interval"
    },
    costsFilter: {
      id: "filter",
      id2: "filter"
    }
  };

  const propsWithSummary = {
    ...props,
    charts: {
      id: "summary"
    },
    costsDates: {
      id: {
        startDate: Moment().startOf('month'),
        endDate: Moment().endOf('month'),
      }
    },
    costsInterval: {
      id: "interval"
    },
    costsFilter: {
      id: "filter"
    }
  };

  const propsWithThreeCharts = {
    ...props,
    charts: {
      id: "bar",
      id2: "pie",
      id3: "bar"
    },
    costsDates: {
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
    costsInterval: {
      id: "interval",
      id2: "interval",
      id3: "interval"
    },
    costsFilter: {
      id: "filter",
      id2: "filter",
      id3: "filter"
    }
  };

  const propsWithThreeChartsAndSummary = {
    ...propsWithThreeCharts,
    charts: {
      id: "bar",
      id2: "pie",
      id3: "summary"
    }
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <CostBreakdownContainer /> component', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Chart/> component if data is available', () => {
    const wrapperInvalid = shallow(<CostBreakdownContainer {...propsWithInvalidCharts}/>);
    let charts = wrapperInvalid.find(Chart);
    expect(Object.keys(propsWithInvalidCharts.charts).length).toBe(2);
    expect(Object.keys(propsWithInvalidCharts.costsDates).length).toBe(0);
    expect(charts.length).toBe(0);
    const wrapperValid = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    charts = wrapperValid.find(Chart);
    expect(Object.keys(propsWithValidCharts.charts).length).toBe(2);
    expect(Object.keys(propsWithValidCharts.costsDates).length).toBe(2);
    expect(charts.length).toBe(2);
    const wrapperThree = shallow(<CostBreakdownContainer {...propsWithThreeCharts}/>);
    charts = wrapperThree.find(Chart);
    expect(Object.keys(propsWithThreeCharts.charts).length).toBe(3);
    expect(Object.keys(propsWithThreeCharts.costsDates).length).toBe(3);
    expect(charts.length).toBe(3);
  });

  it('renders <Infos/> compoeent if data is available', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithSummary}/>);
    let infos = wrapper.find(Infos);
    expect(infos.length).toBe(1);
    const wrapperThreeItems = shallow(<CostBreakdownContainer {...propsWithThreeChartsAndSummary}/>);
    infos = wrapperThreeItems.find(Infos);
    expect(infos.length).toBe(1);
  });

  it('generates default <Chart/> component if no chart available', () => {
    expect(props.initCharts).not.toHaveBeenCalled();
    shallow(<CostBreakdownContainer {...props}/>);
    expect(props.initCharts).toHaveBeenCalled();
  });

  it('can add a summary chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().addSummary({ preventDefault() {} });
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can add a bar chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().addBarChart({ preventDefault() {} });
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can add a pie chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().addPieChart({ preventDefault() {} });
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can add a diff chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().addDiffChart({ preventDefault() {} });
    expect(props.addChart).toHaveBeenCalled();
  });

  it('can reset charts', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.removeChart).not.toHaveBeenCalled();
    wrapper.instance().resetCharts({ preventDefault() {} });
    expect(props.removeChart).toHaveBeenCalled();
  });

  it('adds a chart when there is no chart', () => {
    const wrapper = shallow(<CostBreakdownContainer {...propsWithValidCharts}/>);
    expect(props.initCharts).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(props);
    expect(props.initCharts).toHaveBeenCalled();
  });

  it('does not add a chart when there is available charts', () => {
    const wrapper = shallow(<CostBreakdownContainer {...props}/>);
    expect(props.addChart).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsWithValidCharts);
    expect(props.addChart).not.toHaveBeenCalled();
  });

});
