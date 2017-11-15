import React, { Component } from 'react';

// Form imports
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../../common/forms';
import PropTypes from "prop-types";

import Misc from '../../../misc';

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
    let bill = {
      bucket: values.bucket,
      path: values.path
    };
    this.props.submit(this.props.account, bill);
  };

  render() {
    const actionVerb = (this.props.bill !== undefined ? "Edit" : "Add");

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
      <Panel title={actionVerb + " a bill location"} collapsible defaultCollapse={!this.props.bill}>
        <Form ref={form => { this.form = form; }} onSubmit={this.submit}>

          <div className="form-group">
            <label htmlFor="bucket">S3 Bucket</label>
            <Input
              name="bucket"
              type="text"
              className="form-control"
              value={(this.props.bill !== undefined ? this.props.bill.bucket : "s3://")}
              validations={[Validation.required, Validation.s3BucketFormat]}
            />
          </div>

          <div className="form-group">
            <label htmlFor="path">Path</label>
            <Input
              type="text"
              name="path"
              value={(this.props.bill !== undefined ? this.props.bill.path : undefined)}
              className="form-control"
              validations={[Validation.required, Validation.pathFormat]}
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
  account: PropTypes.number.isRequired,
  bill: PropTypes.shape({
    bucket: PropTypes.string.isRequired,
    path: PropTypes.string.isRequired
  }),
  submit: PropTypes.func.isRequired
};


export default FormComponent;