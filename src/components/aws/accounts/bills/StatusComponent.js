import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Moment from 'moment';
import Spinner from 'react-spinkit';
import PropTypes from "prop-types";

class Item extends Component {

  render() {
    const lastError = (this.props.status.lastError && this.props.status.lastError.length ? (
      <div className="alert alert-warning" role="alert">
        {this.props.status.lastError}
      </div>
    ) : null);
    const pendingBadge = (this.props.status.nextPending ? (<span class="badge">In progress</span>) : null);
    return (
      <ListItem divider className="status-list-item">
        <div>
          <h3>{this.props.status.AwsAccountPretty}</h3>
          <h4>{`s3://${this.props.status.bucket}/${this.props.status.prefix}`}</h4>
        </div>
        <div>
          <div className="info">
            <h5><i className="fa fa-clock-o"/>&nbsp;Last import</h5>
            {lastError}
            Started : {Moment(this.props.status.lastStarted).format('MMM Do Y HH:mm:ss')}
            <br/>
            Finished : {Moment(this.props.status.lastFinished).format('MMM Do Y HH:mm:ss')}
          </div>
          <div className="info">
            <h5><i className="fa fa-cloud-download"/>&nbsp;Next import {pendingBadge}</h5>
            Planned : {Moment(this.props.status.nextStarted).format('MMM Do Y HH:mm')}
          </div>
        </div>
      </ListItem>
    );
  }

}

Item.propTypes = {
  status: PropTypes.shape({
    BillRepositoryId: PropTypes.number.isRequired,
    AwsAccountPretty: PropTypes.string.isRequired,
    AwsAccountId: PropTypes.number.isRequired,
    bucket: PropTypes.string.isRequired,
    prefix: PropTypes.string.isRequired,
    nextStarted: PropTypes.string.isRequired,
    nextPending: PropTypes.bool.isRequired,
    lastStarted: PropTypes.string.isRequired,
    lastFinished: PropTypes.string.isRequired,
    lastError: PropTypes.string.isRequired
  })
};

// Status Component for new AWS Account
class StatusComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    this.props.billsStatusActions.get();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false});
    this.props.billsStatusActions.clear();
  };

  render() {
    const loading = (!this.props.bills.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.bills.error ? ` (${this.props.bills.error.message})` : null);
    const noBills = (this.props.bills.status && (!this.props.bills.values || !this.props.bills.values.length || error) ? <div className="alert alert-warning" role="alert">No bills available{error}</div> : "");

    const values = (this.props.bills.status && this.props.bills.hasOwnProperty("values") ? (this.props.bills.values.map((item, index) => (<Item key={index} status={item}/>))) : null);

    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog}>
          <i className="fa fa-heartbeat"/>
          &nbsp;
          Status
        </button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>Import Status</h1></DialogTitle>

          <DialogContent>

            <List className="status-list">
              {loading}
              {noBills}
              {values}
            </List>

            <DialogActions>

              <button className="btn btn-default pull-left" onClick={this.closeDialog}>
                Close
              </button>

            </DialogActions>

          </DialogContent>
        </Dialog>
      </div>
    );
  }

}

StatusComponent.propTypes = {
  bills: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        BillRepositoryId: PropTypes.number.isRequired,
        AwsAccountPretty: PropTypes.string.isRequired,
        AwsAccountId: PropTypes.number.isRequired,
        bucket: PropTypes.string.isRequired,
        prefix: PropTypes.string.isRequired,
        nextStarted: PropTypes.string.isRequired,
        nextPending: PropTypes.bool.isRequired,
        lastStarted: PropTypes.string.isRequired,
        lastFinished: PropTypes.string.isRequired,
        lastError: PropTypes.string.isRequired
      })
    )
  }),
  billsStatusActions: PropTypes.shape({
    get: PropTypes.func.isRequired,
    clear: PropTypes.func.isRequired,
  }).isRequired,
};

export default StatusComponent;
