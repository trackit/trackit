import React, { Component } from 'react';
import AWS from './aws';

const CostBreakdown = AWS.CostBreakdown;

// HomeContainer Component
class HomeContainer extends Component {

  render() {
    return (
      <div>
        <CostBreakdown/>
      </div>
    );
  }

}

export default HomeContainer;
