import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Spinner from 'react-spinkit';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../../common/forms';
import Popover from '../../../misc/Popover';
import PropTypes from "prop-types";
import Misc from '../../../misc';
import Reports_first from '../../../../assets/report_step_1.png';
import Reports_second from '../../../../assets/report_step_2.png';

const Picture = Misc.Picture;
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
    this.props.clear();
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false});
    this.props.clear();
  };

  submit = (e) => {
    e.preventDefault();
    const formValues = this.form.getValues();
    this.props.submit(formValues);
  };

  componentWillReceiveProps(nextProps) {
    if (nextProps.status && nextProps.status.status && nextProps.status.value && !nextProps.status.hasOwnProperty("error")) {
      this.setState({open: false});
    }
  }

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

    const loading = (this.props.status && !this.props.status.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

     const error = (this.props.status && this.props.status.status && this.props.status.hasOwnProperty("error") ? (
         <div className="alert alert-warning" role="alert">{this.props.status.error.message}</div>
     ) : null);

    const tutorial = (
      <div className="tutorial">

        <ol>
          <li>Go to your <a rel="noopener noreferrer" target="_blank" href="https://s3.console.aws.amazon.com/s3/home">AWS Console S3 page</a>.</li>
          <li>Click <strong>Create bucket</strong> and input a name of your choice for your bucket. You can then complete the next wizard steps without changing the default values.</li>
          <li>Then go to your <a rel="noopener noreferrer" target="_blank" href="https://console.aws.amazon.com/billing/home#/reports">Billing Reports setup page</a> and click <strong>Create report</strong></li>
          <li>
            Choose a report name, select <strong>Hourly</strong> as the <strong>Time unit</strong> and Include <strong>Resources IDs</strong> (see screenshot). You can then click <strong>Next</strong>.
            <br />
            <Picture
              src={Reports_first}
              alt="Reports settings 1 tutorial"
              button={<strong>( Click here to see screenshot )</strong>}
            />
          </li>
          <li>
            In <strong>S3 Bucket</strong> input the name of the bucket you created at <strong>Step 2</strong>, select <strong>GZIP</strong> as <strong>Compression</strong>, then Submit. You can then review your settings and Submit again.
            <br />
            <Picture
              src={Reports_second}
              alt="Reports settings 2 tutorial"
              button={<strong>( Click here to see screenshot )</strong>}
            />
          </li>
          <li>
            You are almost done !
            <br/>
            Please fill the name of the bucket you created at <strong>Step 2</strong> in the Form below. <i className="fa fa-arrow-down"/>
            <br/>
            <strong>That's it ! </strong><i className="fa fa-smile-o"/>
          </li>
        </ol>

      </div>
    );

    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog}>
          {this.props.bill !== undefined ? "Edit" : "Add"}
        </button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>
            <i className="fa fa-shopping-basket red-color"/>
            &nbsp;
            {this.props.bill !== undefined ? "Edit this" : "Add a"} bill location
          </h1></DialogTitle>

          <DialogContent>

            {loading || error}
            {this.props.bill === undefined && tutorial}

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="bucket">S3 Bucket name</label>
                  &nbsp;
                  <Popover info tooltip="Name of the S3 bucket you created"/>
                </div>
                <Input
                  name="bucket"
                  type="text"
                  className="form-control"
                  placeholder="Bucket Name"
                  value={(this.props.bill !== undefined ? this.props.bill.bucket : "")}
                  validations={[Validation.required, Validation.s3BucketNameFormat]}
                />
              </div>

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="bucket">Report path prefix (optional)</label>
                  &nbsp;
                  <Popover info tooltip="If you set a path prefix when creating your report"/>
                </div>
                <Input
                  name="prefix"
                  type="text"
                  className="form-control"
                  placeholder="Optional prefix"
                  value={(this.props.bill !== undefined ? this.props.bill.prefix : "")}
                  validations={[Validation.s3PrefixFormat]}
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
  status: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.object
  }),
  submit: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired
};


export default FormComponent;
