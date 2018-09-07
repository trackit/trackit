import React, { Component } from 'react';
import Components from '../components';

const Dashboard = Components.Dashboard.Dashboard;

// DashboardContainer Component
class DashboardContainer extends Component {

  render() {
    return (
      <div>
        <Dashboard/>
      </div>
    );
  }

}

export default DashboardContainer;
