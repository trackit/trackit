import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from "moment";
import Misc from '../../misc';
import ReactTable from "react-table";
import {formatGigaBytes, formatPrice, formatBytes, formatPercent} from "../../../common/formatters";
import Popover from "@material-ui/core/Popover/Popover";

const Tooltip = Misc.Popover;

export class UnusedStorage extends Component {

  constructor(props) {
    super(props);
    this.state = {
      showPopOver: false
    };
    this.handlePopoverOpen = this.handlePopoverOpen.bind(this);
    this.handlePopoverClose = this.handlePopoverClose.bind(this);
  }

  handlePopoverOpen = (e) => {
    e.preventDefault();
    this.setState({ showPopOver: true });
  };

  handlePopoverClose = (e) => {
    e.preventDefault();
    this.setState({ showPopOver: false });
  };

  render() {
    return (
      <div>
        <Popover
          open={this.state.showPopOver}
          anchorEl={this.anchor}
          onClose={this.handlePopoverClose}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'center',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
          }}
        >
          <div
            className="unusedStorage-list"
            onClick={this.handlePopoverClose}
          >
            {Object.keys(this.props.data).map((item, index) => (<div key={index} className="unusedStorage-item">{item} : {this.props.data[item] >= 0 ? formatBytes(this.props.data[item]) : "No data available"}</div>))}
          </div>
        </Popover>
        <div
          ref={node => {
            this.anchor = node;
          }}
          onClick={this.handlePopoverOpen}
        >
          <Tooltip placement="right" info tooltip="Click to see more details"/>
        </div>
      </div>
    );
  }

}

UnusedStorage.propTypes = {
  data: PropTypes.object.isRequired
};

export class DatabasesComponent extends Component {

  componentWillMount() {
    this.props.getData(this.props.dates.startDate);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.accounts !== this.props.accounts || nextProps.dates !== this.props.dates)
      nextProps.getData(nextProps.dates.startDate);
  }

  render() {
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
            Cell: row => {
              if (row.value === 0 && Object.keys(row.original.costDetail).length === 0) {
                return <span>
                  N/A
                  <Tooltip tooltip='Cost data are unavailable for this timerange. Please check again later.' info triggerStyle={{ fontSize: '0.9em', color: 'inherit' }} />
                </span>;
              } else {
                return (formatPrice(row.value));
              }
            }
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
            columns: [
              {
                Header: 'Total',
                accessor: 'allocatedStorage',
                filterable: false,
                Cell: row => formatGigaBytes(row.value)
              },
              {
                Header: 'Unused',
                accessor: 'freeSpaceAverage',
                filterable: false,
                Cell: row => (row.value & row.value >= 0 ? (
                  <div className="unusedStorageDetails">
                    <span>
                      {formatBytes(row.value)}
                    </span>
                    <UnusedStorage data={{
                      Average: row.value,
                      Minimum: row.original.freeSpaceMinimum,
                      Maximum: row.original.freeSpaceMaximum
                    }}/>
                  </div>
                ) : "No data available")
              },
            ]
          },
          {
            Header: 'CPU',
            columns: [
              {
                Header: 'Average',
                accessor: 'cpuAverage',
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? (
                  <div className="cpu-stats">
                    <Tooltip
                      placement="left"
                      icon={(
                        <div
                          style={{
                            height: '100%',
                            backgroundColor: '#dddddd',
                            borderRadius: '2px',
                            flex: 1
                          }}
                        >
                          <div
                            style={{
                              width: `${row.value}%`,
                              height: '100%',
                              backgroundColor: row.value > 60 ? '#d6413b'
                                : row.value > 30 ? '#ff9800'
                                  : '#4caf50',
                              borderRadius: '2px'
                            }}
                          />
                        </div>
                      )}
                      tooltip={formatPercent(row.value, 2, false)}
                    />
                  </div>
                ) : "No data available")
              },
              {
                Header: 'Peak',
                accessor: 'cpuPeak',
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? (
                  <div className="cpu-stats">
                    <Tooltip
                      placement="right"
                      icon={(
                        <div
                          style={{
                            height: '100%',
                            backgroundColor: '#dddddd',
                            borderRadius: '2px',
                            flex: 1
                          }}
                        >
                          <div
                            style={{
                              width: `${row.value}%`,
                              height: '100%',
                              backgroundColor: row.value > 80 ? '#d6413b'
                                : row.value > 60 ? '#ff9800'
                                  : '#4caf50',
                              borderRadius: '2px'
                            }}
                          />
                        </div>
                      )}
                      tooltip={formatPercent(row.value, 2, false)}
                    />
                  </div>
                ) : "No data available")
              }
            ]
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
