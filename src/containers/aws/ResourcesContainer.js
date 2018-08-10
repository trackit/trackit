import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Components from '../../components';
import Actions from '../../actions';
import s3square from '../../assets/s3-square.png';

const Panel = Components.Misc.Panel;
const VMs = Components.AWS.Resources.VMs;
const Databases = Components.AWS.Resources.Databases;

export class ResourcesContainer extends Component {
  render() {
    return(
      <Panel>
        <div className="clearfix">
          <h3 className="white-box-title no-padding inline-block">
            <img className="white-box-title-icon" src={s3square} alt="AWS square logo"/>
            Resources
          </h3>
        </div>
        <VMs/>
        <Databases/>
      </Panel>
    );
  }
}

export default ResourcesContainer;
