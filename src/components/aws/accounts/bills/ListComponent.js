import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Form from './FormComponent';

class ListItem extends Component {

  constructor(props) {
    super(props);
    this.state = {
      editForm: false
    };
    this.toggleEditForm = this.toggleEditForm.bind(this);
    this.editBill = this.editBill.bind(this);
    this.deleteBill = this.deleteBill.bind(this);
  }

  toggleEditForm = (e) => {
    e.preventDefault();
    const editForm = !this.state.editForm;
    this.setState({ editForm });
  };

  editBill = (body) => {
    console.log(body);
    this.setState({ editForm: false });
  };

  deleteBill = (e) => {
    e.preventDefault();
    this.props.delete(this.props.bill.id);
  };

  render() {

    const actions = (!this.state.editForm ? (
      <div className="btn-group btn-group-xs pull-right" role="group">

        <button type="button" className="btn btn-default" onClick={this.toggleEditForm}>
          <span className="glyphicon glyphicon-pencil" aria-hidden="true"/> Edit
        </button>

        <button type="button" className="btn btn-default" onClick={this.deleteAccount}>
          <span className="glyphicon glyphicon-remove" aria-hidden="true"/> Delete
        </button>

      </div>
    ) : (
      <div className="btn-group btn-group-xs pull-right" role="group">
        <button type="button" className="btn btn-default" onClick={this.toggleEditForm}>
          <span className="glyphicon glyphicon-remove" aria-hidden="true"/> Hide
        </button>
      </div>
    ));

    const editForm = (this.state.editForm ? (
      <Form submit={this.editBill} bill={this.props.bill} account={this.props.account} />
    ) : null);

    return (
      <li className="bill list-group-item">

        <span className="bill-bucket">
          {this.props.bill.bucket}
        </span>

        <span className="bill-path">
          {this.props.bill.path}
        </span>

        {actions}

        {editForm}

      </li>
    );
  }

}

ListItem.propTypes = {
  account: PropTypes.number.isRequired,
  bill: PropTypes.shape({
    bucket: PropTypes.string.isRequired,
    path: PropTypes.string.isRequired
  }),
  edit: PropTypes.func.isRequired,
  delete: PropTypes.func.isRequired
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    let noBills = (!this.props.bills.length ? <div className="alert alert-warning" role="alert">No bills available</div> : "");
    return (
      <ul className="bills list-group list-group-flush">
        {noBills}
        {this.props.bills.map((bill, index) => (
          <ListItem
            key={index}
            bill={bill}
            account={this.props.account}
            edit={this.props.edit}
            delete={this.props.delete}/>
        ))}
      </ul>
    );
  }

}

ListComponent.propTypes = {
  account: PropTypes.number.isRequired,
  bills: PropTypes.arrayOf(
    PropTypes.shape({
      bucket: PropTypes.string.isRequired,
      path: PropTypes.string.isRequired
    })
  ),
  new: PropTypes.func.isRequired,
  edit: PropTypes.func.isRequired,
  delete: PropTypes.func.isRequired
};

export default ListComponent;