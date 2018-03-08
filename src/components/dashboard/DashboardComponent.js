import React, { Component } from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import { Responsive, WidthProvider } from 'react-grid-layout';
import UUID from 'uuid/v4';
import Actions from '../../actions';
import AWS from './aws';

import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';

import s3square from '../../assets/s3-square.png';

const ResponsiveReactGridLayout = WidthProvider(Responsive);

const defaultValues = {
  position: [0,0],
  static: false,
  maxSize: [6, undefined]
};

const itemsSize = {
  header: [6, 1],
  cb_infos: [6,2],
  cb_pie: [3,4],
  cb_bar: [3,4],
  s3_infos: [6,2],
  s3_chart: [2,4],
};

const generateLayout = (item) => {
  return {
    x: (item.hasOwnProperty("position") ? item.position[0] : defaultValues.position[0]),
    y: (item.hasOwnProperty("position") ? item.position[1] : defaultValues.position[1]),
    w: (itemsSize.hasOwnProperty(item.type) ? itemsSize[item.type][0] : 1),
    h: (itemsSize.hasOwnProperty(item.type) ? itemsSize[item.type][1] : 1),
    static: (item.hasOwnProperty("static") ? item.static : defaultValues.static),
    isResizable: false
  };
};

/* istanbul ignore next */
const renderItem = (key, item, child, close=null) => {
  const layout = generateLayout(item);
  let title;
  switch (item.type) {
    case "cb_infos":
    case "cb_pie":
    case "cb_bar":
      title = (
        <div className=" dashboard-item-icon">
          <i className="menu-icon fa fa-area-chart red-color"/>
          &nbsp;
          Cost Breakdown
        </div>
      );
      break;
    case "s3_infos":
    case "s3_chart":
      title = (
        <div className=" dashboard-item-icon">
          <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
          &nbsp;
          S3 Analytics
        </div>
      );
      break;
    default:
      title = null;
  }
  const closeButton = (close !== null ? (
    <div className="close" onClick={(e) => { e.preventDefault(); close(key); }}>
      <i className="fa fa-times close" aria-hidden="true"/>
    </div>
  ) : null);

  return (
    <div className="white-box dashboard-item" key={key} data-grid={layout}>
      {title}
      {closeButton}
      <div className="dashboard-item-content">
        {child}
      </div>
    </div>
  );
};

export class Header extends Component {
  /* istanbul ignore next */
  render() {
    return (
      <div>
        <div className="clearfix">
          <div className="inline-block">
            <h3 className="white-box-title no-padding inline-block">Dashboard</h3>
          </div>
        </div>
        <div className="clearfix">
          <div className="inline-block pull-right">
            <div className="inline-block dashboard-btn-group">
              <div className="inline-block dashboard-btn-group-title">
                <i className="menu-icon fa fa-area-chart red-color"/>
                &nbsp;
                Cost Breakdown
              </div>
              &nbsp;
              <div className="btn-group">
                <button className="btn btn-default inline-block" onClick={(e) => {e.preventDefault(); this.props.addItem("cb_infos");}}>Info</button>
                <button className="btn btn-default inline-block" onClick={(e) => {e.preventDefault(); this.props.addItem("cb_bar");}}>Bar Chart</button>
                <button className="btn btn-default inline-block" onClick={(e) => {e.preventDefault(); this.props.addItem("cb_pie");}}>Pie Chart</button>
              </div>
            </div>
            &nbsp;
            <div className="inline-block dashboard-btn-group">
              <div className="inline-block dashboard-btn-group-title">
                <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
                &nbsp;
                S3 Analytics
              </div>
              &nbsp;
              <div className="btn-group">
                <button className="btn btn-default inline-block" onClick={(e) => {e.preventDefault(); this.props.addItem("s3_infos");}}>Info</button>
                <button className="btn btn-default inline-block" onClick={(e) => {e.preventDefault(); this.props.addItem("s3_chart");}}>Chart</button>
              </div>
            </div>
            &nbsp;
            <button className="btn btn-danger inline-block dashboard-btn-group" onClick={this.props.reset}>Reset dashboard</button>
          </div>
        </div>
      </div>
    );
  }
}

Header.propsTypes = {
  addItem: PropTypes.func.isRequired,
  reset: PropTypes.func.isRequired,
};

const header = {
  type: "header",
  position: [0, 0],
  size: [6, 1],
  static: true
};

// Dashboard Component
export class DashboardComponent extends Component {

  constructor(props) {
    super(props);
    if (!this.props.items || !Object.keys(this.props.items).length)
      this.props.initDashboard();
    this.addItem = this.addItem.bind(this);
    this.removeItem = this.removeItem.bind(this);
    this.updateLayout = this.updateLayout.bind(this);
    this.renderItem = this.renderItem.bind(this);
    this.resetDashboard = this.resetDashboard.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    if (!Object.keys(nextProps.items).length)
      nextProps.initDashboard();
  }

  addItem = (mode) => {
    this.props.addItem({
      ...defaultValues,
      type: mode
    });
  };

  removeItem = (key) => {
    this.props.removeItem(key);
  };

  resetDashboard = (e) => {
    e.preventDefault();
    Object.keys(this.props.items).forEach((id) => {this.props.removeItem(id)});
  };

  updateLayout = (layout) => {
    let items = Object.assign({}, this.props.items);
    layout.forEach((item) => {
      if (items.hasOwnProperty(item.i)) {
        let data = items[item.i];
        data.position = [item.x, item.y];
        data.size = [item.w, item.h];
        items[item.i] = data;
      }
    });
    this.props.updateDashboard(items);
  };

  renderItem = (key, item) => {
    let content = null;
    if (this.props.values &&
      this.props.dates && this.props.dates.hasOwnProperty(key) &&
      this.props.intervals && this.props.intervals.hasOwnProperty(key) &&
      this.props.filters && this.props.filters.hasOwnProperty(key)
    ) {
      switch (item.type) {
        case "s3_infos":
          content = <AWS.S3AnalyticsInfos
            id={key}
            accounts={this.props.accounts}
            values={this.props.values[key]}
            getValues={this.props.getData}
            dates={this.props.dates[key]}
            setDates={this.props.setItemDates}
          />;
          break;
        case "s3_chart":
          content = <AWS.S3AnalyticsCharts
            id={key}
            accounts={this.props.accounts}
            values={this.props.values[key]}
            getValues={this.props.getData}
            dates={this.props.dates[key]}
            setDates={this.props.setItemDates}
            filter={this.props.filters[key]}
            setFilter={this.props.setItemFilter}
          />;
          break;
        case "cb_infos":
          content = <AWS.CostBreakdownInfos
            id={key}
            accounts={this.props.accounts}
            values={this.props.values[key]}
            getValues={this.props.getData}
            dates={this.props.dates[key]}
            setDates={this.props.setItemDates}
            interval={this.props.intervals[key]}
            setInterval={this.props.setItemInterval}
          />;
          break;
        case "cb_pie":
          content = <AWS.CostBreakdownPieChart
            id={key}
            accounts={this.props.accounts}
            values={this.props.values[key]}
            getValues={this.props.getData}
            dates={this.props.dates[key]}
            setDates={this.props.setItemDates}
            filter={this.props.filters[key]}
            setFilter={this.props.setItemFilter}
            interval={this.props.intervals[key]}
            setInterval={this.props.setItemInterval}
          />;
          break;
        case "cb_bar":
          content = <AWS.CostBreakdownBarChart
            id={key}
            accounts={this.props.accounts}
            values={this.props.values[key]}
            getValues={this.props.getData}
            dates={this.props.dates[key]}
            setDates={this.props.setItemDates}
            filter={this.props.filters[key]}
            setFilter={this.props.setItemFilter}
            interval={this.props.intervals[key]}
            setInterval={this.props.setItemInterval}
          />;
          break;
        default:
          content = key;
      }
    }
    return renderItem(key, item, content, this.removeItem)
  };

  render() {
    return (
      <div className="container-fluid">

        <ResponsiveReactGridLayout
          className="layout"
          containerPadding={[0,0]}
          cols={{lg: 6, md: 6, sm: 6, xs: 3, xxs: 3}}
          onLayoutChange={this.updateLayout}
          rowHeight={100}
        >
          {renderItem("header", header, (<Header addItem={this.addItem} reset={this.resetDashboard}/>))}
          {Object.keys(this.props.items).map(key => this.renderItem(key, this.props.items[key]))}
        </ResponsiveReactGridLayout>

      </div>
    );
  }

}

DashboardComponent.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  items: PropTypes.object,
  values: PropTypes.object,
  dates: PropTypes.object,
  intervals: PropTypes.object.isRequired,
  filters: PropTypes.object.isRequired,
  initDashboard: PropTypes.func.isRequired,
  updateDashboard: PropTypes.func.isRequired,
  addItem: PropTypes.func.isRequired,
  removeItem: PropTypes.func.isRequired,
  getData: PropTypes.func.isRequired,
  setItemDates: PropTypes.func.isRequired,
  setItemInterval: PropTypes.func.isRequired,
  setItemFilter: PropTypes.func.isRequired,
  resetItemDates: PropTypes.func.isRequired,
  resetItemInterval: PropTypes.func.isRequired,
  resetItemFilter: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws, dashboard}) => ({
  items: dashboard.items,
  values: dashboard.values,
  dates: dashboard.dates,
  intervals: dashboard.intervals,
  filters: dashboard.filters,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  initDashboard: () => {
    dispatch(Actions.Dashboard.initDashboard());
  },
  updateDashboard: (layout) => {
    dispatch(Actions.Dashboard.updateDashboard(layout));
  },
  addItem: (props) => {
    dispatch(Actions.Dashboard.addItem(UUID(), props));
  },
  removeItem: (id) => {
    dispatch(Actions.Dashboard.removeItem(id));
  },
  getData: (id, type, begin, end, filters) => {
    dispatch(Actions.Dashboard.getData(id, type, begin, end, filters));
  },
  setItemDates: (id, startDate, endDate) => {
    dispatch(Actions.Dashboard.setItemDates(id, startDate, endDate))
  },
  resetItemDates: () => {
    dispatch(Actions.Dashboard.resetItemDates())
  },
  setItemInterval: (id, interval) => {
    dispatch(Actions.Dashboard.setItemInterval(id, interval));
  },
  resetItemInterval: () => {
    dispatch(Actions.Dashboard.resetItemInterval());
  },
  setItemFilter: (id, filter) => {
    dispatch(Actions.Dashboard.setItemFilter(id, filter));
  },
  resetItemFilter: () => {
    dispatch(Actions.Dashboard.resetItemFilter());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(DashboardComponent);