import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../common/forms';
import Popover from '../../misc/Popover';
import PropTypes from 'prop-types';

const Validation = Validations.AWSAccount;

// Form Component for add or edit AWS Account
class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      roleArn: '',
      pretty: '',
      payer: false,
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.submit = this.submit.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
  }

  componentDidMount() {
    this.setState({
      roleArn: this.props.account !== undefined ? this.props.account.roleArn : undefined,
      pretty: this.props.account !== undefined ? this.props.account.pretty : undefined,
      payer: this.props.account !== undefined ? this.props.account.payer : false,
    });
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({
      open: false,
      roleArn: this.props.account !== undefined ? this.props.account.roleArn : undefined,
      pretty: this.props.account !== undefined ? this.props.account.pretty : undefined,
      payer: this.props.account !== undefined ? this.props.account.payer : false,
    });
  };

  submit = (e) => {
    e.preventDefault();
    this.closeDialog(e);
    let account = {
      roleArn: this.state.roleArn,
      pretty: this.state.pretty,
      payer: this.state.payer,
    };
    this.props.submit(account);
  };

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({
      [name]: value
    });
  }


  render() {
    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog} disabled={this.props.disabled}>
          {this.props.account !== undefined ? <i className="fa fa-edit"/> : <i className="fa fa-plus" />}
          &nbsp;
          {this.props.account !== undefined ? "Edit" : "Add"}
        </button>

        <Dialog open={this.state.open} fullWidth>

          <DialogTitle disableTypography><h1>{this.props.account !== undefined ? "Edit this" : "Create an"} account</h1></DialogTitle>

          <DialogContent>

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit} >

              <div className="form-group">
                <div className="input-title">
                  <label htmlFor="roleArn">Role ARN</label>
                  &nbsp;
                  <Popover info tooltip="Amazon Resource Name for your role"/>
                </div>
                <Input
                  name="roleArn"
                  type="text"
                  className="form-control"
                  value={this.state.roleArn}
                  onChange={this.handleInputChange}
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
                  value={this.state.pretty}
                  onChange={this.handleInputChange}
                  className="form-control"
                />
              </div>

              <div className="checkbox">
                <label>
                <input
                  type="checkbox"
                  name="payer"
                  checked={this.state.payer}
                  onChange={this.handleInputChange}
                />
                  This account has billing data
                </label>
              </div>

              <hr />

              <DialogActions>

                <button className="btn btn-default pull-left" onClick={this.closeDialog}>
                  Cancel
                </button>

                <Button
                  className="btn btn-primary btn-block"
                  type="submit"
                >
                  {this.props.account !== undefined ? "Save" : "Create"}
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
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
    payer: PropTypes.bool.isRequired,
  }),
  submit: PropTypes.func.isRequired,
  disabled: PropTypes.bool
};

FormComponent.defaultProps = {
  disabled: false
};

export default FormComponent;
