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
      return (<h2>No data available.</h2>);

    return (
      <div>
        <div className={"col-md-2 col-sm-6 p-t-15 p-b-15 br-sm br-md bb-xs" + (this.props.offset ? " col-md-offset-1" : "")}>
          <ul className="in-col">
            <li>
              <i className="fa fa-shopping-bag fa-2x green-color"/>
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
        <div className="col-md-2 col-sm-6 p-t-15 p-b-15 br-md bb-xs">
          <ul className="in-col">
            <li>
              <i className="fa fa-database fa-2x red-color"/>
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
        <div className="col-md-2 col-sm-4 p-t-15 p-b-15 bb-xs br-sm br-md">
          <ul className="in-col">
            <li>
              <i className="fa fa-globe fa-2x blue-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatPrice(totals.bandwidth_cost)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            bandwidth total cost
          </h4>
        </div>
        <div className="col-md-2 col-sm-4 p-t-15 p-b-15 bb-xs br-sm">
          <ul className="in-col">
            <li>
              <i className="fa fa-hdd-o fa-2x orange-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatPrice(totals.storage_cost)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            storage total cost
          </h4>
        </div>
        <div className="col-md-2 col-sm-4 p-t-15 p-b-15">
          <ul className="in-col">
            <li>
              <i className="fa fa-exchange fa-2x purple-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatPrice(totals.requests_cost)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            requests total cost
          </h4>
        </div>
        <span className="clearfix"></span>
      </div>
    );
  }

}

InfosComponent.propTypes = {
  data: PropTypes.object,
  offset: PropTypes.bool
};

InfosComponent.defaultProps = {
  offset: true
};

export default InfosComponent;
