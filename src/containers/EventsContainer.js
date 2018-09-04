import React, { Component } from 'react';
import moment from 'moment';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import Components from '../components';
import Actions from '../actions';

const TimerangeSelector = Components.Misc.TimerangeSelector;

// EventsContainer Component
class EventsContainer extends Component {
  componentDidMount() {
    if (this.props.dates) {
      const dates = this.props.dates;
      this.props.getData(dates.startDate, dates.endDate);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.dates && (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts))
      nextProps.getData(nextProps.dates.startDate, nextProps.dates.endDate);
  }

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


  render() {
    const propsValues = this.props.values;

    const timerange = (this.props.dates ?  (
      <TimerangeSelector
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.props.setDates}
      />
    ) : null);

    let events = [];
    if (propsValues && propsValues.status && propsValues.values) {
      const abnormalsList = [];
      for (var key in propsValues.values) {
        if (propsValues.values.hasOwnProperty(key)) {
            const dataSet = propsValues.values[key];
            const abnormals = this.getAbnormals(dataSet);
            if (abnormals) {
              for (let i = 0; i < abnormals.length; i++) {
                const element = abnormals[i];
                abnormalsList.push({element, key, dataSet});
              }
            }
        }
      }
      abnormalsList.sort((a, b) => {
        if (moment(a.element.date).isBefore(b.element.date)) {
          return 1;
        } else {
          return -1;
        }
      });
      for (let i = 0; i < abnormalsList.length; i++) {
        const element = abnormalsList[i].element;
        const key = abnormalsList[i].key;
        const dataSet = abnormalsList[i].dataSet;
        events.push(
        <div key={`${element.date}-${key}`}>
          <Components.Events.EventPanel
            dataSet={dataSet}
            abnormalElement={element}
            service={key}
          />
        </div>
        );
        
      }
    }


    return (
      <div>
        <div className="row">
          <div className="col-md-12">
            <div className="white-box">
              <h3 className="white-box-title no-padding inline-block">
                <i className="fa fa-exclamation-triangle"></i>
                &nbsp;
                Events
              </h3>
              <div className="inline-block pull-right">
                {timerange}
              </div>
            </div>
          </div>
        </div>
        {events}
      </div>
    );
  }

}

EventsContainer.propTypes = {
  dates: PropTypes.object.isRequired,
  values: PropTypes.object.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  getData: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws, events}) => ({
  dates: events.dates,
  accounts: aws.accounts.selection,
  values: events.values,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (begin, end) => {
    dispatch(Actions.Events.getData(begin, end));
  },
  setDates: (startDate, endDate) => {
    dispatch(Actions.Events.setDates(startDate, endDate));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(EventsContainer);

