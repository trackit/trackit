import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from 'moment';
import ReactTable from 'react-table';
import Popover from '@material-ui/core/Popover';
import {formatPercent, formatPrice, formatMegaBytes, formatGigaBytes} from '../../../common/formatters';
import Misc from '../../misc';
import Costs from "./misc/Costs";
import Tags from './misc/Tags';

const Tooltip = Misc.Popover;

const getTotalCost = (costs) => {
  let total = 0;
  Object.keys(costs).forEach((key) => total += costs[key]);
  return total;
};

export class ElasticSearchComponent extends Component {

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
    let domains = [];
    if (this.props.data.status && this.props.data.hasOwnProperty("value") && this.props.data.value) {
      domains = this.props.data.value.map((item) => item.domain);
      const reportsDates = this.props.data.value.map((account) => (Moment(account.reportDate)));
      const oldestReport = Moment.min(reportsDates);
      const newestReport = Moment.max(reportsDates);
      reportDate = (<Tooltip info tooltip={"Reports created between " + oldestReport.format("ddd D MMM HH:mm") + " and " + newestReport.format("ddd D MMM HH:mm")}/>);
    }

    const regions = [];
    const types = [];
    if (domains)
      domains.forEach((domain) => {
        if (regions.indexOf(domain.region) === -1)
          regions.push(domain.region);
        if (types.indexOf(domain.type) === -1)
          types.push(domain.type);
      });
    regions.sort();
    types.sort();

    const list = (!loading && !error ? (
      <ReactTable
        data={domains}
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
            accessor: 'domainName',
            minWidth: 150,
            Cell: row => (<strong>{row.value}</strong>)
          },
          {
            Header: 'ID',
            accessor: 'domainId',
            Cell: row => (<strong>{row.value}</strong>)
          },
          {
            Header: 'Type',
            accessor: 'instanceType',
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
            Header: 'Instances',
            accessor: 'instanceCount',
            maxWidth: 100,
            filterable: false
          },
          {
            Header: 'Region',
            accessor: 'region',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {regions.map((region, index) => (<option key={index} value={region}>{region}</option>))}
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
            Header: 'Storage',
            columns: [
              {
                Header: 'Total',
                accessor: 'totalStorageSpace',
                filterable: false,
                Cell: row => formatGigaBytes(row.value)
              },
              {
                Header: 'Unused',
                id: 'freeStorageSpace',
                accessor: d => d.stats.freeSpace,
                filterable: false,
                Cell: row => formatMegaBytes(row.value)
              },
            ]
          },
          {
            Header: 'CPU',
            columns: [
              {
                Header: 'Average',
                id: 'cpuAverage',
                accessor: d => d.stats.cpu.average,
                filterable: false,
                Cell: row => (
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
                )
              },
              {
                Header: 'Peak',
                id: 'cpuPeak',
                accessor: d => d.stats.cpu.peak,
                filterable: false,
                Cell: row => (
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
                )
              }
            ]
          },
          {
            Header: 'Memory Pressure',
            columns: [
              {
                Header: 'Average',
                id: 'jvmMemoryPressureAverage',
                accessor: d => d.stats.JVMMemoryPressure.average,
                filterable: false,
                Cell: row => (
                  <div className="jvm-memory-pressure-stats">
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
                              backgroundColor: row.value > 85 ? '#d6413b'
                                : row.value > 75 ? '#ff9800'
                                  : '#4caf50',
                              borderRadius: '2px'
                            }}
                          />
                        </div>
                      )}
                      tooltip={formatPercent(row.value, 2, false)}
                    />
                  </div>
                )
              },
              {
                Header: 'Peak',
                id: 'jvmMemoryPressurePeak',
                accessor: d => d.stats.JVMMemoryPressure.peak,
                filterable: false,
                Cell: row => (
                  <div className="jvm-memory-pressure-stats">
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
                              backgroundColor: row.value > 85 ? '#d6413b'
                                : row.value > 75 ? '#ff9800'
                                  : '#4caf50',
                              borderRadius: '2px'
                            }}
                          />
                        </div>
                      )}
                      tooltip={formatPercent(row.value, 2, false)}
                    />
                  </div>
                )
              }
            ]
          }
        ]}
        defaultSorted={[{
          id: 'name'
        }]}
        defaultPageSize={10}
        className=" -highlight"
      />
    ) : null);

    return (
      <div className="clearfix resources search-engines">
        <h4 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-search-plus"/>
          &nbsp;
          ElasticSearch
          {reportDate}
        </h4>
        {loading}
        {error}
        {list}
      </div>
    )
  }

}

ElasticSearchComponent.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  data: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.arrayOf(PropTypes.shape({
      account: PropTypes.string.isRequired,
      reportDate: PropTypes.string.isRequired,
      domain: PropTypes.shape({
        domainId: PropTypes.string.isRequired,
        domainName: PropTypes.string.isRequired,
        region: PropTypes.string.isRequired,
        costs: PropTypes.object,
        stats: PropTypes.shape({
          cpu: PropTypes.shape({
            average: PropTypes.number,
            peak: PropTypes.number
          }),
          JVMMemoryPressure: PropTypes.shape({
            average: PropTypes.number,
            peak: PropTypes.number
          }),
          freeSpace: PropTypes.number.isRequired,
        }),
        totalStorageSpace: PropTypes.number.isRequired,
        instanceType: PropTypes.string.isRequired,
        instanceCount: PropTypes.number.isRequired,
        tags: PropTypes.object.isRequired
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
  data: aws.resources.ES
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (date) => {
    dispatch(Actions.AWS.Resources.get.ES(date));
  },
  clear: () => {
    dispatch(Actions.AWS.Resources.clear.ES());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(ElasticSearchComponent);
