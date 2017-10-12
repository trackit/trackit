import React, { Component } from 'react';
import Components from '../components';

// MainContainer Component
class MainContainer extends Component {

  render() {
    return (
      <div>
        <Components.Misc.Header />
        {this.props.children}
      </div>
    );
  }

}

export default MainContainer;
