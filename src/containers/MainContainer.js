import React, { Component } from 'react';
import Components from '../components';

// MainContainer Component
class MainContainer extends Component {

  render() {

    const styles = {
      content: {
        marginLeft: 60,
        padding: 25
      }
    };

    return (
      <div>
        <Components.Misc.Navigation/>
        <div className="content" style={styles.content}>
          {this.props.children}
        </div>
      </div>
    );
  }

}

export default MainContainer;
