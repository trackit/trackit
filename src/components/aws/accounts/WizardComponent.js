import React, {Component} from 'react';
import PropTypes from "prop-types";
import Validations from "../../../common/forms";
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import {CopyToClipboard} from 'react-copy-to-clipboard';
import Spinner from 'react-spinkit';
import Misc from '../../misc';
import BucketForm from './bills/FormComponent';
import RoleCreation from '../../../assets/wizard-creation.png';
import RoleARN from '../../../assets/wizard-rolearn.png';
import Reports_first from '../../../assets/report_step_1.png';
import Reports_second from '../../../assets/report_step_2.png';
import '../../../styles/Wizard.css';

const Popover = Misc.Popover;
const Picture = Misc.Picture;
const Validation = Validations.AWSAccount;

export class StepRoleCreation extends Component {

  constructor() {
    super();
    this.state = {
      copied: null,
    };
    this.toggleCopied = this.toggleCopied.bind(this);
  }

  toggleCopied = (value, result) => {
    if (result) {
      this.setState({ copied: value });
      setTimeout(() => {
        this.setState({ copied: null });
      }, 3000);  
    }
  };

  submit = (e) => {
    e.preventDefault();
    this.props.next();
  };

  render() {
    return (
      <div className="step step-one">

        <Form ref={
          /* istanbul ignore next */
          form => { this.form = form; }
        } onSubmit={this.submit} >

          <div className="tutorial">

            <ol>
              <li>Go to your <a rel="noopener noreferrer" target="_blank" href="https://console.aws.amazon.com/iam/home#/roles">AWS Console IAM Roles page</a>.</li>
              <li>Click on <strong>Create Role</strong></li>
              <li>
                <div>
                  Follow this screenshot to configure your new role correctly,
                  <br/>
                  using the informations provided below and click Next:
                  <br/>
                  <Picture
                    src={RoleCreation}
                    alt="Role creation tutorial"
                    button={<strong>( Click here to see screenshot )</strong>}
                  />
                  <hr/>
                  Account ID : <strong className="value">{this.props.external.accountId}</strong>
                  <CopyToClipboard text={this.props.external.accountId} onCopy={this.toggleCopied}>
                    <div className="badge">
                      <i className="fa fa-clipboard" aria-hidden="true"/>
                      {this.state.copied === this.props.external.accountId && ' Copied'}
                    </div>
                  </CopyToClipboard>
                  <br/>
                  External : <strong className="value">{this.props.external.external}</strong>
                  <CopyToClipboard text={this.props.external.external} onCopy={this.toggleCopied}>
                    <div className="badge">
                      <i className="fa fa-clipboard" aria-hidden="true"/>
                      {this.state.copied === this.props.external.external && ' Copied'}
                    </div>
                  </CopyToClipboard>
                </div>
                <hr/>
              </li>
              {
                this.props.minimalPolicy ?
                <li>Select the policy you created in previous step</li>
                : <li>Select the <strong>ReadOnlyAccess</strong> policy</li>
              }
              <li>Set a name for this new role and validate</li>
            </ol>

          </div>

          <div className="form-group clearfix">
            <div className="btn-group col-md-5" role="group">
              <div className="btn btn-default btn-left" onClick={this.props.close}>Cancel</div>
              <div className="btn btn-default btn-left" onClick={this.props.back}>Previous</div>
            </div>
            <Button className="btn btn-primary col-md-5 btn-right" type="submit">Next</Button>
          </div>

        </Form>

      </div>
    );
  }

}

StepRoleCreation.propTypes = {
  external: PropTypes.shape({
    external: PropTypes.string.isRequired,
    accountId: PropTypes.string.isRequired,
  }),
  next: PropTypes.func.isRequired,
  back: PropTypes.func.isRequired,
  close: PropTypes.func.isRequired,
  minimalPolicy: PropTypes.bool,
};

export class StepNameARN extends Component {

  submit = (e) => {
    e.preventDefault();
    let values = this.form.getValues();
    let account = {
      roleArn: values.roleArn,
      pretty: values.pretty,
      external: this.props.external.external,
      payer: true,
    };
    this.props.submit(account);
  };

  render() {
    const error = (this.props.account && this.props.account.status && this.props.account.hasOwnProperty("error")) ? (
      <div className="alert alert-warning" role="alert">{this.props.account.error.message}</div>
    ) : (null);

    return (
      <div className="step step-two">

        <div className="tutorial">

          <ol>
            <li>In <strong>Roles list</strong> on <a rel="noopener noreferrer" target="_blank" href="https://console.aws.amazon.com/iam/home#/roles">AWS Console IAM Roles page</a>, select the role you created in previous step</li>
            <li>
              <div>
                Copy the Role ARN in <strong>role summary</strong> to the form below.
                <br/>
                Details are available in this screenshot
                <Picture
                  src={RoleARN}
                  alt="Role ARN tutorial"
                  button={<strong>( Click here to see screenshot )</strong>}
                />
              </div>
            </li>
            <li>You can set a name for this account, to help you to manage your accounts easily</li>
          </ol>

        </div>
        {error}
        <Form ref={
          /* istanbul ignore next */
          form => { this.form = form; }
        } onSubmit={this.submit} >

          <div className="form-group">
            <div className="input-title">
              <label htmlFor="roleArn">Role ARN</label>
              &nbsp;
              <Popover info tooltip="Amazon Resource Name for your role ( See step 2 )"/>
            </div>
            <Input
              name="roleArn"
              type="text"
              className="form-control"
              validations={[Validation.required, Validation.roleArnFormat]}
            />
          </div>

          <div className="form-group">
            <div className="input-title">
              <label htmlFor="pretty">Name</label>
              &nbsp;
              <Popover info tooltip="Choose a pretty name"/>
            </div>
            <Input
              type="text"
              name="pretty"
              className="form-control"
            />
          </div>

          <div className="form-group clearfix">
            <div className="btn-group col-md-5" role="group">
              <div className="btn btn-default btn-left" onClick={this.props.close}>Cancel</div>
              <div className="btn btn-default btn-left" onClick={this.props.back}>Previous</div>
            </div>
            <Button className="btn btn-primary col-md-5 btn-right" type="submit">{this.props.account.status ? "Next" : <Spinner className="spinner" name='circle' color="white"/>}</Button>
          </div>

        </Form>

      </div>
    );
  }

}

StepNameARN.propTypes = {
  external: PropTypes.shape({
    external: PropTypes.string.isRequired,
    accountId: PropTypes.string.isRequired,
  }),
  account: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string
    })
  }),
  submit: PropTypes.func.isRequired,
  next: PropTypes.func.isRequired,
  back: PropTypes.func.isRequired,
  close: PropTypes.func.isRequired,
};

export class StepBucket extends Component {

  constructor() {
    super();
    this.state = {
      bucketName: '',
      bucketPrefix: '',
    };
    this.handleInputChange = this.handleInputChange.bind(this);
    this.toggleCopied = this.toggleCopied.bind(this);
  };

  toggleCopied = (value, result) => {
    if (result) {
      this.setState({ copied: value });
      setTimeout(() => {
        this.setState({ copied: null });
      }, 3000);  
    }
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
    this.props.submit({bucket: this.state.bucketName, prefix: this.state.bucketPrefix});
    this.props.next();
  };

  getBucketPolicy() {
    const bucketString = this.state.bucketPrefix.length ? `${this.state.bucketName}/${this.state.bucketPrefix}` : this.state.bucketName;

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


  render() {

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
                  name="bucketName"
                  type="text"
                  className="form-control"
                  placeholder="Bucket Name"
                  value={this.state.bucketName}
                  onChange={this.handleInputChange}
                  validations={[Validation.required, Validation.s3BucketNameFormat]}
                />
              </div>
          </li>
          <li>Still on your  <a rel="noopener noreferrer" target="_blank" href="https://s3.console.aws.amazon.com/s3/home">AWS Console S3 page</a> click on the name of the bucket you just created.</li>
          <li>Go to the <strong>Permissions</strong> tab and select <strong>Bucket Policy</strong></li>
          <li>
            Paste the following into the Bucket policy Editor and <strong>Save</strong>:
            <CopyToClipboard text={this.getBucketPolicy()} onCopy={this.toggleCopied}>
                  <div className="badge">
                    <i className="fa fa-clipboard" aria-hidden="true"/>
                    {this.state.copied && ' Copied'}
                  </div>
            </CopyToClipboard>
            <pre style={{ height: '85px', marginTop: '10px' }}>
              {this.getBucketPolicy()}
            </pre>
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

    return (
      <div>


        <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>
              {tutorial}

              <div className="form-group clearfix">
                <div className="btn btn-default col-md-5 btn-left" onClick={this.props.close}>Cancel</div>
                <Button className="btn btn-primary col-md-5 btn-right" type="submit">Next</Button>
              </div>
            </Form>
      </div>
    );
  }

}

StepBucket.propTypes = {
  external: PropTypes.shape({
    external: PropTypes.string.isRequired,
    accountId: PropTypes.string.isRequired,
  }),
  account: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string
    })
  }),
  bill: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error)
  }),
  submit: PropTypes.func.isRequired,
  close: PropTypes.func.isRequired
};

export class StepPolicy extends Component {
  constructor() {
    super();
    this.state = {
      minimal: false
    };
    this.next = this.next.bind(this);
    this.toggleCopied = this.toggleCopied.bind(this);
  }

  toggleCopied = (value, result) => {
    if (result) {
      this.setState({ copied: value });
      setTimeout(() => {
        this.setState({ copied: null });
      }, 3000);  
    }
  };

  next() {
    this.props.setMinimalPolicy(this.state.minimal);
    this.props.next();
  }

  getPolicy() {
    const bucketString = this.props.bucketPrefix.length ? `${this.props.bucketName}/${this.props.bucketPrefix}` : this.props.bucketName;

    return(
      `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::${bucketString}/*"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "s3:GetBucketLocation",
                    "s3:ListBucket"
                ],
                "Resource": "arn:aws:s3:::${bucketString}"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "sts:GetCallerIdentity",
                    "rds:DescribeDBInstances",
                    "cloudwatch:GetMetricStatistics",
                    "ec2:DescribeRegions",
                    "ec2:DescribeInstances",
                    "ec2:DescribeReservedInstancesListings",
                    "ec2:DescribeReservedInstancesModifications",
                    "ec2:DescribeReservedInstancesOfferings",
                    "ec2:DescribeVolumes"
                ],
                "Resource": "*"
            }
        ]
    }`
    );
  }

  render() {
    const tutorial = (
        <div className="tutorial">
        <ol>
          <li>Go to your <a rel="noopener noreferrer" target="_blank" href="https://console.aws.amazon.com/iam/home#/policies">AWS Console IAM Policies page</a>.</li>
          <li>Click <strong>Create Policy</strong></li>
          <li>Select the <strong>JSON</strong> tab</li>
          <li>
            Paste the following into the JSON Editor:
            <CopyToClipboard text={this.getPolicy()} onCopy={this.toggleCopied}>
                  <div className="badge">
                    <i className="fa fa-clipboard" aria-hidden="true"/>
                    {this.state.copied && ' Copied'}
                  </div>
            </CopyToClipboard>
            <pre style={{ height: '180px', marginTop: '10px' }}>
              {this.getPolicy()}
            </pre>
            <div className="alert alert-info">
              <i className="fa fa-info-circle"></i>
              &nbsp;
              With this policy you only gives TrackIt access to the strict minimum it needs to be functional.
            </div>
          </li>
          <li>Click <strong>Review Policy</strong> to submit</li>
          <li>Give your policy a name and click <strong>Create policy</strong></li>
        </ol>
      </div>
    );

    const permissionChoice = (
      <div>
        <div className="alert alert-info">
          <i className="fa fa-info-circle"></i>
          &nbsp;
          TrackIt works best using AWS full ReadOnly policy.
          <br/>
          <br/>
          If you want to use a more restrictive policy we can generate one for you that strictly gives TrackIt access to the minimum it needs. However you might need to update it for future TrackIt features to work.
        </div>
        <button className="btn btn-primary btn-block" onClick={this.next}><strong>Use AWS ReadOnly policy</strong></button>
        <br/>
        <button className="btn btn-primary btn-block" onClick={() => {this.setState({ minimal: true })}}>Use minimal policy</button>
        <hr/>
      </div>
    );

    return (
      <div>
        {this.state.minimal ? tutorial : permissionChoice}

        <div className="form-group clearfix">
          <div className="btn-group col-md-5" role="group">
            <div className="btn btn-default btn-left" onClick={this.props.close}>Cancel</div>
            <div className="btn btn-default btn-left" onClick={this.props.back}>Previous</div>
          </div>
          {this.state.minimal && <button className="btn btn-primary col-md-5 btn-right" onClick={this.next}>Next</button>}
        </div>

      </div>
    );
  }
}

StepPolicy.propTypes = {
  bucketName: PropTypes.string.isRequired,
  bucketPrefix: PropTypes.string.isRequired,
  next: PropTypes.func.isRequired,
  back: PropTypes.func.isRequired,
  close: PropTypes.func.isRequired,
  setMinimalPolicy: PropTypes.func.isRequired,
}

class Wizard extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      activeStep: 0,
      bucket: '',
      prefix: '',
      minimalPolicy: false,
      billEditMode: false,
    };
    this.nextStep = this.nextStep.bind(this);
    this.previousStep = this.previousStep.bind(this);
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.setBucketValues = this.setBucketValues.bind(this);
    this.submit = this.submit.bind(this);
    this.submitBucket = this.submitBucket.bind(this);
    this.setMinimalPolicy = this.setMinimalPolicy.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    // Everything went well
    if (
      (nextProps.account.status && nextProps.account.value && !nextProps.account.hasOwnProperty("error"))
      && (nextProps.bill.status && nextProps.bill.value && !nextProps.bill.hasOwnProperty("error"))
    ) {
      this.closeDialog();
    }
    // Account was created successfully but not bill
    else if (
      (nextProps.account.status && nextProps.account.value && !nextProps.account.hasOwnProperty("error"))
      && (nextProps.bill.status && nextProps.bill.hasOwnProperty("error"))
    ) {
      this.setState({ billEditMode : true });
    }

  }

  setBucketValues(values) {
    const { bucket, prefix } = values;
    this.setState({ bucket, prefix });
  }

  setMinimalPolicy(value) {
    this.setState({ minimalPolicy: value });
  }

  submit(account) {
    const bucket = {
      bucket: this.state.bucket,
      prefix: this.state.prefix,
    }
    this.props.submitAccount(account, bucket);
  }

  submitBucket(bucket) {
    this.props.submitBucket(this.props.account.value.id, bucket);
  }

  nextStep = () => {
    const activeStep = this.state.activeStep + 1;
    this.setState({activeStep});
  };

  previousStep = () => {
    const activeStep = this.state.activeStep - 1;
    this.setState({activeStep});
  };

  openDialog = (e) => {
    e.preventDefault();
    this.setState({open: true, activeStep: 0});
    this.props.clearAccount();
    this.props.clearBucket();
  };

  closeDialog = (e=null) => {
    if (e)
      e.preventDefault();
    this.setState({open: false, activeStep: 0, bucket: '', prefix: '', billEditMode: false, minimalPolicy: false });
    this.props.clearAccount();
    this.props.clearBucket();
  };

  render() {

    let steps = [
      {
        title: "Create a billing bucket",
        label: "Billing bucket",
        component: <StepBucket account={this.props.account} bill={this.props.bill} next={this.nextStep} submit={this.setBucketValues} close={this.closeDialog}/>
      },
      {
        title: "Policy",
        label: "Policy",
        component: <StepPolicy bucketName={this.state.bucket} setMinimalPolicy={this.setMinimalPolicy} bucketPrefix={this.state.prefix} next={this.nextStep} back={this.previousStep} close={this.closeDialog}/>
      },
      {
        title: "Create a role",
        label: "Role creation",
        component: <StepRoleCreation external={this.props.external} minimalPolicy={this.state.minimalPolicy} next={this.nextStep} back={this.previousStep} close={this.closeDialog}/>
      },
      {
        title: "Add your role",
        label: "Role ARN & Name",
        component: <StepNameARN external={this.props.external} account={this.props.account} submit={this.submit} next={this.nextStep} back={this.previousStep} close={this.closeDialog}/>
      },
    ];

    const stepper = (
      <div>
        <div>
          {steps[this.state.activeStep].component}
        </div>

        <Stepper nonLinear activeStep={this.state.activeStep} className="account-wizard-stepper">
          {steps.map((step, index) => (
              <Step key={index}>
                <StepButton
                  className={"account-wizard-stepper-item " + (this.state.activeStep > index ? "completed" : (this.state.activeStep === index ? "current" : "")) }
                  completed={this.state.activeStep > index}
                >
                  {step.label}
                </StepButton>
              </Step>
            ))}
        </Stepper>
      </div>
    );

    let billEditMode;
    if (this.props.account && this.props.account.value && this.props.account.value.id) {
      billEditMode = (
        <div>
          <div className="alert alert-warning">
            Your account was created successfully but we could not access the billing bucket you specified.
            Please click on the <strong>Add a bill location</strong> button below and try to set it up again. Thank you !
          </div>
          <BucketForm
            account={this.props.account.value && this.props.account.value.id}
            submit={this.submitBucket}
            status={this.props.bill}
            clear={this.props.clearBucket}
          />
          <hr />
          <button className="btn btn-default" onClick={this.closeDialog}>Close</button>
        </div>
      );  
    }

    return(
      <div className="account-wizard">

        <button className="btn btn-default" onClick={this.openDialog}><i className="fa fa-plus"></i>&nbsp;Add</button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>Add an AWS account : {steps[this.state.activeStep].title}</h1></DialogTitle>

          <DialogContent>
            {this.state.billEditMode ? billEditMode : stepper}
          </DialogContent>

        </Dialog>

      </div>
    );
  }

}

Wizard.propTypes = {
  external: PropTypes.shape({
    external: PropTypes.string.isRequired,
    accountId: PropTypes.string.isRequired,
  }),
  account: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string
    })
  }),
  bill: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error)
  }),
  submitAccount: PropTypes.func.isRequired,
  submitBucket: PropTypes.func.isRequired,
  clearAccount: PropTypes.func.isRequired,
  clearBucket: PropTypes.func.isRequired,
};

export default Wizard;
