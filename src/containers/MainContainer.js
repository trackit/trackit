import React, { Component } from 'react';
import Components from '../components';

// MainContainer Component
class MainContainer extends Component {

  render() {
    return (
      <div>
        <Components.Misc.Navigation/>
        <div className="content">
          {this.props.children}
        </div>
      </div>
    );
  }

}

export default MainContainer;
