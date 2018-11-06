import React, { Component } from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';
import { formatPrice } from '../../common/formatters';

const getTotalCost = (costs) => {
  let total = 0;
  Object.keys(costs).forEach((key) => total += costs[key]);
  return total;
};

class SummaryComponent extends Component {
  getMonthTotal(products) {
    let res = 0;
    for (const key in products) {
      if (products.hasOwnProperty(key)) {
        const element = products[key];
        res += element;
      }
    }
    return res;
  }

  getPotentialsSavings(unused) {
    let res = 0;

    if (unused.ec2 && unused.ec2.status && unused.ec2.values) {
      for (let i = 0; i < unused.ec2.values.length; i++) {
        const element = unused.ec2.values[i];
        res += getTotalCost(element.instance.costs);
      }
    }
    return res;
  }

  render() {
    let message;
    let monthCost;
    let variation;

    const months = Object.keys(this.props.costs.months);

    if (!months.length)
      message = (<h4 className="no-data">No data available for this timerange</h4>);
    else {
      let selectedMonthProducts = [];
      let previousMonthProducts = [];
      const parsedMonths = months.map((month) => (moment(month)));

      if (parsedMonths.length === 2) {
        selectedMonthProducts = this.props.costs.months[months[1]].product;
        previousMonthProducts = this.props.costs.months[months[0]].product;
      } else if (parsedMonths[0].isSame(this.props.date, "month")) {
        selectedMonthProducts = this.props.costs.months[months[0]].product;
      }

      if (Object.keys(selectedMonthProducts).length) {
        const currentMonthTotal = this.getMonthTotal(selectedMonthProducts);

        monthCost = (
          <div className="hl-card">
            <ul className="in-col">
              <li>
                <i className="fa fa-credit-card card-icon blue-color"/>
              </li>
              <li>
                <h3 className="no-margin no-padding font-light">
                  {formatPrice(currentMonthTotal)}
                </h3>
              </li>
            </ul>
            <h4 className="card-label p-l-10 m-b-0">
              total spent in {moment(this.props.date).format('MMM Y')}
            </h4>
          </div>
        );

        if (Object.keys(previousMonthProducts).length) {
          const previousTotal = this.getMonthTotal(previousMonthProducts);

          let projectedCurrentMonthTotal = currentMonthTotal;
          if (this.props.currentInterval)
            projectedCurrentMonthTotal = (currentMonthTotal / moment().date()) * parseInt(moment().endOf('month').format("DD"), 10);

          const percentVariation = 0 - (100 - ((projectedCurrentMonthTotal * 100) / previousTotal));

          variation = (
            <div className="hl-card">
              <ul className="in-col">
                <li>
                  <i className="fa fa-area-chart card-icon blue-color"/>
                </li>
                <li>
                  <h3 className={`no-margin no-padding font-light ${percentVariation < 0 ? 'green-color': 'red-color'}`}>
                    {percentVariation > 0 && '+'}{percentVariation.toFixed(2)}%
                  </h3>
                </li>
              </ul>
              <h4 className="card-label p-l-10 m-b-0">
                variation from {moment(this.props.date).subtract(1, 'months').format('MMM Y')}
              </h4>
            </div>
          );
        }
      } else
        message = (<h4 className="no-data">No data available for this timerange</h4>);
    }

    let savingsElement
    if (this.props.unused
      && this.props.unused.ec2
      && this.props.unused.ec2.status) {
      let savings = this.getPotentialsSavings(this.props.unused);

      savingsElement = (
        <div className="hl-card">
          <ul className="in-col">
            <li>
              <i className="fa fa-power-off card-icon blue-color"/>
            </li>
            <li>
              <h3 className={`no-margin no-padding font-light`}>
                {formatPrice(savings.toFixed(2))}
              </h3>
            </li>
          </ul>
          <h4 className="card-label p-l-10 m-b-0">
            potential savings
          </h4>
        </div>
      );
    }

    return (
      <div className="col-md-12">
        <div className="white-box">
          {message}
          {monthCost}
          {variation}
          {savingsElement}
          <div className="clearfix"/>
        </div>
      </div>
    );
  }
}

SummaryComponent.propTypes = {
  costs: PropTypes.shape({
    months : PropTypes.object.isRequired,
  }).isRequired,
  date: PropTypes.object.isRequired,
  currentInterval: PropTypes.bool.isRequired
};

export default SummaryComponent;
