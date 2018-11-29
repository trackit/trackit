import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import List from '@material-ui/core/List';
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import Collapse from "@material-ui/core/Collapse/Collapse";
import Checkbox from "@material-ui/core/Checkbox/Checkbox";
import Status from "../../../../common/awsAccountStatus";

class Item extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
    this.handleClick = this.handleClick.bind(this);
    this.selectAccount = this.selectAccount.bind(this);
  }

  handleClick = () => {
    this.setState(state => ({ open: !state.open }));
  };

  selectAccount = (account) => {
    this.props.select(account);
  };

  render() {
    const status = Status.getAWSAccountStatus(this.props.account);
    const accountBadge = Status.getBadge(status);

    const subaccounts = (this.props.account.hasOwnProperty("subAccounts") ? (
      <Collapse in={this.state.open} timeout="auto" unmountOnExit>
        <List disablePadding className="account-item-details">
          {this.props.account.subAccounts.map((subAccount) => {
            const status = Status.getAWSAccountStatus(subAccount);
            const badge = Status.getBadge(status);
            return (
              <ListItem>
                <Checkbox
                  className={"checkbox " + (this.props.isSelected(subAccount) ? "selected" : "")}
                  checked={this.props.isSelected(subAccount)}
                  onChange={(e) => {e.preventDefault(); this.selectAccount(subAccount)}}
                  disableRipple
                />
                <ListItemText inset primary={subAccount.pretty || subAccount.awsIdentity} />
                <div className="actions">
                  {badge}
                </div>
              </ListItem>
            );
          })}
        </List>
      </Collapse>
    ) : null);

    const prefix = (this.props.account.subAccounts && this.props.account.subAccounts.length) ? (
      <span className="badge blue-bg pull-right">{this.props.account.subAccounts.length} sub accounts</span>
    ) : (null);

    return (
      <div>
        <ListItem className="account-item">
          <Checkbox
            className={"checkbox " + (this.props.isSelected(this.props.account) ? "selected" : "")}
            checked={this.props.isSelected(this.props.account)}
            onChange={(e) => {e.preventDefault(); this.selectAccount(this.props.account)}}
            disableRipple
          />
          <ListItemText
            disableTypography
            className="account-name"
            primary={this.props.account.pretty || this.props.account.awsIdentity}
          />
          <div className="actions">
            {prefix}
            &nbsp;
            {accountBadge}
          </div>
          {this.props.account.hasOwnProperty("subAccounts") ? (this.state.open ? <ExpandLess onClick={this.handleClick}/> : <ExpandMore onClick={this.handleClick}/>) : null}
        </ListItem>
        {subaccounts}
      </div>
    );
  }

}

Item.propTypes = {
  account: PropTypes.shape({
    id: PropTypes.number.isRequired,
    accountOwner: PropTypes.bool.isRequired,
    awsIdentity: PropTypes.string.isRequired,
    roleArn: PropTypes.string.isRequired,
    pretty: PropTypes.string,
    permissionLevel: PropTypes.number,
    payer: PropTypes.bool.isRequired,
    billRepositories: PropTypes.arrayOf(
      PropTypes.shape({
        error: PropTypes.string.isRequired,
        nextPending: PropTypes.bool.isRequired,
        bucket: PropTypes.string.isRequired,
        prefix: PropTypes.string.isRequired
      })
    ),
    subAccounts: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        accountOwner: PropTypes.bool.isRequired,
        awsIdentity: PropTypes.string.isRequired,
        roleArn: PropTypes.string.isRequired,
        pretty: PropTypes.string,
        permissionLevel: PropTypes.number,
        payer: PropTypes.bool.isRequired,
        billRepositories: PropTypes.arrayOf(
          PropTypes.shape({
            error: PropTypes.string.isRequired,
            nextPending: PropTypes.bool.isRequired,
            bucket: PropTypes.string.isRequired,
            prefix: PropTypes.string.isRequired
          })
        ),
      })
    )
  }),
  select: PropTypes.func.isRequired,
  isSelected: PropTypes.func.isRequired
};

export default Item;
