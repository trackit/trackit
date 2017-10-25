import React, { Component } from 'react';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../common/forms';
import PropTypes from "prop-types";

const Validation = Validations.AWSAccount;

class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.submit = this.submit.bind(this);
  }

  submit = (e) => {
    e.preventDefault();
    let values = this.form.getValues();
    this.props.submit({roleArn: values.roleArn, external: values.external, pretty: values.pretty});
  };

  render() {
    return (
      <div className="panel panel-default">

        <div className="panel-heading">
          <h3 className="panel-title">Add an account</h3>
        </div>

        <div className="panel-body">

          <Form ref={form => {
            this.form = form;
          }} onSubmit={this.submit}>

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

            <div className="form-group">
              <label htmlFor="roleArn">Role ARN</label>
              <Input
                name="roleArn"
                type="text"
                className="form-control"
                validations={[Validation.required, Validation.roleArnFormat]}
              />
            </div>

            <div className="form-group">
              <label htmlFor="pretty">Name</label>
              <Input
                type="text"
                name="pretty"
                className="form-control"
                validations={[Validation.required]}
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
  submit: PropTypes.func.isRequired,
  external: PropTypes.string
};


export default FormComponent;