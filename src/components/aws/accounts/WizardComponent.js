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
      <div>

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
                using the informations provided
                <Picture
                  src={RoleCreation}
                  alt="Role creation tutorial"
                  button={<strong>( Click here to see screenshot )</strong>}
                />
              </div>
            </li>
            <li>Select <strong>ReadOnlyAccess</strong> policy</li>
            <li>Set a name to this new role and validate</li>
          </ol>

        </div>

        <Form ref={
          /* istanbul ignore next */
          form => { this.form = form; }
        } onSubmit={this.submit} >

          <div className="form-group">
            <div className="input-title">
              <label htmlFor="externalId">Account ID</label>
              &nbsp;
              <Popover info popOver="Account ID to add in your IAM role trust policy ( See step 5 )"/>
            </div>
            <Input
              type="text"
              name="accountID"
              className="form-control"
              disabled
              value={this.props.external.accountId}
              validations={[Validation.required]}
            />
          </div>

          <div className="form-group">
            <div className="input-title">
              <label htmlFor="externalId">External</label>
              &nbsp;
              <Popover info popOver="External ID to add in your IAM role trust policy ( See step 5 )"/>
            </div>
            <Input
              type="text"
              name="external"
              className="form-control"
              disabled
              value={this.props.external.external}
              validations={[Validation.required]}
            />
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
    this.props.next();
  };

  render() {
    return (
      <div>

        <div className="tutorial">

          <ol>
            <li>In <strong>Role</strong> list, select the role you created in previous step</li>
            <li>
              <div>
                Fill the form below with information available in <strong>role summary</strong>.
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
            <button className="btn btn-default col-md-5 btn-left" onClick={this.props.close}>Cancel</button>
            <Button className="btn btn-primary col-md-5 btn-right" type="submit">Next</Button>
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
  submit: PropTypes.func.isRequired,
  next: PropTypes.func.isRequired,
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
    this.props.submit(this.props.account.id, bill);
    this.props.close(e);
  };

  render() {
    return (
      <div>

        <div>

          <div className="tutorial">

            <ol>
              <li>Fill the form with the location of a <strong>S3 bucket</strong> that contains bills</li>
              <li>You will be able to add more buckets later.</li>
            </ol>

          </div>

        </div>

        <Form ref={
          /* istanbul ignore next */
          form => { this.form = form; }
        } onSubmit={this.submit} >

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
              validations={[Validation.required, Validation.s3BucketFormat]}
            />
          </div>

          <div className="form-group clearfix">
            <button className="btn btn-default col-md-5 btn-left" onClick={this.props.close}>Cancel</button>
            <Button className="btn btn-primary col-md-5 btn-right" type="submit" disabled={!this.props.account}>Done</Button>
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
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
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
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
  }

  nextStep = () => {
    const activeStep = this.state.activeStep + 1;
    this.setState({activeStep});
  };

  openDialog = (e) => {
    e.preventDefault();
    this.setState({open: true, activeStep: 0});
    this.props.clearAccount();
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false, activeStep: 0});
    this.props.clearAccount();
  };

  render() {

    let steps = [
      {
        label: "Role creation",
        component: <StepOne external={this.props.external} next={this.nextStep} close={this.closeDialog}/>
      },{
        label: "Name",
        component: <StepTwo external={this.props.external} submit={this.props.submitAccount} next={this.nextStep} close={this.closeDialog}/>
      },{
        label: "Bill repository",
        component: <StepThree account={this.props.account} submit={this.props.submitBucket} close={this.closeDialog}/>
      }
    ];

    return(
      <div className="account-wizard">

        <button className="btn btn-default" onClick={this.openDialog}>Add</button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>Create an account</h1></DialogTitle>

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
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
  }),
  submitAccount: PropTypes.func.isRequired,
  clearAccount: PropTypes.func.isRequired,
  submitBucket: PropTypes.func.isRequired,
};


export default Wizard;
