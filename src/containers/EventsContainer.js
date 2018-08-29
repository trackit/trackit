import React, { Component } from 'react';
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


  render() {
    const timerange = (this.props.dates ?  (
      <TimerangeSelector
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.props.setDates}
      />
    ) : null);

    return (
      <div>
        {timerange}
      </div>
    );
  }

}

EventsContainer.propTypes = {
  dates: PropTypes.object.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  getData: PropTypes.func.isRequired,
  setDates: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws, events}) => ({
  dates: events.dates,
  accounts: aws.accounts.selection
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

