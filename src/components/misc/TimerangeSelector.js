import React, {Component} from 'react';
import PropTypes from 'prop-types';

import moment from 'moment';
import DateRangePicker from 'react-bootstrap-daterangepicker';

class TimerangeSelector extends Component {

  constructor() {
    super();
    this.handleApply = this.handleApply.bind(this);
  }

  handleApply(event, picker) {
    console.log(event);
    console.log(picker);
    this.props.setDatesFunc(picker.startDate, picker.endDate);
  }

  render() {

    const ranges = {
     'Last 7 Days': [moment().subtract(6, 'days'), moment()],
     'Last 30 Days': [moment().subtract(29, 'days'), moment()],
     'This Month': [moment().startOf('month'), moment()],
     'Last Month': [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')],
     'This Year': [moment().startOf('year'), moment()],
     'Last Year': [moment().subtract(1, 'year').startOf('year'), moment().subtract(1, 'year').endOf('year')]
    };

    return(
      <DateRangePicker
        parentEl={"body"}
        startDate={moment(this.props.startDate)}
        endDate={moment(this.props.endDate)}
        ranges={ranges}
        opens={'left'}
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
    );
  }

};

TimerangeSelector.propTypes = {
  startDate: PropTypes.object.isRequired,
  endDate: PropTypes.object.isRequired,
  setDatesFunc: PropTypes.func.isRequired,
};

export default TimerangeSelector;
