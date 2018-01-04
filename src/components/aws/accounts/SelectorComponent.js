import React, { Component } from 'react';
import { connect } from 'react-redux';
import List, {
  ListItem,
  ListItemText,
} from 'material-ui/List';
import PropTypes from 'prop-types';
import Actions from "../../../actions";

export class Item extends Component {

  constructor(props) {
    super(props);
    this.selectAccount = this.selectAccount.bind(this);
  }

  selectAccount = (e) => {
    e.preventDefault();
    this.props.select(this.props.account);
  };

  render() {
    return (
      <div>

        <ListItem divider>

          <ListItemText
            disableTypography
            primary={this.props.account.pretty || this.props.account.roleArn}
          />

          <button className={"btn btn-" + (this.props.isSelected ? "danger" : "default")} onClick={this.selectAccount}>
            {this.props.isSelected ? "Unselect" : "Select"}
          </button>

        </ListItem>

      </div>
    );
  }

}

Item.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
  }),
  select: PropTypes.func.isRequired,
  isSelected: PropTypes.bool
};

// Selector Component for AWS Accounts
class SelectorComponent extends Component {

  componentWillMount() {
    this.props.getAccounts();
  }

  render() {

    const isSelected = (item) => (this.props.selected.find((value) => (value.id === item.id)) !== undefined);

    let noAccounts = (!this.props.accounts || !this.props.accounts.length ? <div className="alert alert-warning" role="alert">No account available</div> : "");
    let accounts = (this.props.accounts && this.props.accounts.length ? (
      this.props.accounts.map((account, index) => (
        <Item
          key={index}
          account={account}
          select={this.props.select}
          isSelected={isSelected(account)}
        />
      ))
    ) : null);
    return (
      <div id="account-selection">
        <List disablePadding>
          {noAccounts}
          {accounts}
        </List>
      </div>
    );
  }

}

SelectorComponent.propTypes = {
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string,
    })
  ),
  selected: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string,
    })
  ),
  select: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  selected: aws.accounts.selection
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getAccounts: () => {
    dispatch(Actions.AWS.Accounts.getAccounts());
  },
  select: (account) => {
    dispatch(Actions.AWS.Accounts.selectAccount(account));
  },
  clear: () => {
    dispatch(Actions.AWS.Accounts.clearAccountSelection());
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(SelectorComponent);
