import React, { Component } from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';
import { formatPrice } from '../../common/formatters';

function daysInMonth (month, year) {
    return new Date(year, month, 0).getDate();
}


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

    render() {
        const currentMonthProducts = this.props.costs.month[Object.keys(this.props.costs.month)[0]].product;
        const currentMonthTotal = this.getMonthTotal(currentMonthProducts)
        const previousProducts = this.props.costs.previousMonth[Object.keys(this.props.costs.previousMonth)[0]].product;
        const previousTotal = this.getMonthTotal(previousProducts);

        let projectedCurrentMonthTotal = 0
        // Selected month is current month
        if (moment(Object.keys(this.props.costs.month)[0]).month() === moment().month()) {
            projectedCurrentMonthTotal = (currentMonthTotal / moment().date()) * daysInMonth(moment().month(), moment().year());
        } else {
            projectedCurrentMonthTotal = currentMonthTotal;
        }
        const percentVariation = 0 - (100 - ((projectedCurrentMonthTotal * 100) / previousTotal));
        return (
            <div className="col-md-12">
                <div className="white-box">
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
                    <div className="clearfix"></div>
                </div>
            </div>
        );
    }
}

SummaryComponent.propTypes = {
    costs: PropTypes.shape({
        month : PropTypes.object.isRequired,
        previousMonth: PropTypes.object.isRequired,
    }).isRequired,
    date: PropTypes.object.isRequired,
}

export default SummaryComponent;