import React, {Component} from 'react';
import {connect} from 'react-redux';

import Components from '../../components';
import s3square from '../../assets/s3-square.png';

const Panel = Components.Misc.Panel;
const SingleAccountSelector = Components.AWS.Accounts.SingleAccountSelector;

// S3AnalyticsContainer Component
export class ReportsContainer extends Component {
  render() {
    return (
      <Panel>
          <div className="clearfix">
            <h3 className="white-box-title no-padding inline-block">
              <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
              AWS Reports
            </h3>
            <div className="inline-block pull-right">
              <SingleAccountSelector/>
            </div>
          </div>
      </Panel>
    );
  }
}

export default connect()(ReportsContainer);
