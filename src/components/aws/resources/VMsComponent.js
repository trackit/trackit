import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from 'moment';
import ReactTable from 'react-table';
import Popover from '@material-ui/core/Popover';
import {formatPercent, formatBytes, formatPrice} from '../../../common/formatters';
import Misc from '../../misc';

const Tooltip = Misc.Popover;

export class Volumes extends Component {

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
            horizontal: 'center',
          }}
        >
          <div
            className="volumes-list"
            onClick={this.handlePopoverClose}
          >
            {Object.keys(this.props.volumes).filter(key => (key !== "total")).map((volume, index) => (<div key={index} className="volumes-item">{volume} : {formatBytes(this.props.volumes[volume])}</div>))}
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

Volumes.propTypes = {
  volumes: PropTypes.object.isRequired
};

export class Tags extends Component {

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
            className="tags-list"
            onClick={this.handlePopoverClose}
          >
            {Object.keys(this.props.tags).map((tag, index) => (<div key={index} className="tags-item">{tag} : {this.props.tags[tag]}</div>))}
          </div>
        </Popover>
        <div
          ref={node => {
            this.anchor = node;
          }}
          onClick={this.handlePopoverOpen}
        >
          <Tooltip placement="left" icon={<i className="fa fa-tags"/>} tooltip="Click to show tags"/>
        </div>
      </div>
    );
  }

}

Tags.propTypes = {
  tags: PropTypes.object.isRequired
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
      const reportsDates = this.props.data.value.map((account) => (Moment(account.reportDate)));
      const oldestReport = Moment.min(reportsDates);
      const newestReport = Moment.max(reportsDates);
      reportDate = (<Tooltip info tooltip={"Reports created between " + oldestReport.format("ddd D MMM HH:mm") + " and " + newestReport.format("ddd D MMM HH:mm")}/>);
      instances = [].concat.apply([], this.props.data.value.map((account) => (account.instances)));
    }

    const regions = [];
    const types = [];
    const purchasings = [];
    if (instances)
      instances.forEach((instance) => {
        instance.ioRead.total = Object.keys(instance.ioRead).map((volume) => (instance.ioRead[volume])).reduce((a, b) => (a+b));
        instance.ioWrite.total = Object.keys(instance.ioWrite).map((volume) => (instance.ioWrite[volume])).reduce((a, b) => (a+b));
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
            accessor: 'cost',
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
                ) : "N/A")
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
                ) : "N/A")
              }
            ]
          },
          {
            Header: 'IO',
            columns: [
              {
                Header: 'Read',
                accessor: 'ioRead',
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
                accessor: 'ioWrite',
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
                accessor: 'networkIn',
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? formatBytes(row.value) : "N/A")
              },
              {
                Header: 'Out',
                accessor: 'networkOut',
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
        <h3 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-desktop"/>
          &nbsp;
          VMs
          {reportDate}
        </h3>
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
      instances: PropTypes.arrayOf(
        PropTypes.shape({
          id: PropTypes.string.isRequired,
          state: PropTypes.string.isRequired,
          region: PropTypes.string.isRequired,
          cpuAverage: PropTypes.number.isRequired,
          cpuPeak: PropTypes.number.isRequired,
          ioRead: PropTypes.object.isRequired,
          ioWrite: PropTypes.object.isRequired,
          networkIn: PropTypes.number.isRequired,
          networkOut: PropTypes.number.isRequired,
          keyPair: PropTypes.string.isRequired,
          type: PropTypes.string.isRequired,
          tags: PropTypes.object.isRequired
        })
      )
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
