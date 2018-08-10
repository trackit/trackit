import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Validation from '../../../common/forms/AWSAccountForm';
import Selector from '../../misc/Selector';

export class AccountSelectorComponent extends Component {

  render() {
    const listedAccounts = (this.props.accounts.values && this.props.accounts.values.length > 0 ? this.props.accounts.values : null);
    const availableAccounts = {};
    if (!listedAccounts) {
      return null;
    }
    listedAccounts.forEach((account) => {
      const key = (this.props.idFromARN ? Validation.getAccountIDFromRole(account.roleArn) : account.id);
      availableAccounts[key] = account.pretty;
    });
    return(
      <Selector values={availableAccounts} selected={this.props.account} selectValue={this.props.selectAccount}/>
    );
  }

}

AccountSelectorComponent.propTypes = {
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        roleArn: PropTypes.string.isRequired,
        pretty: PropTypes.string,
        bills: PropTypes.arrayOf(
          PropTypes.shape({
            bucket: PropTypes.string.isRequired,
            path: PropTypes.string.isRequired
          })
        ),
      })
    ),
  }),
  account: PropTypes.string,
  selectAccount: PropTypes.func.isRequired,
  idFromARN: PropTypes.bool
};

AccountSelectorComponent.defaultProps = {
  idFromARN: false
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
});


export default connect(mapStateToProps)(AccountSelectorComponent);
