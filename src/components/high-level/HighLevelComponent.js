import React, {Component} from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import moment from "moment/moment";
import Spinner from "react-spinkit";

import Actions from '../../actions';

import IntervalNavigator from '../misc/IntervalNavigator';
import Summary from './SummaryComponent';
import TopSpendings from './TopSpendingsComponent';
import TopTags from './TopTagsComponent';
import History from './HistoryComponent';
import Events from './EventsComponent';
import Unused from './TopUnusedComponent';

const defaultDates = {
  startDate: moment().startOf('month'),
  endDate: moment().endOf('month')
};


// HighLevelComponent Component
export class HighLevelComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      currentInterval: true
    };
    this.isCurrentInterval = this.isCurrentInterval.bind(this);
    this.handleKeyChange = this.handleKeyChange.bind(this);
  }

  componentDidMount() {
    if (this.props.dates) {
      this.props.getData(this.props.dates.startDate, this.props.dates.endDate);
      this.props.getEvents(this.props.dates.startDate, this.props.dates.endDate);
      this.props.getUnusedEC2(this.props.dates.startDate);
      this.props.getTagsKeys(this.props.dates.startDate, this.props.dates.endDate);
    }
    else
      this.props.setDates(defaultDates.startDate, defaultDates.endDate);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.dates && (this.props.dates !== nextProps.dates || this.props.accounts !== nextProps.accounts)) {
      nextProps.getData(nextProps.dates.startDate, nextProps.dates.endDate);
      nextProps.getEvents(nextProps.dates.startDate, nextProps.dates.endDate);
      nextProps.getUnusedEC2(nextProps.dates.startDate);
      nextProps.getTagsKeys(nextProps.dates.startDate, nextProps.dates.endDate);
    }
  }

  isCurrentInterval(currentInterval) {
    if (this.state.currentInterval !== currentInterval)
      this.setState({currentInterval})
  }

  handleKeyChange(key) {
      this.props.setTagsKeySelected(key);
      this.props.getTagsValues(this.props.dates.startDate, this.props.dates.endDate, key);
  }

  render() {
    const timerange = (this.props.dates ?  (
      <IntervalNavigator
        startDate={this.props.dates.startDate}
        endDate={this.props.dates.endDate}
        setDatesFunc={this.props.setDates}
        interval={'month'}
        hideIntervalSelector={true}
        isCurrentInterval={this.isCurrentInterval}
      />
    ) : null);

    let costLoader;
    let costError;
    let summary;
    let topSpendings;
    let history;

    if (this.props.costs) {
      if (!this.props.costs.status)
        costLoader = <Spinner className="spinner" name='circle'/>;
      else if (this.props.costs.hasOwnProperty("error"))
        costError = <div className="alert alert-warning" role="alert">Error while getting data ({this.props.costs.error.message})</div>;
      else if (this.props.costs.values) {
        if (this.props.costs.values.months) {
          summary = <Summary
            costs={this.props.costs.values}
            date={this.props.dates.startDate}
            currentInterval={this.state.currentInterval}
            unused={this.props.unused}
          />;
          topSpendings = <TopSpendings
            costs={this.props.costs.values}
            date={this.props.dates.startDate}
            currentInterval={this.state.currentInterval}
          />;
        }
        if (this.props.costs.values.history)
          history = <History
            history={this.props.costs.values.history}
          />;
      }
    }

    let tagsLoader;
    let tagsError;
    let tags;
    if (this.props.tags && this.props.tags.keys) {
      if (!this.props.tags.keys.status)
        tagsLoader = <Spinner className="spinner" name='circle'/>;
      else if (this.props.tags.keys.hasOwnProperty("error"))
        tagsError = <div className="alert alert-warning" role="alert">Error while getting data
          ({this.props.tags.keys.error.message})</div>;
      else if (this.props.tags.keys.values)
        tags = <TopTags
          date={this.props.dates.startDate}
          keys={this.props.tags.keys.values}
          selected={this.props.tags.selected}
          currentInterval={this.state.currentInterval}
          setSelected={this.handleKeyChange}
          costs={this.props.tags.costs}
        />;
    }

    let eventsLoader;
    let eventsError;
    let events;
    if (this.props.events) {
      if (!this.props.events.status)
        eventsLoader = <Spinner className="spinner" name='circle'/>;
      else if (this.props.events.hasOwnProperty("error"))
        eventsError = <div className="alert alert-warning" role="alert">Error while getting data
          ({this.props.events.error.message})</div>;
      else if (this.props.events.values)
        events = <Events
          events={this.props.events.values}
          date={this.props.dates.startDate}
        />;
    }

   let unusedLoader;
   let unusedError;
   let unused;
   if (this.props.unused) {
     if (!this.props.unused.ec2.status)
       unusedLoader = <Spinner className="spinner" name='circle'/>;
     else if (this.props.unused.ec2.hasOwnProperty("error"))
       unusedError = <div className="alert alert-warning" role="alert">Error while getting data
         ({this.props.unused.ec2.error.message})</div>;
     else if (this.props.unused.ec2.values)
       unused = <Unused
         unused={this.props.unused}
         date={this.props.dates.startDate}
       />;
   }

    const status = ((eventsLoader || eventsError || costLoader || costError || tagsLoader || tagsError || unusedLoader || unusedError) ? (
      <div className="col-md-12">
        <div className="white-box">
          {eventsLoader}
          {eventsError}
          {costLoader}
          {costError}
          {tagsLoader}
          {tagsError}
          {unusedLoader}
          {unusedError}
        </div>
      </div>
    ) : null);

    return (
      <div>
        <div className="row">
          <div className="col-md-12">
            <div className="white-box">
              <div className="clearfix">
                <h3 className="white-box-title no-padding inline-block">
                  <i className="fa fa-home"></i>
                  &nbsp;
                  Home
                </h3>
                <div className="inline-block pull-right">
                  {timerange}
                </div>
              </div>
            </div>
          </div>
        </div>
        <div className="row">
          {status}
          {summary}
          {history}
        </div>
        <div className="row">
          {topSpendings}
          {tags}
        </div>
        <div className="row">
          {events}
          {unused}
        </div>
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
  tags: highlevel.tags,
  unused: highlevel.unused,
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
  getUnusedEC2: (date) => {
    dispatch(Actions.Highlevel.getUnusedEC2(date))
  },
  getTagsKeys: (begin, end) => {
    dispatch(Actions.Highlevel.getTagsKeys(begin, end))
  },
  getTagsValues: (begin, end, key) => {
    dispatch(Actions.Highlevel.getTagsValues(begin, end, key))
  },
  setTagsKeySelected: (key) => {
    dispatch(Actions.Highlevel.selectTagsKey(key))
  },
  setDates: (startDate, endDate) => {
    dispatch(Actions.Highlevel.setDates(startDate, endDate))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(HighLevelComponent);
