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
                return 'badge green-bg'
            case 1:
                return 'badge orange-bg'
            case 2:
                return 'badge red-bg'
            case 3:
                return 'badge red-bg'
            default:
                return 'badge red-bg'
        }
    }
    
    render() {
        const propsValues = this.props.events;

        let events = [];
        if (propsValues) {
            const abnormalsList = [];
            for (var key in propsValues) {
              if (propsValues.hasOwnProperty(key)) {
                  const dataSet = propsValues[key];
                  const abnormals = this.getAbnormals(dataSet);
                  if (abnormals) {
                    for (let i = 0; i < abnormals.length; i++) {
                      const element = abnormals[i];
                      abnormalsList.push({element, key});
                    }
                  }
              }
            }
            abnormalsList.sort((a, b) => {
              if ((a.element.cost - a.element.upper_band) < (b.element.cost - b.element.upper_band)) {
                return 1;
              } else {
                return -1;
              }
            });
            for (let i = 0; i < abnormalsList.length; i++) {
              const element = abnormalsList[i].element;
              const key = abnormalsList[i].key;
              events.push(
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
            }
        }

        let noEventsMessage;
        let table;
        if (!events.length) {
            noEventsMessage= <h4 className="hl-panel-title">
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