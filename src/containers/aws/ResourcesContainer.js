import React, { Component } from 'react';
import {connect} from 'react-redux';
import Components from '../../components';
import s3square from '../../assets/s3-square.png';
import PropTypes from "prop-types";
import Actions from "../../actions";

const Panel = Components.Misc.Panel;
const IntervalNavigator = Components.Misc.IntervalNavigator;
const VMs = Components.AWS.Resources.VMs;
const Databases = Components.AWS.Resources.Databases;
const ElasticSearch = Components.AWS.Resources.ElasticSearch;

export class ResourcesContainer extends Component {

  render() {
    return(
      <Panel>
        <div className="clearfix">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            Resources
          </h3>
          <div className="inline-block pull-right">
            <IntervalNavigator
              startDate={this.props.dates.startDate}
              endDate={this.props.dates.endDate}
              setDatesFunc={this.props.setDates}
              interval="month"
              hideIntervalSelector={true}
            />
          </div>
        </div>
        <VMs/>
        <Databases/>
        <ElasticSearch/>
      </Panel>
    );
  }
}

ResourcesContainer.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  dates: PropTypes.object,
  setDates: PropTypes.func.isRequired,
  resetDates: PropTypes.func.isRequired,
  clearDates: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  dates: aws.resources.dates,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  setDates: (start, end) => {
    dispatch(Actions.AWS.Resources.setDates(start, end));
  },
  resetDates: () => {
    dispatch(Actions.AWS.Resources.resetDates());
  },
  clearDates: () => {
    dispatch(Actions.AWS.Resources.clearDates());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesContainer);
