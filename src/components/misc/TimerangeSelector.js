import React, {Component} from 'react';
import PropTypes from 'prop-types';

import Moment from 'moment';
import DateRangePicker from 'react-bootstrap-daterangepicker';

class TimerangeSelector extends Component {

  constructor(props) {
    super(props);
    this.handleApply = this.handleApply.bind(this);
    this.handleInterval = this.handleInterval.bind(this);
  }

  handleApply(event, picker) {
    this.props.setDatesFunc(picker.startDate, picker.endDate);
  }

  handleInterval(event) {
    this.props.setIntervalFunc(event.target.value);
  }

  render() {

    const ranges = {
     'Last 7 Days': [Moment().subtract(6, 'days'), Moment()],
     'Last 30 Days': [Moment().subtract(29, 'days'), Moment()],
     'This Month': [Moment().startOf('month'), Moment()],
     'Last Month': [Moment().subtract(1, 'month').startOf('month'), Moment().subtract(1, 'month').endOf('month')],
     'This Year': [Moment().startOf('year'), Moment()],
     'Last Year': [Moment().subtract(1, 'year').startOf('year'), Moment().subtract(1, 'year').endOf('year')]
    };

    return(
      <div>
        <DateRangePicker
          parentEl="body"
          startDate={Moment(this.props.startDate)}
          endDate={Moment(this.props.endDate)}
          ranges={ranges}
          opens="left"
          onApply={this.handleApply}
        >
            <button className="btn btn-default">
              <i className="fa fa-calendar"/>
              &nbsp;
              {this.props.startDate.format('MMM Do Y')}
              &nbsp;
              <i className="fa fa-long-arrow-right"/>
              &nbsp;
              {this.props.endDate.format('MMM Do Y')}
            </button>
        </DateRangePicker>
        {(this.props.interval && this.props.setIntervalFunc) &&
          <select value={this.props.interval} onChange={this.handleInterval}>
            <option value="daily">Daily</option>
            <option value="monthly">Monthly</option>
            <option value="yearly">Yearly</option>
          </select>
        }
      </div>
    );
  }

}

TimerangeSelector.propTypes = {
  startDate: PropTypes.object.isRequired,
  endDate: PropTypes.object.isRequired,
  setDatesFunc: PropTypes.func.isRequired,
  interval: PropTypes.string, //only if interval is needed
  setIntervalFunc: PropTypes.func, //only if interval is needed
};

export default TimerangeSelector;
