import React, { Component } from 'react';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../common/forms';
import PropTypes from 'prop-types';

import Misc from '../../misc';

const Panel = Misc.Panel;

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
    const actionVerb = (this.props.account !== undefined ? "Edit" : "Add");

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

    const button = (this.props.account !== undefined ? (
      <div>
        <span className="glyphicon glyphicon-pencil" aria-hidden="true"/>&nbsp;Save
      </div>
    ) : (
      <div>
        <i className="fa fa-plus" />&nbsp;Add
      </div>
    ));

    return (
      <Panel title={actionVerb + " an account"} collapsible defaultCollapse>
        <Form
          ref={
            /* istanbul ignore next */
            (form) => {this.form = form;}
          }
          onSubmit={this.submit}>

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
              {button}
            </Button>
          </div>

        </Form>
      </Panel>
    );
  }

}

FormComponent.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
    bills: PropTypes.arrayOf(
      PropTypes.shape({
        bucket: PropTypes.string.isRequired,
        path: PropTypes.string.isRequired
      })
    ),
  }),
  submit: PropTypes.func.isRequired,
  external: PropTypes.string
};


export default FormComponent;
