import React, {Component} from 'react';
import PropTypes from 'prop-types';
import Moment from 'moment';
import IntervalSelector from './IntervalSelector';

class IntervalNavigator extends Component {

  constructor(props) {
    super(props);
    this.previousDate = this.previousDate.bind(this);
    this.nextDate = this.nextDate.bind(this);
    this.updateInterval = this.updateInterval.bind(this);
  }

  previousDate = (e) => {
    e.preventDefault();
    let start;
    let end;
    switch (this.props.interval) {
      case "year":
        start = this.props.startDate.subtract(1, 'years');
        end = this.props.endDate.subtract(1, 'years');
        break;
      case "month":
        start = this.props.startDate.subtract(1, 'months').startOf('months');
        end = this.props.endDate.subtract(1, 'months').endOf('months');
        break;
      case "week":
        start = this.props.startDate.subtract(1, 'weeks').startOf('isoWeek');
        end = this.props.endDate.subtract(1, 'weeks').endOf('isoWeek');
        break;
      case "day":
      default:
        start = this.props.startDate.subtract(1, 'days');
        end = start;
    }
    this.props.setDatesFunc(start, end);
  };

  nextDate = (e) => {
    e.preventDefault();
    let start;
    let end;
    switch (this.props.interval) {
      case "year":
        start = this.props.startDate.add(1, 'years');
        end = this.props.endDate.add(1, 'years');
        break;
      case "month":
        start = this.props.startDate.add(1, 'months');
        end = this.props.endDate.add(1, 'months').endOf('months');
        break;
      case "week":
        start = this.props.startDate.add(1, 'weeks').startOf('isoWeek');
        end = this.props.endDate.add(1, 'weeks').endOf('isoWeek');
        break;
      case "day":
      default:
        start = this.props.startDate.add(1, 'days');
        end = start;
    }
    this.props.setDatesFunc(start, end);
  };

  getDate() {
    switch (this.props.interval) {
      case "year":
        return this.props.startDate.format('Y');
      case "month":
        return this.props.startDate.format('MMM Y');
      case "week":
        return (
          <div className="inline-block">
            {this.props.startDate.format('MMM Do Y')}
            &nbsp;
            <i className="fa fa-long-arrow-right"/>
            &nbsp;
            {this.props.endDate.format('MMM Do Y')}
          </div>
        );
      case "day":
      default:
        return this.props.startDate.format('MMM Do Y');
    }
  }

  updateInterval(interval) {
    this.props.setIntervalFunc(interval);
    let start;
    let end;
    switch (interval) {
      case "year":
        start = Moment().startOf('year');
        end = Moment().endOf('year');
        break;
      case "month":
        start = Moment().subtract(1, 'month').startOf('month');
        end = Moment().subtract(1, 'month').endOf('month');
        break;
      case "week":
        start = Moment().subtract(1, 'month').endOf('month').startOf('isoWeek');
        end = Moment().subtract(1, 'month').endOf('month').endOf('isoWeek');
        break;
      case "day":
      default:
        start = Moment().subtract(1, 'month').endOf('month');
        end = start;
    }
    return this.props.setDatesFunc(start, end);
  }

  isCurrentInterval() {
    let now;
    switch (this.props.interval) {
      case "year":
        now = Moment().endOf('year');
        break;
      case "month":
        now = Moment().endOf('month');
        break;
      case "week":
        now = Moment().endOf('isoWeek');
        break;
      case "day":
      default:
        now = Moment();
    }
    return this.props.endDate.isSameOrAfter(now);
  }

  render() {
    const currentInterval = this.isCurrentInterval();
    if (this.props.isCurrentInterval)
      this.props.isCurrentInterval(currentInterval);
    return(
      <div className="inline-block">
        <div className="inline-block btn-group">
          <button className="btn btn-default" onClick={this.previousDate}>
            <i className="fa fa-caret-left"/>
          </button>
          <div className="btn btn-default no-click">
            <i className="fa fa-calendar"/>
            &nbsp;
            {this.getDate()}
          </div>
          <button className="btn btn-default" disabled={(this.props.lockFuture ? currentInterval : false)} onClick={this.nextDate}>
            <i className="fa fa-caret-right"/>
          </button>
        </div>
        {!this.props.hideIntervalSelector && <IntervalSelector interval={this.props.interval} setInterval={this.updateInterval} availableIntervals={this.props.availableIntervals}/>}
      </div>
    );
  }

}

IntervalNavigator.defaultProps = {
  hideIntervalSelector: false,
  lockFuture: true,
  isCurrentInterval: null
};

IntervalNavigator.propTypes = {
  startDate: PropTypes.object.isRequired,
  endDate: PropTypes.object.isRequired,
  setDatesFunc: PropTypes.func.isRequired,
  interval: PropTypes.string,
  availableIntervals: PropTypes.arrayOf(PropTypes.string),
  setIntervalFunc: PropTypes.func,
  hideIntervalSelector: PropTypes.bool,
  lockFuture: PropTypes.bool,
  isCurrentInterval: PropTypes.func
};

export default IntervalNavigator;
