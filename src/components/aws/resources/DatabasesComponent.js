import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from "moment";
import Misc from '../../misc';
import ReactTable from "react-table";
import {formatGigaBytes, formatPrice} from "../../../common/formatters";

const Tooltip = Misc.Popover;

export class DatabasesComponent extends Component {

  componentWillMount() {
    this.props.getData(this.props.dates.startDate);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.accounts !== this.props.accounts || nextProps.dates !== this.props.dates)
      nextProps.getData(nextProps.dates.startDate);
  }

  render() {
    if (this.props.dates.startDate.isBefore(Moment().startOf('months')))
      return (
        <div className="clearfix resources dbs">
          <h3 className="white-box-title no-padding inline-block">
            <i className="menu-icon fa fa-database"/>
            &nbsp;
            Databases
          </h3>
          <div className="alert alert-warning" role="alert">Report not available for this month</div>
        </div>
      );

    const loading = (!this.props.data.status ? (<Spinner className="spinner" name='circle'/>) : null);
    const error = (this.props.data.error ? (<div className="alert alert-warning" role="alert">Error while getting data ({this.props.data.error.message})</div>) : null);

    let reportDate = null;
    let instances = [];
    if (this.props.data.status && this.props.data.hasOwnProperty("value") && this.props.data.value) {
      const reportsDates = this.props.data.value.map((account) => (Moment(account.reportDate)));
      const oldestReport = Moment.min(reportsDates);
      const newestReport = Moment.max(reportsDates);
      reportDate = (<Tooltip info tooltip={"Reports created between " + oldestReport.format("ddd D MMM HH:mm") + " and " + newestReport.format("ddd D MMM HH:mm")}/>);
      instances = [].concat.apply([], this.props.data.value.map((account) => (account.instances)));
    }

    const availabilityZones = [];
    const dbInstanceClasses = [];
    const engines = [];
    if (instances)
      instances.forEach((instance) => {
        if (availabilityZones.indexOf(instance.availabilityZone) === -1)
          availabilityZones.push(instance.availabilityZone);
        if (dbInstanceClasses.indexOf(instance.dbInstanceClass) === -1)
          dbInstanceClasses.push(instance.dbInstanceClass);
        if (engines.indexOf(instance.engine) === -1)
          engines.push(instance.engine);
      });
    availabilityZones.sort();
    dbInstanceClasses.sort();
    engines.sort();

    const list = (!loading && !error ? (
      <ReactTable
        data={instances}
        noDataText="No instances available"
        filterable
        defaultFilterMethod={(filter, row) => String(row[filter.id]).toLowerCase().includes(filter.value)}
        columns={[
          {
            Header: 'Name',
            accessor: 'dbInstanceIdentifier',
            minWidth: 150,
            Cell: row => (<strong>{row.value}</strong>)
          },
          {
            Header: 'Type',
            accessor: 'dbInstanceClass',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {dbInstanceClasses.map((type, index) => (<option key={index} value={type}>{type}</option>))}
              </select>
            )
          },
          {
            Header: 'Region',
            accessor: 'availabilityZone',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {availabilityZones.map((region, index) => (<option key={index} value={region}>{region}</option>))}
              </select>
            )
          },
          {
            Header: 'Cost',
            id: 'cost',
            accessor: d => d.cost || 0,
            filterable: false,
            Cell: row => (formatPrice(row.value))
          },
          {
            Header: 'Engine',
            accessor: 'engine',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {engines.map((region, index) => (<option key={index} value={region}>{region}</option>))}
              </select>
            )
          },
          {
            Header: 'Multi-AZ',
            accessor: 'multiAZ',
            maxWidth: 100,
            Cell: row => (<i className={"fa " + (row.value === "true" ? "fa-check-circle" : "fa-times-circle")}/>),
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === String(row[filter.id]))),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                <option value="true">Yes</option>
                <option value="false">No</option>
              </select>
            )
          },
          {
            Header: 'Storage',
            accessor: 'allocatedStorage',
            filterable: false,
            Cell: row => formatGigaBytes(row.value)
          },
        ]}
        defaultSorted={[{
          id: 'name'
        }]}
        defaultPageSize={10}
        className=" -highlight"
      />
    ) : null);

    return (
      <div className="clearfix resources dbs">
        <h3 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-database"/>
          &nbsp;
          Databases
          {reportDate}
        </h3>
        {loading}
        {error}
        {list}
      </div>
    )
  }

}

DatabasesComponent.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  data: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.arrayOf(PropTypes.shape({
      account: PropTypes.string.isRequired,
      reportDate: PropTypes.string.isRequired,
      instances: PropTypes.arrayOf(PropTypes.shape({
        dbInstanceIdentifier: PropTypes.string.isRequired,
        dbInstanceClass: PropTypes.string.isRequired,
        availabilityZone: PropTypes.string.isRequired,
        engine: PropTypes.string.isRequired,
        multiAZ: PropTypes.bool.isRequired,
        allocatedStorage: PropTypes.number.isRequired
      }))
    }))
  }),
  getData: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
  dates: PropTypes.object,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.selection,
  dates: aws.resources.dates,
  data: aws.resources.RDS
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (date) => {
    dispatch(Actions.AWS.Resources.get.RDS(date));
  },
  clear: () => {
    dispatch(Actions.AWS.Resources.clear.RDS());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(DatabasesComponent);
