import React, {Component} from 'react';
import PropTypes from "prop-types";
import Misc from "../../misc";
import Spinner from "react-spinkit";
import TagsChart from './TagsChartComponent';

const TimerangeSelector = Misc.TimerangeSelector;
const Selector = Misc.Selector;
const Tooltip = Misc.Popover;

const filters = {
  product: "Product",
  region: "Region",
  availabilityzone: "Availability Zone",
  account: "Account",
};

class Header extends Component {

  constructor(props) {
    super(props);
    this.close = this.close.bind(this);
    this.setDates = this.setDates.bind(this);
    this.selectTag = this.selectTag.bind(this);
    this.selectFilter = this.selectFilter.bind(this);
  }

  close = (e) => {
    e.preventDefault();
    this.props.close(this.props.id);
  };

  setDates = (start, end) => {
    this.props.setDates(this.props.id, start, end);
  };

  selectTag = (tag) => {
    this.props.selectKey(this.props.id, tag);
  };

  selectFilter = (filter) => {
    this.props.selectFilter(this.props.id, filter);
  };

  render() {
    const close = (this.props.close ? (
      <button className="btn btn-danger" onClick={this.close}><i className="fa fa-times"/></button>
    ) : null);

    let loading = null;
    let keys = (!this.props.keys || !this.props.keys.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);
    if (this.props.keys && this.props.keys.status && this.props.keys.hasOwnProperty("values") && this.props.keys.values.length) {
      const values = {};
      this.props.keys.values.forEach((key) => {values[key] = key;});
      keys = (
        <Selector
          values={values}
          selected={this.props.tag}
          selectValue={this.selectTag}
        />
      );
      loading = (!this.props.values || !this.props.values.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);
    }

    const tooltip = (
      <div className="inline-block">
        <Tooltip icon={<i className="fa fa-info-circle btn btn-default"/>} tooltip="To hide/display a data series please click on its name in the legend. Double-click will display only this item." placement="left"/>
      </div>
    );

    return (
      <div className="clearfix">

        <div className="inline-block pull-left">
          <div className="dashboard-item-icon">
            <i className="menu-icon fa fa-pie-chart"/>
            &nbsp;
            Pie Chart
          </div>
          {loading}
        </div>

        <div className="inline-block pull-right">

          {tooltip}

          {keys}

          <div className="inline-block">
            <Selector
              values={filters}
              selected={this.props.filter}
              selectValue={this.selectFilter}
            />
          </div>

          <div className="inline-block">
            <TimerangeSelector
              startDate={this.props.dates.startDate}
              endDate={this.props.dates.endDate}
              setDatesFunc={this.setDates}
            />
          </div>

          {close}

        </div>

      </div>
    );
  }

}

Header.propTypes = {
  id: PropTypes.string.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  keys: PropTypes.object,
  getKeys: PropTypes.func.isRequired,
  tag: PropTypes.string.isRequired,
  selectKey: PropTypes.func.isRequired,
  values: PropTypes.object,
  getValues: PropTypes.func.isRequired,
  dates: PropTypes.object.isRequired,
  setDates: PropTypes.func.isRequired,
  filter: PropTypes.string.isRequired,
  selectFilter: PropTypes.func.isRequired,
  close: PropTypes.func
};

class ChartComponent extends Component {

  componentWillMount() {
    this.props.getKeys(this.props.id, this.props.dates.startDate, this.props.dates.endDate);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts) {
      nextProps.getKeys(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate);
      nextProps.selectKey(nextProps.id, "");
    }
    else if (this.props.keys !== nextProps.keys && nextProps.keys.status
      && nextProps.keys.hasOwnProperty("values") && nextProps.keys.values.length)
      nextProps.selectKey(nextProps.id, nextProps.keys.values[0]);
    else if ((this.props.tag !== nextProps.tag && nextProps.tag !== "") || (this.props.filter !== nextProps.filter && nextProps.filter !== ""))
      nextProps.getValues(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate, nextProps.filter, nextProps.tag);
  }

  render() {
    let error = null;

    if (this.props.keys && this.props.keys.status) {
      if (this.props.keys.hasOwnProperty("error"))
        error = (<div className="alert alert-warning m-t-20" role="alert">Data not available ({this.props.keys.error.message})</div>);
      else if (this.props.keys.hasOwnProperty("values") && !this.props.keys.values.length)
        error = (<div className="alert alert-warning m-t-20" role="alert">You need to set tags (No keys available for this timerange)</div>);
    }
    if (this.props.values && this.props.keys.values && this.props.keys.hasOwnProperty("error"))
      error = (<div className="alert alert-warning m-t-20" role="alert">Data not available ({this.props.keys.error.message})</div>);

    const chart = (error === null && this.props.values && this.props.values.status && this.props.values.hasOwnProperty("values") ? (
      <TagsChart
        values={this.props.values.values}
        legend
        height={450}
        filter={filters[this.props.filter]}
      />
    ) : null);

    return (
      <div className="clearfix">
        <Header {...this.props}/>
        {error}
        {chart}
      </div>
    )
  }

}

ChartComponent.propTypes = {
  id: PropTypes.string.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  keys: PropTypes.object,
  getKeys: PropTypes.func.isRequired,
  tag: PropTypes.string.isRequired,
  selectKey: PropTypes.func.isRequired,
  values: PropTypes.object,
  getValues: PropTypes.func.isRequired,
  dates: PropTypes.object.isRequired,
  setDates: PropTypes.func.isRequired,
  filter: PropTypes.string.isRequired,
  selectFilter: PropTypes.func.isRequired,
  close: PropTypes.func
};

export default ChartComponent;