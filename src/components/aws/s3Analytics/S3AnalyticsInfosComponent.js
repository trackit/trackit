import React, {Component} from 'react';
import PropTypes from 'prop-types';

import {formatBytes, formatPrice} from '../../../common/formatters';

// S3AnalyticsInfosComponent Component
class S3AnalyticsInfosComponent extends Component {

  extractTotals(data) {
    const res = {
      buckets: 0,
      size: 0,
      bandwidth_cost: 0,
      storage_cost: 0
    };

    for (let i = 0; i < data.length; i += 1) {
      const tmp = data[i];
      res.buckets++;
      res.size += tmp.size;
      res.bandwidth_cost += tmp.bw_cost;
      res.storage_cost += tmp.storage_cost;
    }
    return res;
  }

  render() {
    const totals = this.extractTotals(this.props.data);

    return (
      <div>
        <div className="col-md-3 col-sm-6 p-t-15 p-b-15 br-sm br-md bb-xs">
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
        <div className="col-md-3 col-sm-6 p-t-15 p-b-15 br-md bb-xs">
          <ul className="in-col">
            <li>
              <i className="fa fa-database fa-2x red-color"/>
            </li>
            <li>
              <h3 className="no-margin no-padding font-light">
                {formatBytes(totals.size)}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            total size
          </h4>
        </div>
        <div className="col-md-3 col-sm-6 p-t-15 p-b-15 bb-xs br-sm br-md">
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
        <div className="col-md-3 col-sm-6 p-t-15 p-b-15">
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
        <span className="clearfix"></span>
      </div>
    );
  }

}

S3AnalyticsInfosComponent.propTypes = {
  data: PropTypes.array.isRequired
};

export default S3AnalyticsInfosComponent;
