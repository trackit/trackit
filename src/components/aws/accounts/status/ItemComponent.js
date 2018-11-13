import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ListItem from "@material-ui/core/ListItem/ListItem";
import ListItemText from "@material-ui/core/ListItemText/ListItemText";
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

  selectAccount = (e) => {
    e.preventDefault();
    this.props.select(this.props.account);
  };

  render() {
    const status = Status.getAWSAccountStatus(this.props.account);
    const accountBadge = Status.getBadge(status);
/*
    const subaccounts = (this.props.account.hasOwnProperty("children") ? (
      <Collapse in={this.state.open} timeout="auto" unmountOnExit>
        <List disablePadding className="account-item-details">
          <ListItem>
            <Checkbox
              className={"checkbox " + (this.props.isSelected ? "selected" : "")}
              checked={this.props.isSelected}
              onChange={this.selectAccount}
              disableRipple
            />
            <ListItemText inset primary="Subaccount 1" />
          </ListItem>
          <ListItem>
            <Checkbox
              className={"checkbox " + (this.props.isSelected ? "selected" : "")}
              checked={this.props.isSelected}
              onChange={this.selectAccount}
              disableRipple
            />
            <ListItemText inset primary="Subaccount 2" />
          </ListItem>
        </List>
      </Collapse>
    ) : null);
*/
    return (
      <div>
        <ListItem className="account-item">
          <Checkbox
            className={"checkbox " + (this.props.isSelected ? "selected" : "")}
            checked={this.props.isSelected}
            onChange={this.selectAccount}
            disableRipple
          />
          <ListItemText
            disableTypography
            className="account-name"
            primary={this.props.account.pretty || this.props.account.roleArn}
          />
          <div className="actions">
            {accountBadge}
          </div>
          {/*this.props.account.hasOwnProperty("children") ? (this.state.open ? <ExpandLess onClick={this.handleClick}/> : <ExpandMore onClick={this.handleClick}/>) : null*/}
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
  }),
  select: PropTypes.func.isRequired,
  isSelected: PropTypes.bool
};

export default Item;
