import React, { Component } from 'react';
import Components from '../../../components';
import { connect } from 'react-redux';
import PropTypes from "prop-types";

import Actions from "../../../actions";

const Panel = Components.Misc.Panel;
const List = Components.User.List;
const Form = Components.User.ViewerForm;

class ViewersContainer extends Component {

  constructor(props) {
    super(props);
    this.state = {
      addViewerDialogOpen: false
    };
    this.openDialog = this.openDialog.bind(this);
  }

  componentWillMount() {
    this.props.getViewers();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.lastCreated !== nextProps.lastCreated)
      this.props.getViewers();
  }

  openDialog = (addViewerDialogOpen) => (event) => {
    event.preventDefault();
    this.setState({ addViewerDialogOpen });
  };

  render() {
    return (
      <Panel>
        <div>
          <h3 className="white-box-title no-padding inline-block">
            <i className="fa fa-users white-box-title-icon" aria-hidden="true"/>
            Read-only users
          </h3>
          <div className="inline-block pull-right">
            <Form submit={this.props.viewerActions.new} viewer={this.props.lastCreated} clear={this.props.viewerActions.clearNew}/>
          </div>
        </div>
        <List
          viewers={this.props.viewers}
          viewerActions={this.props.viewerActions}
        />
      </Panel>
    )
  }
}

ViewersContainer.propTypes = {
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
  lastCreated: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.shape({
      id: PropTypes.number.isRequired,
      email: PropTypes.string.isRequired,
      password: PropTypes.string
    }),
  }),
  getViewers: PropTypes.func.isRequired,
  viewerActions: PropTypes.shape({
    new: PropTypes.func.isRequired,
    clearNew: PropTypes.func.isRequired
  }).isRequired,
};

const mapStateToProps = ({ user }) => ({
  viewers: user.viewers.all,
  lastCreated: user.viewers.creation,
});

const mapDispatchToProps = (dispatch) => ({
  getViewers: () => dispatch(Actions.User.getViewers()),
  viewerActions: {
    new: (email) => {
      dispatch(Actions.User.createViewer(email));
    },
    clearNew: () => {
      dispatch(Actions.User.clearCreate());
    }
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(ViewersContainer);
