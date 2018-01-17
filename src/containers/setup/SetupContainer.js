import React, { Component } from 'react';

import Panels from '.';
import '../../styles/Setup.css';

// Setup Container for Management Panels
class SetupContainer extends Component {

  render() {
    return (
      <div>
        <Panels.AWS.Accounts/>
      </div>
    );
  }

}

export default SetupContainer;
