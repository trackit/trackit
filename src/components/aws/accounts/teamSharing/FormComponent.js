import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Spinner from 'react-spinkit';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Misc from '../../../misc';
import Validations from '../../../../common/forms';
import PropTypes from "prop-types";

const Selector = Misc.Selector;
const Validation = Validations.AWSAccount;

const permissions = {
  2: "Read-Only",
  1: "Standard",
  0: "Administrator"
};

// Form Component for new AWS Account
class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      permission: (props.AccountViewer ? props.AccountViewer.level : 2),
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.setPermission = this.setPermission.bind(this);
    this.submit = this.submit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.status && nextProps.status.status && nextProps.status.values && !nextProps.status.hasOwnProperty("error")) {
      this.setState({open: false});
    }
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({
      open: true
    });
    this.props.clear();
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false});
    this.props.clear();
  };

  setPermission(permission) {
    this.setState({permission: parseInt(permission, 10)});
  }

  submit = (e) => {
    e.preventDefault();
    const formValues = this.form.getValues();
    this.props.submit(formValues.email, this.state.permission);
  };

  render() {
    const loading = (this.props.status && !this.props.status.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

    const error = (this.props.status && this.props.status.status && this.props.status.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">{this.props.status.error.message}</div>
    ) : null);

    const allowedPermissions = {};
    Object.keys(permissions).forEach((permission) => {
      if (this.props.permissionLevel <= permission)
        allowedPermissions[permission] = permissions[permission];
    });

    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog} disabled={this.props.disabled}>
          {this.props.AccountViewer !== undefined ? <i className="fa fa-edit"/> : <i className="fa fa-plus"/>}
          &nbsp;
          {this.props.AccountViewer !== undefined ? "Edit" : "Invite a user"}
        </button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>
            <i className="fa fa-users red-color"/>
            &nbsp;
            {this.props.AccountViewer !== undefined ? "Edit this shared user" : "Invite a user"}
          </h1></DialogTitle>

          <DialogContent>

            {loading || error}

            <Form ref={
                /* istanbul ignore next */
                form => { this.form = form; }
              }
              className="team-sharing-form"
              onSubmit={this.submit}>


              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="guestmail">Email address</label>
                </div>
                <Input
                  name="email"
                  type="email"
                  className="form-control"
                  placeholder="john.doe@domain.com"
                  value={this.props.AccountViewer ? this.props.AccountViewer.email : undefined}
                  disabled={this.props.AccountViewer !== undefined}
                  onChange={this.handleInputChange}
                  validations={[Validation.required, Validation.guestMailFormat]}
                />
              </div>

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="guestpermission">Permission Level</label>
                </div>
                <Selector
                  values={allowedPermissions}
                  selected={this.state.permission}
                  selectValue={this.setPermission}
                />
              </div>

              <div className="alert alert-info">
                <strong>Administrator</strong> : Can add viewers and edit / delete viewers, bill repository or AWS Account (Has same rights as the account owner)
                <br/>
                <strong>Standard</strong> : Can add viewers and edit / delete viewers with Standard or Read-Only rights
                <br/>
                <strong>Read-Only</strong> : Can only see account data
              </div>

              <DialogActions>

                <button className="btn btn-default pull-left" onClick={this.closeDialog}>
                  Cancel
                </button>

                <Button
                  className="btn btn-primary btn-block"
                  type="submit"
                >
                  {this.props.AccountViewer !== undefined ? "Save" : "Invite"}
                </Button>

              </DialogActions>

            </Form>

          </DialogContent>
        </Dialog>
      </div>
    );
  }

}

FormComponent.propTypes = {
  account: PropTypes.number.isRequired,
  AccountViewer: PropTypes.shape({
    sharedId: PropTypes.number.isRequired,
    email: PropTypes.string.isRequired,
    level: PropTypes.number.isRequired,
    userId: PropTypes.number.isRequired,
    sharingStatus: PropTypes.bool.isRequired
  }),
  status: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.object
  }),
  submit: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
  permissionLevel: PropTypes.number.isRequired,
  disabled: PropTypes.bool
};

FormComponent.defaultProps = {
  disabled: false
};

export default FormComponent;
