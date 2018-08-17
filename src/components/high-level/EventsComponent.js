import React, { Component } from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';
import { Link } from 'react-router-dom';
import { formatPrice } from '../../common/formatters';


class EventsComponent extends Component {
  getAbnormals(data) {
    const res = [];
    for (let i = 0; i < data.length; i++) {
      const element = data[i];
      if (element.abnormal) {
        res.push(element);
      }
    }
    return res.length ? res : null;
  }

  getBadgeClasses(level) {
    switch (level) {
      case 0:
        return 'badge green-bg';
      case 1:
        return 'badge orange-bg';
      case 2:
      case 3:
      default:
        return 'badge red-bg';
    }
  }

  render() {
    const propsValues = this.props.events;

    let events = [];

    if (propsValues) {
      const abnormalsList = [];

      Object.keys(propsValues).forEach((account) => {
        Object.keys(propsValues[account]).forEach((key) => {
          const event = propsValues[account][key];
          const abnormals = event.filter((item) => (item.abnormal));
          abnormals.forEach((element) => {
            abnormalsList.push({element, key, event});
          });
        });
      });

      abnormalsList.sort((a, b) => ((moment(a.element.date).isBefore(b.element.date)) ? 1 : -1));

      events = abnormalsList.map((abnormal) => {
        const element = abnormal.element;
        const key = abnormal.key;
        return (
          <tr key={`${element.date}-${key}`}>
            <td className="badge-cell"><span className={`${this.getBadgeClasses(element.level)} p-l-10 p-r-10`}>{element.pretty_level}</span></td>
            <td><strong>{key.length ? key : 'Unknown service'}</strong></td>
            <td>{moment(element.date).format('MMM, Do')}</td>
            <td>
              {formatPrice(element.cost)}
            </td>
            <td>
              <strong>+{(element.cost - element.upper_band).toFixed(2)}</strong>
            </td>
          </tr>
        );
      })
    }

    let noEventsMessage;
    let table;
    if (!events.length) {
      noEventsMessage= <h4 className="hl-panel-title m-t-20">
        <i className="fa fa-check"></i>
        &nbsp;
        All good. No detected events.
      </h4>;
    } else {
      table = (
        <div className="hl-table-wrapper">
          <table className="hl-table">
            <thead>
            <tr>
              <th></th>
              <th>Service</th>
              <th>Date</th>
              <th>Cost</th>
              <th>Exceed by</th>
            </tr>
            </thead>
            <tbody>
            {events}
            </tbody>
          </table>
        </div>
      );
    }

    return (
      <div className="col-md-6">
        <div className="white-box hl-panel">
          <h4 className="m-t-0 hl-panel-title">
            {moment(this.props.date).format('MMM Y')} Events
          </h4>
          <Link to="/app/events" className="hl-details-link">
            More details
          </Link>
          <hr className="m-b-0"/>
          {table}
          {noEventsMessage}
          <div className="clearfix"></div>
        </div>
      </div>
    );
  }
}

EventsComponent.propTypes = {
  events: PropTypes.object.isRequired,
  date: PropTypes.object.isRequired,
}

export default EventsComponent;