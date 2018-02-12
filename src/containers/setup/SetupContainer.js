import React, { Component } from 'react';

import Panels from '.';
import '../../styles/Setup.css';

// Setup Container for Management Panels
class SetupContainer extends Component {

  render() {
    /*
    ** Passing down the match prop so the panel
    ** can access the URL param from react-router
    */
    return (
      <div>
        <Panels.AWS.Accounts match={this.props.match}/>
      </div>
    );
  }

}

export default SetupContainer;
