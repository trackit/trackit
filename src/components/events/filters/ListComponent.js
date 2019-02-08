import React, { Component } from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Checkbox from "@material-ui/core/Checkbox/Checkbox";
import Spinner from 'react-spinkit';
import PropTypes from "prop-types";
import Misc from '../../misc';
import * as Filters from '../../../common/eventsFilters';
import Form from './FormComponent';

const Dialog = Misc.Dialog;
const Popover = Misc.Popover;
const Delete = Misc.DeleteConfirmation;

class Item extends Component {

  constructor(props) {
    super(props);
    this.switchFilterStatus = this.switchFilterStatus.bind(this);
  }

  switchFilterStatus = (e) => {
    e.preventDefault();
    const newBody = this.props.filter;
    newBody.disabled = !this.props.filter.disabled;
    this.props.set(newBody);
  };

  render() {
    let desc = null;

    if (this.props.filter.hasOwnProperty("desc") && this.props.filter.desc.length)
      desc = (<Popover info tooltip={this.props.filter.desc}/>);

    return (
      <ListItem divider className="filter-list-item">
        <div className="item-info">
          <Checkbox
            className={"checkbox " + (!this.props.filter.disabled ? "selected" : "")}
            checked={!this.props.filter.disabled}
            onChange={this.switchFilterStatus}
            disableRipple
          />
          <div className="item-details">
            <div className="item-name">
              {this.props.filter.name}
              &nbsp;
              {desc}
            </div>
            <div className="item-filter">
              {Filters.showFilter(this.props.filter)}
            </div>
          </div>
        </div>
        <div className="item-actions">
          <Form
            filter={this.props.filter}
            status={this.props.status}
            submit={this.props.set}
            clear={this.props.clear}
          />
          <Delete entity={`filter named ${this.props.filter.name}`} confirm={() => this.props.delete(this.props.filter.id)}/>
        </div>
      </ListItem>
    );
  }

}

Item.propTypes = {
  filter: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    desc: PropTypes.string.isRequired,
    rule: PropTypes.string.isRequired,
    data: PropTypes.isRequired,
    disabled: PropTypes.bool.isRequired
  }),
  status: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    value: PropTypes.object
  }),
  set: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
  delete: PropTypes.func.isRequired,
};

// List Component for Anomalies Filters
class ListComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
    this.setFilters = this.setFilters.bind(this);
    this.deleteFilter = this.deleteFilter.bind(this);
  }

  componentWillMount() {
    this.props.actions.get();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.filterEdition.status !== nextProps.filterEdition.status && nextProps.filterEdition.status)
      this.props.actions.get();
  }

  componentWillUnmount() {
    this.props.actions.clear();
  }

  setFilters(newFilter) {
    const filters = [...this.props.filters.values];
    if (newFilter.hasOwnProperty("id") && newFilter.id !== null)
      filters[newFilter.id] = newFilter;
    else
      filters.push(newFilter);
    this.props.actions.set(filters);
  }

  deleteFilter(id) {
    const filters = this.props.filters.values;
    filters.splice(id, 1);
    this.props.actions.set(filters);
  }

  render() {

    const loading = (!this.props.filters.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.filters.error ? ` (${this.props.filters.error.message})` : null);
    const noFilters = (this.props.filters.status && (!this.props.filters.values || !this.props.filters.values.length || error) ? <div className="alert alert-warning" role="alert">No filters available{error}</div> : "");

    const values = (this.props.filters.status && this.props.filters.hasOwnProperty("values") && this.props.filters.values ? (
      this.props.filters.values.map((item, index) => (
        <Item
          key={index}
          filter={item}
          status={this.props.filterEdition}
          set={this.setFilters}
          clear={this.props.actions.clearSet}
          delete={this.deleteFilter}
        />))
    ) : null);


    let enabledFilters;
    if (this.props.filters.status && this.props.filters.hasOwnProperty("values") && this.props.filters.values) {
      const count = this.props.filters.values.filter((filter) => (!filter.disabled)).length;
      if (count)
        enabledFilters = <span className="filters-badge badge">{count}</span>;
    }

    return (
      <Dialog
        buttonName={<span><i className="fa fa-filter"/> Filters {enabledFilters}</span>}
        disabled={this.props.disabled}
        title="Events Filters"
        secondActionName="Close"
        onOpen={this.getFilters}
        onClose={this.clearFilters}
        titleChildren={<Form
          status={this.props.filterEdition}
          submit={this.setFilters}
          clear={this.props.actions.clearSet}
        />}
      >

        <List className="filters-list">
          {loading}
          {noFilters}
          {values}
        </List>

      </Dialog>
    );
  }

}

ListComponent.propTypes = {
  filters: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        name: PropTypes.string.isRequired,
        desc: PropTypes.string.isRequired,
        rule: PropTypes.string.isRequired,
        data: PropTypes.isRequired,
        disabled: PropTypes.bool.isRequired
      })
    )
  }),
  filterEdition: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.array
  }),
  actions: PropTypes.shape({
    get: PropTypes.func.isRequired,
    clear: PropTypes.func.isRequired,
    set: PropTypes.func.isRequired,
    clearSet: PropTypes.func.isRequired,
  }).isRequired,
};

export default ListComponent;
