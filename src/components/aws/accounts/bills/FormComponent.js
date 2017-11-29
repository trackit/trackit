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
    let values = this.form.getValues();
    let bill = {
      bucket: values.bucket,
      path: values.path
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

          <DialogTitle>{this.props.bill !== undefined ? "Edit this" : "Add a"} bill location</DialogTitle>

          <DialogContent>

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>

              <div className="form-group">
                <label htmlFor="bucket">S3 Bucket</label>
                <Input
                  name="bucket"
                  type="text"
                  className="form-control"
                  value={(this.props.bill !== undefined ? this.props.bill.bucket : "s3://")}
                  validations={[Validation.required, Validation.s3BucketFormat]}
                />
              </div>

              <div className="form-group">
                <label htmlFor="path">Path</label>
                <Input
                  type="text"
                  name="path"
                  value={(this.props.bill !== undefined ? this.props.bill.path : undefined)}
                  className="form-control"
                  validations={[Validation.required, Validation.pathFormat]}
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
    path: PropTypes.string.isRequired
  }),
  submit: PropTypes.func.isRequired
};


export default FormComponent;
