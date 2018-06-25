import React, { Component } from 'react';
import Components from '../../../components';
import { connect } from 'react-redux';

import List, {
  ListItem,
  ListItemText,
} from 'material-ui/List';

import Dialog, {
  DialogContent,
  DialogTitle,
} from 'material-ui/Dialog';

import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';

import Actions from "../../../actions";

const Panel = Components.Misc.Panel;

class NewViewerForm extends Component {
  state = { email: '', password: null }

  createViewer = () => this.props.createViewer(this.state.email)

  render() {
    return (
      <div>
        <TextField
          onChange={ event => this.setState({ email: event.target.value }) }
          value={ this.state.email }
          fullWidth
          label='Email'
          helperText='Email for the user you will create and give read-only access to your data. The password will be generated later.'
        />
        <Button
          onClick={ this.createViewer }
        >
          Create
        </Button>
      </div>
    )
  }
}

class ViewersContainer extends Component {
  state = { addViewerDialogOpen: false };

  componentWillMount() {
    console.log('viewers container props', this.props);
    this.props.getViewers();
  }

  openDialog = addViewerDialogOpen => event => {
    event.preventDefault();
    this.setState({ addViewerDialogOpen });
  }

  render() {
    console.log('viewer container props', this.props)
    return (
      <Panel>
        <Dialog
          open={this.state.addViewerDialogOpen}
          onBackdropClick={this.openDialog(false)}
          fullWidth
        >
          <DialogTitle disableTypography><h1>My Team</h1></DialogTitle>
          <DialogContent>
            <NewViewerForm createViewer={ this.props.createViewer } />
          </DialogContent>
        </Dialog>
        <div>
          <h3 className="white-box-title no-padding inline-block">Read-only users</h3>
          <div className="inline-block pull-right">
            <button className="btn btn-default" onClick={this.openDialog(true)}>Add</button>
          </div>
          <List disablePadding className='accounts-list'>
          {
            (this.props.viewers.values || []).map(viewer => (
              <ListItem key={viewer.id}>
                <ListItemText
                  primary={viewer.email}
                  secondary={ `Password: ${viewer.password || '●●●●●●●●'}` }
                />
              </ListItem>
            ))
          }
          </List>
        </div>
      </Panel>
    )
  }
}

const mapStateToProps = state => {console.log(state); return {
  viewers: state.user.viewers.all,
  lastViewerCreated: state.user.viewers.lastCreated,
}};

const mapDispatchToProps = dispatch => ({
  getViewers: () => dispatch(Actions.User.getViewers()),
  createViewer: email => dispatch(Actions.User.createViewer(email)),
})

export default connect(mapStateToProps, mapDispatchToProps)(ViewersContainer);
