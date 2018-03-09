import React, { Component } from 'react';
import Components from '../components';

const Dashboard = Components.Dashboard.Dashboard;

// HomeContainer Component
class HomeContainer extends Component {

  render() {
    return (
      <div>
        <Dashboard/>
      </div>
    );
  }

}

export default HomeContainer;
