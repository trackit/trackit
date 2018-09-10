import React, { Component } from 'react';
import Components from '../components';

const HighLevel = Components.HighLevel.HighLevel;

// HomeContainer Component
class HomeContainer extends Component {

  render() {
    return (
      <div>
        <HighLevel/>
      </div>
    );
  }

}

export default HomeContainer;
