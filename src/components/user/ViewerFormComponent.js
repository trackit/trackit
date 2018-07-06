import React, { Component } from 'react';
import PropTypes from "prop-types";
import Dialog, {
  DialogActions,
  DialogContent,
  DialogTitle,
} from 'material-ui/Dialog';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Spinner from 'react-spinkit';
import {CopyToClipboard} from 'react-copy-to-clipboard';
import Validations from '../../common/forms';

const Validation = Validations.Auth;

class ViewerFormComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      addViewerDialogOpen: false
    };
    this.openDialog = this.openDialog.bind(this);
    this.submit = this.submit.bind(this);
  }

  componentWillMount() {
    this.props.clear();
  }

  openDialog = (addViewerDialogOpen) => (event) => {
    event.preventDefault();
    if (!addViewerDialogOpen)
      this.props.clear();
    this.setState({ addViewerDialogOpen });
  };

  submit = (e) => {
    e.preventDefault();
    const values = this.form.getValues();
    this.props.submit(values.email);
  };

  render() {
    const loading = (!this.props.viewer.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.viewer.error ? (<div className="alert alert-warning" role="alert">{this.props.viewer.error.message}</div>) : null);

    const password = (this.props.viewer.status && this.props.viewer.value ? (
      <div>
        <div className="alert alert-info">
          Please save this password safely, you will not be able to see it again later.
        </div>

        <div className="form-group">
          <div className="input-title">
            <label htmlFor="email">Email</label>
          </div>
          <Input
            name="email"
            type="email"
            className="form-control"
            value={this.props.viewer.value.email}
            validations={[Validation.required, Validation.email]}
            disabled
          />
        </div>

        <div className="form-group">
          <div className="input-title">
            <label htmlFor="email">Password</label>
            <CopyToClipboard text={this.props.viewer.value.password}>
              <div className="badge viewer">
                <i className="fa fa-clipboard" aria-hidden="true"/>
              </div>
            </CopyToClipboard>
          </div>
          <Input
            name="password"
            type="text"
            className="form-control"
            value={this.props.viewer.value.password}
            validations={[Validation.required]}
            disabled
          />
        </div>
      </div>
    ) : null);

    const passwordActions = (password ? (
      <DialogActions>
        <button className="btn btn-default pull-left" onClick={this.openDialog(false)}>
          Close
        </button>
      </DialogActions>
    ) : null);

    const form = (!loading && !password ? (
      <div>
        <div className="alert alert-info">
          Email for the user you will create and give read-only access to your data. The password will be generated later.
        </div>
        <div className="form-group">
          <div className="input-title">
            <label htmlFor="email">Email</label>
          </div>
          <Input
            name="email"
            type="email"
            className="form-control"
            validations={[Validation.required, Validation.email]}
          />
        </div>
      </div>
    ) : null);

    const formActions = (form ? (
      <DialogActions>
        <button className="btn btn-default pull-left" onClick={this.openDialog(false)}>
          Cancel
        </button>
        <Button
          className="btn btn-primary btn-block"
          type="submit"
        >
          Add
        </Button>
      </DialogActions>
    ) : null);

    return (
      <div>
        <button className="btn btn-default" onClick={this.openDialog(true)}>
          Add
        </button>
        <Dialog
          open={this.state.addViewerDialogOpen}
          fullWidth
        >
          <DialogTitle disableTypography><h1>Add a viewer</h1></DialogTitle>
          <DialogContent>
            {loading}
            {error}
            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit} >
              {form}
              {password}
              {formActions}
              {passwordActions}
            </Form>
          </DialogContent>
        </Dialog>
      </div>
    )
  }
}

ViewerFormComponent.propTypes = {
  submit: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
  viewer: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.shape({
      id: PropTypes.number.isRequired,
      email: PropTypes.string.isRequired,
      password: PropTypes.string
    }),
  }),
};

export default ViewerFormComponent;
