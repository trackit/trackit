import React, {Component} from 'react';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';
import Misc from '../../misc';
import {formatPrice} from '../../../common/formatters';

const TimerangeSelector = Misc.TimerangeSelector;

class InfosComponent extends Component {

  constructor(props) {
    super(props);
    this.close = this.close.bind(this);
    this.setDates = this.setDates.bind(this);
  }

  componentWillMount() {
    this.props.getCosts(this.props.id, this.props.dates.startDate, this.props.dates.endDate, ['region', 'product']);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.dates !== nextProps.dates ||
      this.props.accounts !== nextProps.accounts)
      nextProps.getCosts(nextProps.id, nextProps.dates.startDate, nextProps.dates.endDate, ['region', 'product']);
  }

  close = (e) => {
    e.preventDefault();
    this.props.close(this.props.id);
  };

  setDates = (start, end) => {
    this.props.setDates(this.props.id, start, end);
  };

  extractTotals() {
    if (!this.props.values.values.hasOwnProperty("region"))
      return null;

    const res = {
      cost: 0,
      services: 0,
      regions: 0
    };

    let products = [];

    Object.keys(this.props.values.values.region).forEach((key) => {
      const item = this.props.values.values.region[key];
      res.regions++;
      Object.keys(item.product).forEach((name) => {
        if (products.indexOf(name) < 0)
          products.push(name);
        res.cost += item.product[name];
      });
    });

    res.services = products.length;

    return res;
  }

  render() {
    const icon = (this.props.icon ? (
      <div className="dashboard-item-icon">
        <i className="fa fa-list"/>
        &nbsp;
        Summary
      </div>
    ) : null);

    const loading = (!this.props.values || !this.props.values.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

    const close = (this.props.close ? (
      <button className="btn btn-danger" onClick={this.close}>Remove this chart</button>
    ) : null);

    const error = (this.props.values && this.props.values.status && this.props.values.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">Data not available ({this.props.values.error.message})</div>
    ) : null);

    const timerange = (this.props.setDates ? (
      <TimerangeSelector
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.setDates}
      />
    ) : null);

    const totals = (this.props.values && this.props.values.status && this.props.values.hasOwnProperty("values") ? this.extractTotals() : null);

    const noData = (!totals && !loading ? (<h4 className="no-data">No data available.</h4>) : null);

    let values = null;

    if (totals && !noData) {
      values = (
        <div>
          <div className="cb-card">
            <ul className="in-col">
              <li>
                <i className="fa fa-credit-card card-icon blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {formatPrice(totals.cost)}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              total cost
            </h4>
          </div>
          <div className="cb-card">
            <ul className="in-col">
              <li>
                <i className="fa fa-th-list card-icon blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {totals.services}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              services
            </h4>
          </div>
          <div className="cb-card">
            <ul className="in-col">
              <li>
                <i className="fa fa-globe card-icon blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {totals.regions}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              regions
            </h4>
          </div>
          <span className="clearfix"></span>
        </div>
      );
    }

    const metaBloc = (
      <div className="clearfix">
      <div className="inline-block pull-left">
        {icon}
        {loading}
        {error}
      </div>
      <div className="inline-block pull-right">
        <div className="inline-block">
          {timerange}
        </div>
        &nbsp;
        {close}
      </div>
    </div>
    );

    return (
      <div>
        {(icon || loading || error || timerange || close) && metaBloc}
        {noData}
        {values}
      </div>
    );
  }

}

InfosComponent.propTypes = {
  id: PropTypes.string.isRequired,
  values: PropTypes.object,
  dates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  accounts: PropTypes.arrayOf(PropTypes.object),
  getCosts: PropTypes.func.isRequired,
  setDates: PropTypes.func,
  close: PropTypes.func,
  icon: PropTypes.bool
};

InfosComponent.defaultProps = {
  icon: true
};

export default InfosComponent;
