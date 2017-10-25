import React, { Component } from 'react';
import PropTypes from 'prop-types';

class ListItem extends Component {

  render() {
    return (
      <div className="account list-group-item">
        {this.props.account.pretty}
      </div>
    );
  }

}

ListItem.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    userId: PropTypes.number.isRequired,
    pretty: PropTypes.string.isRequired
  })
};

// List Component for AWS Accounts
class ListComponent extends Component {

  render() {
    let noAccounts = (!this.props.accounts.length ? <div className="alert alert-warning" role="alert">No account available</div> : "");
    return (
      <div className="accounts list-group list-group-flush">
        {noAccounts}
        {this.props.accounts.map((account) => (<ListItem key={account.id} account={account}/>))}
      </div>
    );
  }

}

ListComponent.propTypes = {
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      userId: PropTypes.number.isRequired,
      pretty: PropTypes.string.isRequired
    })
  )
};

export default ListComponent;