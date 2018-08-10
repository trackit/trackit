import React, {Component} from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Actions from "../../../actions";
import Spinner from "react-spinkit";
import Moment from "moment";
import Misc from '../../misc';

const Popover = Misc.Popover;

class DatabasesComponent extends Component {

  componentWillMount() {
    if (this.props.account)
      this.props.getData(this.props.account);
  }

  componentWillReceiveProps(nextProps) {
    if (!nextProps.account)
      nextProps.clear();
    else if (nextProps.account !== this.props.account)
      nextProps.getData(nextProps.account);
  }

  render() {
    const loading = (!this.props.data.status ? (<Spinner className="spinner" name='circle'/>) : null);
    const error = (this.props.data.error ? ` (${this.props.data.error.message})` : null);

    const reportDate = (this.props.data.status && this.props.data.hasOwnProperty("value") && this.props.data.value ? (
      <Popover info tooltip={"Report created " + Moment(this.props.data.value.reportDate).fromNow()}/>
    ) : null);

    return (
      <div className="clearfix resources">
        <h3 className="white-box-title no-padding inline-block">
          <i className="menu-icon fa fa-database"/>
          &nbsp;
          Databases
          {reportDate}
        </h3>
        {loading}
        {error}
        {JSON.stringify(this.props.data)}
      </div>
    )
  }

}

DatabasesComponent.propTypes = {
  account: PropTypes.string,
  data: PropTypes.object,
  getData: PropTypes.func.isRequired,
  clear: PropTypes.func.isRequired,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  account: aws.resources.account,
  data: aws.resources.RDS
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  getData: (accountId) => {
    dispatch(Actions.AWS.Resources.get.RDS(accountId));
  },
  clear: () => {
    dispatch(Actions.AWS.Resources.clear.RDS());
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(DatabasesComponent);
