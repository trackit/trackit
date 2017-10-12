import React, { Component } from 'react';
import { connect } from 'react-redux';

// DeclareSetupComponent Component
class DeclareSetupComponent extends Component {

    render() {
      return (
        <div>
          Bonjour
        </div>
      );
    }

}

// connect method from react-router connects the component with redux store
export default connect()(DeclareSetupComponent);
