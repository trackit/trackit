import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from "moment";
import Misc from '../../misc';
import Tags from './misc/Tags';
import ReactTable from "react-table";
import {formatGigaBytes, formatPrice, formatBytes, formatPercent} from "../../../common/formatters";
import UnusedStorage from './misc/UnusedStorage';
import Costs from "./misc/Costs";

const Tooltip = Misc.Popover;
const Collapsible = Misc.Collapsible;

const getTotalCost = (costs) => {
  let total = 0;
  Object.keys(costs).forEach((key) => total += costs[key]);
  return total;
};

export class RdsComponent extends Component {

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
      instances = this.props.data.value.map((item) => item.instance);
      const reportsDates = this.props.data.value.map((item) => (Moment(item.reportDate)));
      const oldestReport = Moment.min(reportsDates);
      const newestReport = Moment.max(reportsDates);
      reportDate = (<Tooltip info tooltip={"Reports created between " + oldestReport.format("ddd D MMM HH:mm") + " and " + newestReport.format("ddd D MMM HH:mm")}/>);
    }

    const availabilityZones = [];
    const types = [];
    const engines = [];
    if (instances)
      instances.forEach((instance) => {
        if (availabilityZones.indexOf(instance.availabilityZone) === -1)
          availabilityZones.push(instance.availabilityZone);
        if (types.indexOf(instance.type) === -1)
          types.push(instance.type);
        if (engines.indexOf(instance.engine) === -1)
          engines.push(instance.engine);
      });
    availabilityZones.sort();
    types.sort();
    engines.sort();

    const list = (!loading && !error ? (
      <ReactTable
        data={instances}
        noDataText="No instances available"
        filterable
        defaultFilterMethod={(filter, row) => String(row[filter.id]).toLowerCase().includes(filter.value)}
        columns={[
          {
            Header: 'Tags',
            accessor: 'tags',
            maxWidth: 50,
            filterable: false,
            Cell: row => ((row.value && Object.keys(row.value).length) ?
              (<Tags tags={row.value}/>) :
              (<Tooltip placement="left" icon={<i className="fa fa-tag disabled"/>} tooltip="No tags"/>))
          },
          {
            Header: 'Name',
            accessor: 'id',
            minWidth: 150,
            Cell: row => (<strong>{row.value}</strong>)
          },
          {
            Header: 'Type',
            accessor: 'type',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {types.map((type, index) => (<option key={index} value={type}>{type}</option>))}
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
            accessor: 'costs',
            filterable: false,
            sortMethod: (a, b) => (a && b && getTotalCost(a) > getTotalCost(b) ? 1 : -1),
            Cell: row => (row.value && Object.keys(row.value).length !== 0 ? (
                <div className="unusedStorageDetails">
                    <span>
                      {formatPrice(getTotalCost(row.value))}
                    </span>
                  <Costs costs={row.value}/>
                </div>
              ) : (
                <span>
                  N/A
                  <Tooltip tooltip='Cost data are unavailable for this timerange. Please check again later.' info triggerStyle={{ fontSize: '0.9em', color: 'inherit' }} />
                </span>
              )
            )
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
                id: 'unusedStorage',
                accessor: d => d.stats.freeSpace || null,
                filterable: false,
                Cell: row => (row.value && row.value.average !== undefined && row.value.average >= 0 ? (
                  <div className="unusedStorageDetails">
                    <span>
                      {formatBytes(row.value.average)}
                    </span>
                    <UnusedStorage data={row.value}/>
                  </div>
                ) : "N/A")
              },
            ]
          },
          {
            Header: 'CPU',
            columns: [
              {
                Header: 'Average',
                id: 'averageCPU',
                accessor: d => d.stats.cpu.average || null,
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
                ) : "N/A")
              },
              {
                Header: 'Peak',
                id: 'peakCPU',
                accessor: d => d.stats.cpu.peak || null,
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
                ) : "N/A")
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

    const header = (
      <div>
        <h4 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-database"/>
          &nbsp;
          RDS
          {reportDate}
        </h4>
        {loading}
        {error}
      </div>
    );

    return (
      <Collapsible
        className="clearfix resources dbs"
        header={header}
      >
        {list}
      </Collapsible>
    );
  }

}

RdsComponent.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  data: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.arrayOf(PropTypes.shape({
      account: PropTypes.string.isRequired,
      reportDate: PropTypes.string.isRequired,
      instance: PropTypes.shape({
        id: PropTypes.string.isRequired,
        type: PropTypes.string.isRequired,
        availabilityZone: PropTypes.string.isRequired,
        engine: PropTypes.string.isRequired,
        multiAZ: PropTypes.bool.isRequired,
        allocatedStorage: PropTypes.number.isRequired,
        tags: PropTypes.object,
        costs: PropTypes.object,
        stats: PropTypes.shape({
          cpu: PropTypes.shape({
            average: PropTypes.number,
            peak: PropTypes.number
          }),
          freeSpace: PropTypes.shape({
            minimum: PropTypes.number,
            maximum: PropTypes.number,
            average: PropTypes.number
          })
        })
      })
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

export default connect(mapStateToProps, mapDispatchToProps)(RdsComponent);
