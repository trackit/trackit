import React, {Component} from 'react';
import PropTypes from 'prop-types';
import Spinner from 'react-spinkit';

import {formatGigaBytes, formatPrice} from '../../../common/formatters';

// S3AnalyticsInfosComponent Component
class InfosComponent extends Component {

  extractTotals() {
    if (!this.props.data.hasOwnProperty("values") || !Object.keys(this.props.data.values).length)
      return null;

    const res = {
      buckets: 0,
      size: 0,
      bandwidth_cost: 0,
      storage_cost: 0,
      requests_cost: 0
    };

    Object.keys(this.props.data.values).forEach((key) => {
      const item = this.props.data.values[key];
      res.buckets++;
      res.size += item.GbMonth;
      res.bandwidth_cost += item.BandwidthCost;
      res.storage_cost += item.StorageCost;
      res.requests_cost += item.RequestsCost;
    });

    return res;
  }

  render() {
    if (!this.props.data || !this.props.data.status)
      return (<Spinner className="spinner clearfix" name='circle'/>);

    if (this.props.data && this.props.data.status && this.props.data.hasOwnProperty("error"))
      return (<div className="alert alert-warning" role="alert">Data not available ({this.props.data.error.message})</div>);

    const totals = this.extractTotals();

    if (!totals)
      return (<h4 className="no-data">No data available.</h4>);

    /* istanbul ignore next */
    return (
      <div>
        <div className="s3-card">
          <ul className="in-col">
            <li className="hidden-sm hidden-xs">
              <i className="fa fa-database card-icon blue-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatGigaBytes(totals.size)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            total size
          </h4>
        </div>
        <div className="s3-card">
          <ul className="in-col">
            <li className="hidden-sm hidden-xs">
              <i className="fa fa-shopping-bag card-icon blue-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {totals.buckets}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            total buckets
          </h4>
        </div>
        <div className="s3-card">
          <ul className="in-col">
            <li className="hidden-sm hidden-xs">
              <i className="fa fa-globe card-icon blue-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatPrice(totals.bandwidth_cost)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            bandwidth cost
          </h4>
        </div>
        <div className="s3-card">
          <ul className="in-col">
            <li className="hidden-sm hidden-xs">
              <i className="fa fa-hdd-o card-icon blue-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatPrice(totals.storage_cost)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            storage cost
          </h4>
        </div>
        <div className="s3-card">
          <ul className="in-col">
            <li className="hidden-sm hidden-xs">
              <i className="fa fa-exchange card-icon blue-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatPrice(totals.requests_cost)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            requests cost
          </h4>
        </div>
        <span className="clearfix"></span>
      </div>
    );
  }

}

InfosComponent.propTypes = {
  data: PropTypes.object,
};

InfosComponent.defaultProps = {
};

export default InfosComponent;
