import React, { Component } from 'react';
import { connect } from 'react-redux';

import Components from '../../../components';
import Actions from "../../../actions";

const List = Components.AWS.Access.List;

// MainContainer Component
class AccessManagementContainer extends Component {

  componentWillMount() {
    this.props.getAccess();
  }

  render() {
    return (
      <div>
        <List/>
        {this.props.access | "None"}
      </div>
    );
  }

}

const mapStateToProps = (state) => ({access: state.aws.access});

const mapDispatchToProps = (dispatch) => ({
  getAccess: () => {
    dispatch(Actions.AWS.Access.getAccess())
  }
});

export default connect(mapStateToProps, mapDispatchToProps)(AccessManagementContainer);
