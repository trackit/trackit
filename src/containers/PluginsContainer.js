import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import Components from '../components';
import Actions from '../actions';
import Spinner from "react-spinkit";

import "../styles/Plugins.css";

// PluginsContainer Component
class PluginsContainer extends Component {
  componentDidMount() {
    this.props.getData();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.accounts !== nextProps.accounts)
      nextProps.getData();
  }

  accountExists(values, account) {
    for (let i = 0; i < values.length; i++) {
      const element = values[i];
      if (element.account === account)
        return i;
    }
    return null;
  }

  getAccountLabel(account) {
    if (this.props.allAccounts && this.props.allAccounts.status && this.props.allAccounts.values) {
      const all = this.props.allAccounts.values;
      for (let i = 0; i < all.length; i++) {
        const element = all[i];
        if (element.roleArn.split(':')[4] === account) {
          return element.pretty;
        }
      }
    }
    return null;
  }

  groupResultsByAccounts(values) {
    const res = [];
    for (let i = 0; i < values.length; i++) {
      const element = values[i];
      const existsIndex = this.accountExists(res, element.account);
      if (existsIndex != null && res[existsIndex] && res[existsIndex].results) {
        res[existsIndex].results.push(element);
      } else {
        const label = this.getAccountLabel.bind(this)(element.account);
        res.push({
          account: element.account,
          label,
          results: [element]
        });
      }
    }
    return res;
  }

  render() {

    const loading = (!this.props.data.status ? (<Spinner className="spinner" name='circle'/>) : null);

    const error = (this.props.data.error ? ` (${this.props.data.error.message})` : null);
    const noData = (this.props.data.status && (!this.props.data.values || !this.props.data.values.length || error) ? <div className="alert alert-warning" role="alert">No account available{error}</div> : "");

    const spinnerAndError = (loading || noData ? (
      <div className="white-box">
        {loading}
        {noData}
      </div>
    ) : null);

    const accounts = [];

    if (this.props.data && this.props.data.status && this.props.data.values) {
      const formattedData = this.groupResultsByAccounts(this.props.data.values);
      for (let i = 0; i < formattedData.length; i++) {
        const element = formattedData[i];
        accounts.push(<Components.Plugins.Account account={element} key={element.account} />);
      }
    }

    return (
      <div>
        <div className="row">
          <div className="col-md-12">
            <div className="white-box">
              <h3 className="white-box-title no-padding inline-block">
                <i className="fa fa-check-square-o"></i>
                &nbsp;
                Recommandations
              </h3>
            </div>
            {spinnerAndError}
            {accounts}
          </div>
        </div>
      </div>
    );
  }

}

PluginsContainer.propTypes = {
  data: PropTypes.object.isRequired,
  accounts: PropTypes.arrayOf(PropTypes.object),
  allAccounts: PropTypes.object,
  getData: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({plugins, aws}) => ({
  data: plugins.data,
  accounts: aws.accounts.selection,
  allAccounts: aws.accounts.all,
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: () => {
    dispatch(Actions.Plugins.getData());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(PluginsContainer);
