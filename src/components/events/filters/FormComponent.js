import React, { Component } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Spinner from 'react-spinkit';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Select from 'react-validation/build/select';
import Button from 'react-validation/build/button';
import Validations from '../../../common/forms';
import PropTypes from "prop-types";
import * as Filters from '../../../common/eventsFilters';

const Validation = Validations.AWSAccount;

// Form Component for new AWS Account
class FormComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false,
      name : (props.filter !== undefined ? props.filter.name : ""),
      desc : (props.filter !== undefined ? props.filter.desc : ""),
      rule : (props.filter !== undefined ? props.filter.rule : Object.keys(Filters.filters)[0]),
      data : (props.filter !== undefined ? props.filter.data : ""),
      enabled : (props.filter !== undefined ? props.filter.enabled : true),
      error: null
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.submit = this.submit.bind(this);
    this.toggleFilterDataValue = this.toggleFilterDataValue.bind(this);
    this.toggleFilterDataMultipleValues = this.toggleFilterDataMultipleValues.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
    this.addMultipleInputItem = this.addMultipleInputItem.bind(this);
    this.removeMultipleInputItem = this.removeMultipleInputItem.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    this.setState({
      open: true,
      name : (this.props.filter !== undefined ? this.props.filter.name : ""),
      desc : (this.props.filter !== undefined ? this.props.filter.desc : ""),
      rule : (this.props.filter !== undefined ? this.props.filter.rule : Object.keys(Filters.filters)[0]),
      data : (this.props.filter !== undefined ? this.props.filter.data : ""),
      disabled : (this.props.filter !== undefined ? this.props.filter.disabled : false),
      error: null
    });
    this.props.clear();
  };

  closeDialog = (e) => {
    e.preventDefault();
    this.setState({open: false, error: null});
    this.props.clear();
  };

  handleInputChange(event) {
    const misc = {};
    const target = event.target;
    const name = target.name;
    const type = target.type;
    let value = target.value;
    switch (type) {
      case 'checkbox':
        value = target.checked;
        break;
      case 'number':
        if (value === "")
          value = 0;
        else
          value = parseFloat(value);
        break;
      default:
        value = target.value;
    }
    if (name === "rule")
      misc.data = "";
    if (name === "data" && type === "select-multiple") {
      value = (Array.isArray(this.state.data) ? this.state.data : []);
      const newValue = parseInt(target.value, 10);
      const newValueIndex = value.indexOf(newValue);
      if (newValueIndex === -1)
        value.push(newValue);
      else
        value.splice(newValueIndex, 1);
    }

    this.setState({
      [name]: value,
      ...misc
    });
  }

  toggleFilterDataValue(event) {
    event.preventDefault();
    const target = event.target;
    const value = parseInt(target.value, 10);
    const data = (Array.isArray(this.state.data) ? this.state.data : []);
    const index = data.indexOf(value);
    if (index === -1)
      data.push(value);
    else
      data.splice(index, 1);
    this.setState({data});
  }

  toggleFilterDataMultipleValues(event, values) {
    event.preventDefault();
    this.setState({data: values});
  }

  handleMultipleInputData(event, index=null) {
    event.preventDefault();
    const target = event.target;
    let value = target.value;
    const data = (Array.isArray(this.state.data) ? this.state.data : []);
    if (value !== null) {
      data[index] = value;
      this.setState({data});
    }
  }

  removeMultipleInputItem(event, index=null) {
    event.preventDefault();
    const data = (Array.isArray(this.state.data) ? this.state.data : []);
    if (data.length > index)
      data.splice(index, 1);
    this.setState({data});
  }

  addMultipleInputItem(event) {
    event.preventDefault();
    const data = (Array.isArray(this.state.data) ? this.state.data : []);
    data.push("");
    this.setState({data});
  }

  checkMultipleInput() {
    const errors = [];
    const data = (Array.isArray(this.state.data) ? this.state.data : []);
    data.forEach((item, index) => {
      data.forEach((elem, idx) => {
        const msg = `${item} is duplicated`;
        if (elem === item && index !== idx && errors.indexOf(msg) === -1)
          errors.push(msg);
      })
    });
    if (errors.length)
      return errors.map((err, index) => (<div className="alert alert-warning" role="alert" key={index}>{err}</div>));
    return null;
  }

  submit = (e) => {
    e.preventDefault();
    const body = {
      name: this.state.name,
      desc: this.state.desc,
      rule: this.state.rule,
      data: this.state.data,
      disabled: this.state.disabled,
      id: (this.props.filter && this.props.filter.hasOwnProperty("id") ? this.props.filter.id : null)
    };
    const error = !Filters.checkFilterValue(this.state.rule, this.state.data);
    if (error)
      this.setState({error : "Data format is invalid."});
    else {
      this.props.submit(body);
      this.setState({error: null});
    }
  };

  componentWillReceiveProps(nextProps) {
    if (nextProps.status && nextProps.status.status && nextProps.status.values && !nextProps.status.hasOwnProperty("error")) {
      this.setState({open: false, error: null});
    }
  }

  getFilterInput() {
    if (this.state.rule) {
      const inputType = Filters.getFilterInput(this.state.rule);
      if (inputType) {
        let values;
        switch (inputType.format) {
          case "checkbox":
            if (Array.isArray(inputType.values)) {
              const days = [...inputType.values];
              const arrays = [];
              while (days.length > 0)
                arrays.push(days.splice(0, 7));
              values = arrays.map((week, index) => (
                <div key={index} className="week">
                  {week.map((value, idx) => (
                    <button
                      key={idx}
                      className={"btn btn-default " + (this.state.data.indexOf(value) !== -1 ? "active" : "")}
                      value={value}
                      onClick={this.toggleFilterDataValue}
                    >
                      {value}
                    </button>
                    ))}
                </div>
              ));
            } else {
              values = Object.keys(inputType.values).map((value, index) => (
                <button
                  key={index}
                  className={"btn btn-default " + (this.state.data.indexOf(parseInt(value, 10)) !== -1 ? "active" : "")}
                  value={value}
                  onClick={this.toggleFilterDataValue}
                >
                  {inputType.values[index]}
                </button>
              ));
            }
            return (
              <div className={"filter-btn-group " + this.state.rule}>
                <div className="filter-btn-group-items">
                  {values}
                </div>
                <div className="filter-btn-group-actions">
                  <button className="btn btn-default" onClick={(e) => this.toggleFilterDataMultipleValues(e, (Array.isArray(inputType.values) ? [...inputType.values] : Object.keys(inputType.values).map((value) => parseInt(value, 10))))}>
                    Select all
                  </button>
                  <button className="btn btn-default" onClick={(e) => this.toggleFilterDataMultipleValues(e, [])}>
                    Unselect all
                  </button>
                </div>
              </div>
            );
          case "array":
            values = [...(Array.isArray(this.state.data) ? this.state.data : [])];
            return (<div className="multiple-values">
              {values.map((item, index) => (
                <div key={index} className="multiple-values-item">
                  <input
                    name="data"
                    type={inputType.type}
                    className="form-control"
                    placeholder="Filter Value"
                    value={item}
                    onChange={(e) => {this.handleMultipleInputData(e, index)}}
                    {...inputType.customProps}
                  />
                  <button className="btn btn-default" type="button" onClick={(e) => this.removeMultipleInputItem(e, index)} disabled={values.length <= 1}>
                    <i className="fa fa-times"/>
                  </button>
                </div>
              ))}
              <button className="btn btn-default" type="button" onClick={(e) => this.addMultipleInputItem(e)}>
                <i className="fa fa-plus"/>
              </button>
              {this.checkMultipleInput()}
            </div>);
          case "input":
          default:
            return (<Input
              name="data"
              type={inputType.type}
              className="form-control"
              placeholder="Filter Value"
              value={this.state.data}
              onChange={this.handleInputChange}
              validations={[Validation.required]}
              {...inputType.customProps}
            />);
        }
      }
    }
    return (<Input
      name="data"
      type="text"
      className="form-control"
      placeholder="Filter Value"
      value={this.state.data}
      onChange={this.handleInputChange}
      validations={[Validation.required]}
    />)
  }

  render() {
    const loading = (this.props.status && !this.props.status.status ? (<Spinner className="spinner clearfix" name='circle'/>) : null);

    let error = (this.props.status && this.props.status.status && this.props.status.hasOwnProperty("error") ? (
      <div className="alert alert-warning" role="alert">{this.props.status.error.message}</div>
    ) : null);

    if (error === null)
      error = (this.state.error !== null ? (
        <div className="alert alert-warning" role="alert">{this.state.error}</div>
      ) : null);

    return (
      <div>

        <button className="btn btn-default" onClick={this.openDialog} disabled={this.props.disabled}>
          {this.props.filter !== undefined ? <i className="fa fa-edit"/> : <i className="fa fa-plus"/>}
          &nbsp;
          {this.props.filter !== undefined ? "Edit" : "Add a filter"}
        </button>

        <Dialog open={this.state.open} fullWidth maxWidth="md">

          <DialogTitle disableTypography><h1>
            <i className="fa fa-filter red-color"/>
            &nbsp;
            {this.props.filter !== undefined ? "Edit this" : "Add a"} filter
          </h1></DialogTitle>

          <DialogContent>

            {loading || error}

            <Form ref={
              /* istanbul ignore next */
              form => { this.form = form; }
            } onSubmit={this.submit}>
              <div>
                <div className="form-group">
                  <div className="input-title">
                    <label htmlFor="name">Name (optional)</label>
                  </div>
                  <Input
                    name="name"
                    type="text"
                    className="form-control"
                    placeholder="Filter Name"
                    value={this.state.name}
                    onChange={this.handleInputChange}
                    validations={[]}
                  />
                </div>
                <div className="form-group">
                  <div className="input-title">
                    <label htmlFor="desc">Description (optional)</label>
                  </div>
                  <Input
                    name="desc"
                    type="text"
                    className="form-control"
                    placeholder="Filter Description"
                    value={this.state.desc}
                    onChange={this.handleInputChange}
                    validations={[]}
                  />
                </div>
                <div className="form-group">
                  <div className="input-title">
                    <label htmlFor="rule">Type</label>
                  </div>
                  <Select
                    name="rule"
                    className="form-control"
                    value={this.state.rule}
                    onChange={this.handleInputChange}
                    validations={[Validation.required]}
                  >
                    {Object.keys(Filters.filters).map((name, index) => (<option key={index} value={name}>{Filters.filters[name].pretty}</option>))}
                  </Select>
                </div>
                <div className="form-group">
                  <div className="input-title">
                    <label htmlFor="data">Value</label>
                  </div>
                  {this.getFilterInput()}
                </div>
              </div>

              <DialogActions>

                <button className="btn btn-default pull-left" onClick={this.closeDialog}>
                  Cancel
                </button>

                <Button
                  className="btn btn-primary btn-block"
                  type="submit"
                >
                  {this.props.filter !== undefined ? "Save" : "Add"}
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
  filter: PropTypes.shape({
    name: PropTypes.string.isRequired,
    desc: PropTypes.string.isRequired,
    rule: PropTypes.string.isRequired,
    data: PropTypes.isRequired,
    disabled: PropTypes.bool.isRequired
  }),
  status: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.array
  }),
  submit: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
  disabled: PropTypes.bool
};

FormComponent.defaultProps = {
  disabled: false
};

export default FormComponent;
