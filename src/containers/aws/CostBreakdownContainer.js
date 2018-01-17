import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions';
import s3square from '../../assets/s3-square.png';

const TimerangeSelector = Components.Misc.TimerangeSelector;
const Selector = Components.Misc.Selector;
const Panel = Components.Misc.Panel;
const Chart = Components.AWS.CostBreakdown.Chart;

export class CostBreakdownContainer extends Component {

  componentWillMount() {
    this.props.getCosts(this.props.costsDates.startDate, this.props.costsDates.endDate, [this.props.costsFilter, this.props.costsInterval]);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.costsDates !== nextProps.costsDates ||
      this.props.costsInterval !== nextProps.costsInterval ||
      this.props.costsFilter !== nextProps.costsFilter ||
      this.props.accounts !== nextProps.accounts)
      nextProps.getCosts(nextProps.costsDates.startDate, nextProps.costsDates.endDate, [nextProps.costsFilter, nextProps.costsInterval]);
  }

  render() {

    const filters = {
      all: "Total",
      account: "Account",
      product: "Product",
      region: "Region"
    };

    return(
      <Panel>

        <div className="clearfix">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            Cost Breakdown
          </h3>

          <div className="inline-block pull-right">
            <div className="inline-block">
              <Selector
                values={filters}
                selected={this.props.costsFilter}
                selectValue={this.props.setCostsFilter}
              />
            </div>
            <div className="inline-block">
              <TimerangeSelector
                startDate={this.props.costsDates.startDate}
                endDate={this.props.costsDates.endDate}
                setDatesFunc={this.props.setCostsDates}
                interval={this.props.costsInterval}
                setIntervalFunc={this.props.setCostsInterval}
              />
            </div>
          </div>

        </div>

        <div>
          <Chart values={this.props.costsValues} interval={this.props.costsInterval} filter={this.props.costsFilter}/>
        </div>

      </Panel>
    );
  }
}

CostBreakdownContainer.propTypes = {
  costsValues: PropTypes.object,
  costsDates: PropTypes.shape({
    startDate: PropTypes.object,
    endDate: PropTypes.object,
  }),
  costsInterval: PropTypes.string.isRequired,
  costsFilter: PropTypes.string.isRequired,
  getCosts: PropTypes.func.isRequired,
  setCostsDates: PropTypes.func.isRequired,
  setCostsInterval: PropTypes.func.isRequired,
  setCostsFilter: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  costsValues: aws.costs.values,
  costsDates: aws.costs.dates,
  costsInterval: aws.costs.interval,
  costsFilter: aws.costs.filter,
  accounts: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getCosts: (begin, end, filters) => {
    dispatch(Actions.AWS.Costs.getCosts(begin, end, filters));
  },
  setCostsDates: (startDate, endDate) => {
    dispatch(Actions.AWS.Costs.setCostsDates(startDate, endDate))
  },
  setCostsInterval: (interval) => {
    dispatch(Actions.AWS.Costs.setCostsInterval(interval));
  },
  setCostsFilter: (filter) => {
    dispatch(Actions.AWS.Costs.setCostsFilter(filter));
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(CostBreakdownContainer);
