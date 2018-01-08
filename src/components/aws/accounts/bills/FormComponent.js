import React, { Component } from 'react';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogTitle,
} from 'material-ui/Dialog';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../../common/forms';
import Popover from '../../../misc/Popover';
import PropTypes from "prop-types";

const Validation = Validations.AWSAccount;

// Form Component for new AWS Account
class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.submit = this.submit.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false});
  };

  submit = (e) => {
    e.preventDefault();
    this.closeDialog(e);
    const formValues = this.form.getValues();
    const bucketValues = Validation.getS3BucketValues(formValues.bucket);
    let bill = {
      bucket: bucketValues[0],
      prefix: bucketValues[1]
    };
    this.props.submit(this.props.account, bill);
  };

  render() {
    const button = (this.props.bill !== undefined ? (
      <div>
        <span className="glyphicon glyphicon-pencil" aria-hidden="true"/>&nbsp;Save
      </div>
    ) : (
      <div>
        <i className="fa fa-plus" />&nbsp;Add
      </div>
    ));

    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog}>
          {this.props.bill !== undefined ? "Edit" : "Add"}
        </button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>{this.props.bill !== undefined ? "Edit this" : "Add a"} bill location</h1></DialogTitle>

          <DialogContent>

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="bucket">S3 Bucket</label>
                  &nbsp;
                  <Popover info popOver="Name of S3 bucket and path to bills"/>
                </div>
                <Input
                  name="bucket"
                  type="text"
                  className="form-control"
                  value={(this.props.bill !== undefined ? this.props.bill.bucket : "")}
                  validations={[Validation.required, Validation.s3BucketFormat]}
                />
              </div>

              <DialogActions>

                <button className="btn btn-default pull-left" onClick={this.closeDialog}>
                  Cancel
                </button>

                <Button
                  className="btn btn-primary btn-block"
                  type="submit"
                >
                  {this.props.bill !== undefined ? "Save" : "Add"}
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
  bill: PropTypes.shape({
    bucket: PropTypes.string.isRequired,
    prefix: PropTypes.string.isRequired
  }),
  submit: PropTypes.func.isRequired
};


export default FormComponent;
