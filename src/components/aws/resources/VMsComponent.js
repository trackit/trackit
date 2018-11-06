import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from 'moment';
import ReactTable from 'react-table';
import {formatPercent, formatBytes, formatPrice} from '../../../common/formatters';
import Misc from '../../misc';
import Tags from './misc/Tags';
import Volumes from './misc/Volumes';
import Costs from './misc/Costs';

const Tooltip = Misc.Popover;

const getTotalCost = (costs) => {
  let total = 0;
  Object.keys(costs).forEach((key) => total += costs[key]);
  return total;
};

export class VMsComponent extends Component {

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

    const regions = [];
    const types = [];
    const purchasings = [];
    if (instances)
      instances.forEach((instance) => {
        if (instance.stats && instance.stats.volumes && instance.stats.volumes.length) {
          instance.stats.volumes.read.total = Object.keys(instance.stats.volumes.read).map((volume) => (instance.stats.volumes.read[volume])).reduce((a, b) => (a + b));
          instance.stats.volumes.write.total = Object.keys(instance.stats.volumes.write).map((volume) => (instance.stats.volumes.write[volume])).reduce((a, b) => (a + b));
        }
        if (regions.indexOf(instance.region) === -1)
          regions.push(instance.region);
        if (types.indexOf(instance.type) === -1)
          types.push(instance.type);
        if (purchasings.indexOf(instance.purchasing) === -1)
          purchasings.push(instance.purchasing);
      });
    regions.sort();
    types.sort();
    purchasings.sort();

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
            id: 'name',
            accessor: row => (row.tags.hasOwnProperty("Name") ? row.tags.Name : ""),
            minWidth: 150,
            filterMethod: (filter, row) =>
              (row.tags.hasOwnProperty("Name") ? String(row.tags.Name) : "").toLowerCase().includes(filter.value),
            Cell: row => (<strong>{row.value}</strong>)
          },
          {
            Header: 'ID',
            accessor: 'id',
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
            Header: 'Purchasing Option',
            accessor: 'purchasing',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {purchasings.map((purchasing, index) => (<option key={index} value={purchasing}>{purchasing}</option>))}
              </select>
            )
          },
          {
            Header: 'CPU',
            columns: [
              {
                Header: 'Average',
                id: 'cpuAverage',
                accessor: d => d.stats.cpu.average,
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
                id: 'cpuPeak',
                accessor: d => d.stats.cpu.peak,
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
          {
            Header: 'IO',
            columns: [
              {
                Header: 'Read',
                id: 'ioRead',
                accessor: d => d.stats.volumes.read || null,
                minWidth: 120,
                filterable: false,
                sortMethod: (a, b) => (a && b && a.total > b.total ? 1 : -1),
                Cell: row => (row.value ? (row.value.total >= 0 ? (
                  <div className="ioDetails">
                    <span>
                      {formatBytes(row.value.total)}
                    </span>
                    <Volumes volumes={row.value}/>
                  </div>
                ) : "N/A") : null)
              },
              {
                Header: 'Write',
                id: 'ioWrite',
                accessor: d => d.stats.volumes.write || null,
                minWidth: 120,
                filterable: false,
                sortMethod: (a, b) => (a && b && a.total > b.total ? 1 : -1),
                Cell: row => (row.value ? (row.value.total >= 0 ? (
                  <div className="ioDetails">
                    <span>
                      {formatBytes(row.value.total)}
                    </span>
                    <Volumes volumes={row.value}/>
                  </div>
                ) : "N/A") : null)
              }
            ]
          },
          {
            Header: 'Network',
            columns: [
              {
                Header: 'In',
                id: 'networkIn',
                accessor: d => d.stats.network.in,
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? formatBytes(row.value) : "N/A")
              },
              {
                Header: 'Out',
                id: 'networkOut',
                accessor: d => d.stats.network.out,
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? formatBytes(row.value) : "N/A")
              }
            ]
          },
          {
            Header: 'Key Pair',
            accessor: 'keyPair'
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
      <div className="clearfix resources vms">
        <h4 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-desktop"/>
          &nbsp;
          VMs
          {reportDate}
        </h4>
        {loading}
        {error}
        {list}
      </div>
    )
  }

}

VMsComponent.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  data: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.arrayOf(PropTypes.shape({
      account: PropTypes.string.isRequired,
      reportDate: PropTypes.string.isRequired,
      instance: PropTypes.shape({
        id: PropTypes.string.isRequired,
        state: PropTypes.string.isRequired,
        region: PropTypes.string.isRequired,
        keyPair: PropTypes.string.isRequired,
        type: PropTypes.string.isRequired,
        purchasing: PropTypes.string.isRequired,
        tags: PropTypes.object.isRequired,
        costs: PropTypes.object,
        stats: PropTypes.shape({
          cpu: PropTypes.shape({
            average: PropTypes.number,
            peak: PropTypes.number
          }),
          network: PropTypes.shape({
            in: PropTypes.number,
            out: PropTypes.number
          }),
          volumes: PropTypes.shape({
            read: PropTypes.object,
            write: PropTypes.object
          }),
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
  data: aws.resources.EC2
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (date) => {
    dispatch(Actions.AWS.Resources.get.EC2(date));
  },
  clear: () => {
    dispatch(Actions.AWS.Resources.clear.EC2());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(VMsComponent);
