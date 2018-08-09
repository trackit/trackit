import React, { Component } from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Spinner from 'react-spinkit';
import PropTypes from 'prop-types';
//import Misc from '../misc';

//const DeleteConfirmation = Misc.DeleteConfirmation;

export class Item extends Component {

  constructor(props) {
    super(props);
    this.deleteViewer = this.deleteViewer.bind(this);
  }

  deleteViewer = () => {
    console.log("Viewer deletion is not available yet");
//    this.props.viewerActions.delete(this.props.viewer.id);
  };

  render() {
    return (
      <div>

        <ListItem divider>

          <ListItemText
            disableTypography
            className="viewer-name"
            primary={this.props.viewer.email}
          />

          <div className="actions">

            <div className="inline-block">
              {/*'<DeleteConfirmation entity="viewer" confirm={this.deleteViewer}/>'*/}
            </div>

          </div>

        </ListItem>

      </div>
    );
  }

}

Item.propTypes = {
  viewer: PropTypes.shape({
    id: PropTypes.number.isRequired,
    email: PropTypes.string.isRequired,
    password: PropTypes.string,
  }),
  viewerActions: PropTypes.shape({
    delete: PropTypes.func,
  }).isRequired,
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    const loading = (!this.props.viewers.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.viewers.error ? ` (${this.props.viewers.error.message})` : null);
    const noViewers = (this.props.viewers.status && (!this.props.viewers.values || !this.props.viewers.values.length || error) ? <div className="alert alert-warning" role="alert">No viewer available{error}</div> : "");

    const viewers = (this.props.viewers.status && this.props.viewers.values && this.props.viewers.values.length ? (
      this.props.viewers.values.map((viewer, index) => (
        <Item
          key={index}
          viewer={viewer}
          viewerActions={this.props.viewerActions}
        />
      ))
    ) : null);

    return (
      <List disablePadding className="viewers-list">
        {loading}
        {noViewers}
        {viewers}
      </List>
    );
  }

}

ListComponent.propTypes = {
  viewers: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        email: PropTypes.string.isRequired,
        password: PropTypes.string
      })
    ),
  }),
  viewerActions: PropTypes.shape({
    delete: PropTypes.func,
  }).isRequired,
};

export default ListComponent;
