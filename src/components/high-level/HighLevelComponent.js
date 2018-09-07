import React, {Component} from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import moment from "moment/moment";

import Actions from '../../actions';

import IntervalNavigator from '../misc/IntervalNavigator';
import StatusBadges from '../aws/accounts/StatusBadgesComponent';
import Summary from './SummaryComponent';
import TopSpendings from './TopSpendingsComponent';
import History from './HistoryComponent';
import Events from './EventsComponent';

const defaultDates = {
  startDate: moment().startOf('month'),
  endDate: moment().endOf('month')
};


// HighLevelComponent Component
export class HighLevelComponent extends Component {

  componentDidMount() {
    if (this.props.dates) {
      this.props.getData(this.props.dates.startDate, this.props.dates.endDate);
      this.props.getEvents(this.props.dates.startDate, this.props.dates.endDate);
    }
    else
      this.props.setDates(defaultDates.startDate, defaultDates.endDate);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.dates && (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts)) {
      nextProps.getData(nextProps.dates.startDate, nextProps.dates.endDate);
      nextProps.getEvents(nextProps.dates.startDate, nextProps.dates.endDate);
    }
  }

  render() {
    const timerange = (this.props.dates ?  (
      <IntervalNavigator
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.props.setDates}
        interval={'month'}
        hideIntervalSelector={true}
      />
    ) : null);

    let badges;

    if (this.props.costs && this.props.costs.status) {
      badges = (
        <StatusBadges
          values={
            this.props.costs ? (
              this.props.costs.status ? this.props.costs.values : {}
            ) : {}
          }
        />
      );
    }

    let summary;
    if (this.props.costs && this.props.costs.status && this.props.costs.values && this.props.costs.values.month && this.props.costs.values.previousMonth) {
      summary = <Summary
        costs={this.props.costs.values}
        date={this.props.dates.startDate}
      />;
    }

    let topSpendings;
    if (this.props.costs && this.props.costs.status && this.props.costs.values && this.props.costs.values.month && this.props.costs.values.previousMonth) {
      topSpendings = <TopSpendings
        costs={this.props.costs.values}
        date={this.props.dates.startDate}
      />;
    }

    let history;
    if (this.props.costs && this.props.costs.status && this.props.costs.values && this.props.costs.values.history) {
      history = <History
        history={this.props.costs.values.history}
      />;
    }

    let events;
    if (this.props.events && this.props.events.status && this.props.events.values) {
      events = <Events
        events={this.props.events.values}
        date={this.props.dates.startDate}
      />;
    }

  

    return (
        <div>
            <div className="col-md-12">
              <div className="white-box">
                  <div className="clearfix">
                      <h3 className="white-box-title no-padding inline-block">
                          <i className="fa fa-home"></i>
                          &nbsp;
                          Home
                          {badges}
                      </h3>
                      <div className="inline-block pull-right">
                          {timerange}
                      </div>
                  </div>
              </div>
            </div>
            {summary}
            {history}
            {topSpendings}
            {events}
        </div>


    );
  }

}

HighLevelComponent.propTypes = {
  dates: PropTypes.shape({
    startDate: PropTypes.object.isRequired,
    endDate: PropTypes.object.isRequired,
  }),
  setDates: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = ({highlevel, aws}) => ({
  dates: highlevel.dates,
  costs: highlevel.costs,
  events: highlevel.events,
  accounts: aws.accounts.selection,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (begin, end) => {
    dispatch(Actions.Highlevel.getCosts(begin, end))
  },
  getEvents: (begin, end) => {
    dispatch(Actions.Highlevel.getEvents(begin, end))
  },
  setDates: (startDate, endDate) => {
    dispatch(Actions.Highlevel.setDates(startDate, endDate))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(HighLevelComponent);
