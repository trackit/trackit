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
import Reports_first from '../../../assets/report_step_1.png';
import Reports_second from '../../../assets/report_step_2.png';
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
    this.props.submit(this.props.account.value.id, formValues);
  };


  componentWillReceiveProps(nextProps) {
    if (nextProps.bill.status && nextProps.bill.value)
      nextProps.close();
  }

  render() {
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

    const error = (this.props.bill && this.props.bill.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">{this.props.bill.error.message}</div>
    ) : null);


    return (
      <div>

            {tutorial}
            {error}

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="bucket">S3 Bucket name</label>
                  &nbsp;
                  <Popover info popOver="Name of the S3 bucket you created"/>
                </div>
                <Input
                  name="bucket"
                  type="text"
                  className="form-control"
                  placeholder="Bucket Name"
                  value={""}
                  validations={[Validation.required, Validation.s3BucketNameFormat]}
                />
              </div>

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="bucket">Report path prefix (optional)</label>
                  &nbsp;
                  <Popover info popOver="If you set a path prefix when creating your report"/>
                </div>
                <Input
                  name="prefix"
                  type="text"
                  className="form-control"
                  placeholder="Optional prefix"
                  value={""}
                  validations={[Validation.s3PrefixFormat]}
                />
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
