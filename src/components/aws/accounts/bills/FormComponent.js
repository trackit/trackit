import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Spinner from 'react-spinkit';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import {CopyToClipboard} from 'react-copy-to-clipboard';
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
      open: false,
      step: 0,
      bucket : (this.props.bill !== undefined ? this.props.bill.bucket : ""),
      prefix : (this.props.bill !== undefined ? this.props.bill.prefix : "")
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.submit = this.submit.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({
      open: true,
      bucket : (this.props.bill !== undefined ? this.props.bill.bucket : ""),
      prefix : (this.props.bill !== undefined ? this.props.bill.prefix : "")
    });
    this.props.clear();
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false, step: 0});
    this.props.clear();
  };

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({
      [name]: value
    });
  }

  submit = (e) => {
    e.preventDefault();
    if (!this.state.step) {
      this.setState({ step : 1 })
    } else {
      const body = {
        bucket: this.state.bucket,
        prefix: this.state.prefix
      };
      this.props.submit(body);
    }
  };

  componentWillReceiveProps(nextProps) {
    if (nextProps.status && nextProps.status.status && nextProps.status.value && !nextProps.status.hasOwnProperty("error")) {
      this.setState({open: false, step: 0});
    }
  }

  getBucketPolicy() {
    const bucketString = this.state.prefix.length ? `${this.state.bucket}/${this.state.prefix}` : this.state.bucket;

    return(
      `{
        "Version": "2008-10-17",
        "Id": "PolicyAccessTrackitBucket",
        "Statement": [
          {
            "Sid": "Stmt1",
            "Effect": "Allow",
            "Principal": {
              "AWS": "arn:aws:iam::386209384616:root"
            },
            "Action": [
              "s3:GetBucketAcl",
              "s3:GetBucketPolicy"
            ],
            "Resource": "arn:aws:s3:::${bucketString}"
          },
          {
            "Sid": "Stmt2",
            "Effect": "Allow",
            "Principal": {
              "AWS": "arn:aws:iam::386209384616:root"
            },
            "Action": [
              "s3:PutObject"
            ],
            "Resource": "arn:aws:s3:::${bucketString}/*"
          }
        ]
      }`
    );
  }


  getPolicy(bucket, prefix) {
      const bills = [];
      if (this.props.bills && this.props.bills.status && this.props.bills.values && this.props.bills.values.length) {
        for (let i = 0; i < this.props.bills.values.length; i++) {
          const element = this.props.bills.values[i];
          const bucketString = element.prefix.length ? `${element.bucket}/${element.prefix}` : element.bucket;
          bills.push(
            `arn:aws:s3:::${bucketString}`
          );
        }
      }
      const newBucketString = prefix.length ? `${bucket}/${prefix}` : bucket;
      bills.push(`arn:aws:s3:::${newBucketString}`);

      const base = {
        Version: "2012-10-17",
        Statement: [
            {
                Effect: "Allow",
                Action: "s3:GetObject",
                Resource: bills
            },
            {
                Effect: "Allow",
                Action: [
                    "s3:GetBucketLocation",
                    "s3:ListBucket"
                ],
                Resource: bills
            },
            {
                Effect: "Allow",
                Action: [
                    "sts:GetCallerIdentity",
                    "rds:DescribeDBInstances",
                    "cloudwatch:GetMetricStatistics",
                    "ec2:DescribeRegions",
                    "ec2:DescribeInstances"
                ],
                Resource: "*"
            }
        ]
    }
    return(JSON.stringify(base, null, 2));

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

     const hidePolicy = (
      <div className="alert alert-info">
        <i className="fa fa-arrow-up"></i>
        &nbsp;
        Please enter your bucket name in step 3 to see the policy
      </div>
    );

    const policyElement = (
      <pre style={{ height: '85px', marginTop: '10px' }}>
        {this.getBucketPolicy()}
      </pre>
    );

      

    const tutorial = (
      <div className="tutorial">

        <ol>
          <li>Go to your <a rel="noopener noreferrer" target="_blank" href="https://s3.console.aws.amazon.com/s3/home">AWS Console S3 page</a>.</li>
          <li>Click <strong>Create bucket</strong> and input a name of your choice for your bucket. You can then complete the next wizard steps without changing the default values.</li>
          <li>
            Please fill the name of the bucket you created at <strong>Step 2</strong> in the Form below. <i className="fa fa-arrow-down"/>
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
                  value={this.state.bucket}
                  onChange={this.handleInputChange}
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
                  value={this.state.prefix}
                  onChange={this.handleInputChange}
                  validations={[Validation.s3PrefixFormat]}
                />
              </div>
          </li>
          <li>Still on your  <a rel="noopener noreferrer" target="_blank" href="https://s3.console.aws.amazon.com/s3/home">AWS Console S3 page</a> click on the name of the bucket you just created.</li>
          <li>Go to the <strong>Permissions</strong> tab and select <strong>Bucket Policy</strong></li>
          <li>
            Paste the following into the Bucket policy Editor and <strong>Save</strong>:
            {this.state.bucket.length ? (
            <CopyToClipboard text={this.getBucketPolicy()}>
                  <div className="badge">
                    <i className="fa fa-clipboard" aria-hidden="true"/>
                  </div>
            </CopyToClipboard>

            ): null}
            {this.state.bucket.length ? policyElement : hidePolicy}
          </li>
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
        </ol>

      </div>
    );

    const tutorialEditPolicy = (
      <div className="tutorial">
        <div className="alert alert-info">
          Please note that if you used our generated AWS Policy you will have to update it when adding or editing a bill location.
          <br/>
          <br/>
          To do that, go to your <a href="https://console.aws.amazon.com/iam/home#/policies">AWS Policies List</a>, select TrackIt policy and Edit it.
          Paste the following policy into the JSON Editor and submit.
          <br/>
          <br/>
          <strong>
            If you used AWS ReadOnlyAccess policy when setting up your account please ignore this step
          </strong>
        </div>
        <h5>
          Updated policy
          <CopyToClipboard text={this.getPolicy(this.state.bucket, this.state.prefix)}>
            <div className="badge">
              <i className="fa fa-clipboard" aria-hidden="true"/>
            </div>
          </CopyToClipboard>
        </h5>
        <pre style={{ height: '180px', marginTop: '10px' }}>
          {this.getPolicy(this.state.bucket, this.state.prefix)}
        </pre>
      </div>
    );

    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog} disabled={this.props.disabled}>
          {this.props.bill !== undefined ? <i className="fa fa-edit"/> : <i className="fa fa-plus"/>}
          &nbsp;
          {this.props.bill !== undefined ? "Edit" : "Add a bill location"}
        </button>

        <Dialog open={this.state.open} fullWidth maxWidth="md">

          <DialogTitle disableTypography><h1>
            <i className="fa fa-shopping-basket red-color"/>
            &nbsp;
            {this.props.bill !== undefined ? "Edit this" : "Add a"} bill location
          </h1></DialogTitle>

          <DialogContent>

            {loading || error}

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>
              {this.props.bill === undefined && this.state.step === 0 && tutorial}
              {this.props.bill !== undefined &&
                <div>
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
                      value={this.state.bucket}
                      onChange={this.handleInputChange}
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
                      value={this.state.prefix}
                      onChange={this.handleInputChange}
                      validations={[Validation.s3PrefixFormat]}
                    />
                  </div>
                </div>
              }

              {this.state.step === 1 && tutorialEditPolicy}

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
  clear: PropTypes.func.isRequired,
  bills: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        error: PropTypes.string.isRequired,
        bucket: PropTypes.string.isRequired,
        prefix: PropTypes.string.isRequired
      })
    )
  }),
  disabled: PropTypes.bool
};

FormComponent.defaultProps = {
  disabled: false
};

export default FormComponent;
