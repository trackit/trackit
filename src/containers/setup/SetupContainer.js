import React, { Component } from 'react';

import Panels from '.';

// SetupContainer Component
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
