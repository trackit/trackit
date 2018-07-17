import React, {Component} from 'react';
import PropTypes from "prop-types";
import Validations from "../../../common/forms";
import Dialog, {
  DialogContent,
  DialogTitle,
} from 'material-ui/Dialog';
import Stepper, {
  Step,
  StepButton
} from 'material-ui/Stepper';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import {CopyToClipboard} from 'react-copy-to-clipboard';
import Spinner from 'react-spinkit';
import Misc from '../../misc';
import RoleCreation from '../../../assets/wizard-creation.png';
import RoleARN from '../../../assets/wizard-rolearn.png';
import '../../../styles/Wizard.css';

const Popover = Misc.Popover;
const Picture = Misc.Picture;
const Validation = Validations.AWSAccount;

export class StepOne extends Component {

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
              <li>Go to your <strong>AWS Console</strong></li>
              <li>In <strong>Services</strong> panel, select <strong>IAM</strong></li>
              <li>Choose <strong>Role</strong> on the left side menu</li>
              <li>Click on <strong>Create Role</strong></li>
              <li>
                <div>
                  Follow this screenshot to configure your new role correctly,
                  <br/>
                  using the informations provided below :
                  <br/>
                  <Picture
                    src={RoleCreation}
                    alt="Role creation tutorial"
                    button={<strong>( Click here to see screenshot )</strong>}
                  />
                  <hr/>
                  Account ID : <strong className="value">{this.props.external.accountId}</strong>
                  <CopyToClipboard text={this.props.external.accountId}>
                    <div className="badge">
                      <i className="fa fa-clipboard" aria-hidden="true"/>
                    </div>
                  </CopyToClipboard>
                  <br/>
                  External : <strong className="value">{this.props.external.external}</strong>
                  <CopyToClipboard text={this.props.external.external}>
                    <div className="badge">
                      <i className="fa fa-clipboard" aria-hidden="true"/>
                    </div>
                  </CopyToClipboard>
                </div>
                <hr/>
              </li>
              <li>Select <strong>ReadOnlyAccess</strong> policy</li>
              <li>Set a name for this new role and validate</li>
            </ol>

          </div>

          <div className="form-group clearfix">
            <button className="btn btn-default col-md-5 btn-left" onClick={this.props.close}>Cancel</button>
            <Button className="btn btn-primary col-md-5 btn-right" type="submit">Next</Button>
          </div>

        </Form>

      </div>
    );
  }

}

StepOne.propTypes = {
  external: PropTypes.shape({
    external: PropTypes.string.isRequired,
    accountId: PropTypes.string.isRequired,
  }),
  next: PropTypes.func.isRequired,
  close: PropTypes.func.isRequired
};

export class StepTwo extends Component {

  submit = (e) => {
    e.preventDefault();
    let values = this.form.getValues();
    let account = {
      roleArn: values.roleArn,
      pretty: values.pretty,
      external: this.props.external.external
    };
    this.props.submit(account);
  };

  componentWillReceiveProps(nextProps) {
    if (nextProps.account.status && nextProps.account.value && !nextProps.account.hasOwnProperty("error"))
      nextProps.next();
  }

  render() {
    const error = (this.props.account && this.props.account.status && this.props.account.hasOwnProperty("error")) ? (
      <div className="alert alert-warning" role="alert">{this.props.account.error.message}</div>
    ) : (null);

    return (
      <div className="step step-two">

        <div className="tutorial">

          <ol>
            <li>In <strong>Role</strong> list, select the role you created in previous step</li>
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
              <Popover info popOver="Amazon Resource Name for your role ( See step 2 )"/>
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
              <Popover info popOver="Choose a pretty name"/>
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

StepTwo.propTypes = {
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
  close: PropTypes.func.isRequired
};

export class StepThree extends Component {

  submit = (e) => {
    e.preventDefault();
    const formValues = this.form.getValues();
    const bucketValues = Validation.getS3BucketValues(formValues.bucket);
    let bill = {
      bucket: bucketValues[0],
      prefix: bucketValues[1]
    };
    this.props.submit(this.props.account.value.id, bill);
  };

  componentWillReceiveProps(nextProps) {
    if (nextProps.bill.status && nextProps.bill.value)
      nextProps.close();
  }

  render() {
    const error = (this.props.bill && this.props.bill.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">{this.props.bill.error.message}</div>
    ) : null);

    return (
      <div className="step step-three">

        {error}

        <div className="tutorial">

          <ol>
            <li>Fill the form with the location of a <strong>S3 bucket</strong> that contains bills
              <br/>
              Example : <code>s3://my.bucket/bills</code>
            </li>
            <li>You will be able to add more buckets later.</li>
          </ol>

        </div>

        <Form
          ref={
            /* istanbul ignore next */
            form => { this.form = form; }
          }
          onSubmit={this.submit}
        >

          <div className="form-group">
            <div className="input-title">
              <label htmlFor="bucket">S3 Bucket</label>
              &nbsp;
              <Popover info popOver="Name of S3 bucket and path to bills"/>
            </div>
            <div className="input-group">
              <div className="input-group-addon">s3://</div>
              <Input
                name="bucket"
                type="text"
                className="form-control"
                placeholder="<bucket-name>/<path>"
                validations={[Validation.required, Validation.s3BucketFormat]}
              />
            </div>
          </div>

          <div className="form-group clearfix">
            <div className="btn btn-default col-md-5 btn-left" onClick={this.props.close}>Cancel</div>
            <Button className="btn btn-primary col-md-5 btn-right" type="submit" disabled={!this.props.account}>{!this.props.bill || this.props.bill.status ? "Done" : <Spinner className="spinner" name='circle' color="white"/>}</Button>
          </div>

        </Form>

      </div>
    );
  }

}

StepThree.propTypes = {
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

class Wizard extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      activeStep: 0
    };
    this.nextStep = this.nextStep.bind(this);
    this.previousStep = this.previousStep.bind(this);
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
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
    this.setState({open: false, activeStep: 0});
    this.props.clearAccount();
    this.props.clearBucket();
  };

  render() {

    let steps = [
      {
        title: "Create a role",
        label: "Role creation",
        component: <StepOne external={this.props.external} next={this.nextStep} close={this.closeDialog}/>
      },{
        title: "Add your role",
        label: "Name",
        component: <StepTwo external={this.props.external} account={this.props.account} submit={this.props.submitAccount} next={this.nextStep} back={this.previousStep} close={this.closeDialog}/>
      },{
        title: "Add a bill repository",
        label: "Bill repository",
        component: <StepThree account={this.props.account} bill={this.props.bill} submit={this.props.submitBucket} close={this.closeDialog}/>
      }
    ];

    return(
      <div className="account-wizard">

        <button className="btn btn-default" onClick={this.openDialog}>Add</button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>Add an AWS account : {steps[this.state.activeStep].title}</h1></DialogTitle>

          <DialogContent>

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
  clearAccount: PropTypes.func.isRequired,
  submitBucket: PropTypes.func.isRequired,
  clearBucket: PropTypes.func.isRequired,
};

export default Wizard;
