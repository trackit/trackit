import React, { Component } from 'react';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../common/forms';
import PropTypes from "prop-types";

const Validation = Validations.AWSAccount;

// Form Component for new AWS Account
class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.submit = this.submit.bind(this);
  }

  submit = (e) => {
    e.preventDefault();
    let values = this.form.getValues();
    let account = {
      roleArn: values.roleArn,
      pretty: values.pretty
    };
    if (this.props.account === undefined && this.props.external)
      account.external = values.external;
    this.props.submit(account);
  };

  render() {
    const actionTitle = (this.props.account !== undefined ? "Edit" : "Add") + " an account";
    const external = (this.props.account !== undefined ? "" : (
      <div className="form-group">
        <label htmlFor="externalId">External</label>
        <Input
          type="text"
          name="external"
          className="form-control"
          disabled
          value={this.props.external}
          validations={[Validation.required]}
        />
      </div>
    ));
    return (
      <div className="panel panel-default">

        <div className="panel-heading">
          <h3 className="panel-title">{actionTitle}</h3>
        </div>

        <div className="panel-body">

          <Form ref={form => {
            this.form = form;
          }} onSubmit={this.submit}>

            {external}

            <div className="form-group">
              <label htmlFor="roleArn">Role ARN</label>
              <Input
                name="roleArn"
                type="text"
                className="form-control"
                value={(this.props.account !== undefined ? this.props.account.roleArn : undefined)}
                validations={[Validation.required, Validation.roleArnFormat]}
              />
            </div>

            <div className="form-group">
              <label htmlFor="pretty">Name</label>
              <Input
                type="text"
                name="pretty"
                value={(this.props.account !== undefined ? this.props.account.pretty : undefined)}
                className="form-control"
              />
            </div>

            <div>
              <Button
                className="btn btn-primary btn-block"
                type="submit"
              >
                <i className="fa fa-plus" />
                &nbsp;
                Add
              </Button>
            </div>

          </Form>

        </div>
      </div>
    );
  }

}

FormComponent.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    userId: PropTypes.number.isRequired,
    pretty: PropTypes.string
  }),
  submit: PropTypes.func.isRequired,
  external: PropTypes.string
};


export default FormComponent;