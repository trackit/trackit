import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions';
import s3square from '../../assets/s3-square.png';

const Panel = Components.Misc.Panel;
const AccountSelector = Components.AWS.Accounts.AccountSelector;
const VMs = Components.AWS.Resources.VMs;
const Databases = Components.AWS.Resources.Databases;
const Storage = Components.AWS.Resources.Storage;

export class ResourcesContainer extends Component {
  render() {
    return(
      <Panel>
        <div className="clearfix">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            Resources
          </h3>
          <div className="inline-block pull-right">
            Selected account :
            <AccountSelector
              account={this.props.account}
              selectAccount={this.props.selectAccount}
              idFromARN={true}
            />
          </div>
        </div>
        <VMs/>
      </Panel>
    );
  }
}

ResourcesContainer.propTypes = {
  account: PropTypes.string,
  selectAccount: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  account: aws.resources.account,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  selectAccount: (accountId) => {
    dispatch(Actions.AWS.Resources.selectAccount(accountId));
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesContainer);
