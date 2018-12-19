import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from 'moment';
import ReactTable from 'react-table';
import {formatMillisecondsDuration, formatBytes, formatMegaBytes} from '../../../common/formatters';
import Misc from '../../misc';
import Tags from './misc/Tags';
const Tooltip = Misc.Popover;

const formatInvocationPercentage = (invocations) => {
  if (invocations && invocations.total >= 0 && invocations.failed >= 0) {
    const success = (100 - (invocations.failed * 100 / invocations.total)).toFixed(2);
    const style = {
      color: (success > 90 ? '#4caf50' : (success > 75 ? '#ff9800' : '#d6413b'))
    };
    return (
      <div className="success-percentage">
        <span className="success-percentage-value" style={style}>{success} %</span>
        <Tooltip placement="right" info tooltip={`Failed : ${invocations.failed}`}/>
      </div>
    );
  } else
    return "N/A";
};

export class LambdaComponent extends Component {

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
    let lambdas = [];
    if (this.props.data.status && this.props.data.hasOwnProperty("value") && this.props.data.value) {
      lambdas = this.props.data.value.map((item) => item.function);
      const reportsDates = this.props.data.value.map((item) => (Moment(item.reportDate)));
      const oldestReport = Moment.min(reportsDates);
      const newestReport = Moment.max(reportsDates);
      reportDate = (<Tooltip info tooltip={"Reports created between " + oldestReport.format("ddd D MMM HH:mm") + " and " + newestReport.format("ddd D MMM HH:mm")}/>);
    }

    const runtimes = [];
    if (lambdas)
      lambdas.forEach((lambda) => {
        if (runtimes.indexOf(lambda.runtime) === -1)
          runtimes.push(lambda.runtime);
      });
    runtimes.sort();

    const list = (!loading && !error ? (
      <ReactTable
        data={lambdas}
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
            accessor: 'name',
            minWidth: 250,
            Cell: row => (
              <div className="function-name">
                <span><strong>{row.value}</strong></span>
                {row.original.description.length ? (<Tooltip placement="right" info tooltip={row.original.description}/>) : null}
              </div>)
          },
          {
            Header: 'Version',
            accessor: 'version',
          },
          {
            Header: 'Runtime',
            accessor: 'runtime',
            filterMethod: (filter, row) => (filter.value === "all" ? true : (filter.value === row[filter.id])),
            Filter: ({ filter, onChange }) => (
              <select
                onChange={event => onChange(event.target.value)}
                style={{ width: "100%" }}
                value={filter ? filter.value : "all"}
              >
                <option value="all">Show All</option>
                {runtimes.map((type, index) => (<option key={index} value={type}>{type}</option>))}
              </select>
            )
          },
          {
            Header: 'Size',
            accessor: 'size',
            minWidth: 100,
            Cell: row => (formatBytes(row.value))
          },
          {
            Header: 'Memory',
            accessor: 'memory',
            minWidth: 100,
            Cell: row => (formatMegaBytes(row.value))
          },
          {
            Header: 'Invocations',
            columns: [
              {
                Header: 'Total',
                id: 'invocationsTotal',
                accessor: d => d.stats.invocations.total,
                minWidth: 100,
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? row.value : "N/A")
              },
              {
                Header: 'Success',
                id: 'invocationsSuccess',
                accessor: d => d.stats.invocations,
                minWidth: 125,
                filterable: false,
                Cell: row => (formatInvocationPercentage(row.value))
              }
            ]
          },
          {
            Header: 'Duration',
            columns: [
              {
                Header: 'Average',
                id: 'durationAverage',
                accessor: d => d.stats.duration.average,
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? formatMillisecondsDuration(row.value) : "N/A")
              },
              {
                Header: 'Maximum',
                id: 'durationMaximum',
                accessor: d => d.stats.duration.maximum,
                filterable: false,
                Cell: row => (row.value && row.value >= 0 ? formatMillisecondsDuration(row.value) : "N/A")
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
      <div className="clearfix resources lambda">
        <h4 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-code"/>
          &nbsp;
          Lambda
          {reportDate}
        </h4>
        {loading}
        {error}
        {list}
      </div>
    )
  }

}

LambdaComponent.propTypes = {
  accounts: PropTypes.arrayOf(PropTypes.object),
  data: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.arrayOf(PropTypes.shape({
      account: PropTypes.string.isRequired,
      reportDate: PropTypes.string.isRequired,
      function: PropTypes.shape({
        name: PropTypes.string.isRequired,
        description: PropTypes.string.isRequired,
        version: PropTypes.string.isRequired,
        lastModified: PropTypes.string.isRequired,
        runtime: PropTypes.string.isRequired,
        tags: PropTypes.object.isRequired,
        size: PropTypes.number.isRequired,
        memory: PropTypes.number.isRequired,
        stats: PropTypes.shape({
          invocations: PropTypes.shape({
            total: PropTypes.number,
            failed: PropTypes.number
          }),
          duration: PropTypes.shape({
            average: PropTypes.number,
            maximum: PropTypes.number
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
  data: aws.resources.Lambdas
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (date) => {
    dispatch(Actions.AWS.Resources.get.lambdas(date));
  },
  clear: () => {
    dispatch(Actions.AWS.Resources.clear.lambdas());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(LambdaComponent);
