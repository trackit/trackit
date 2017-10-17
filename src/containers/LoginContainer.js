import React, { Component } from 'react';
import { connect } from 'react-redux';
// import PropTypes from 'prop-types';

import Actions from '../actions';

// LoginContainer Component
class LoginContainer extends Component {

  componentDidMount() {}

  submit() {
    this.props.login('testUsername', 'testPassword');
  }

  render() {
    return (
      <div className="container-fluid">
        <h2>LOGIN VIEW</h2>
        <br/>
        <button className="btn btn-default" onClick={this.submit.bind(this)}>LOGIN</button>
      </div>
    );
  }

}

LoginContainer.propTypes = {};

const mapStateToProps = () => ({

});

const mapDispatchToProps = (dispatch) => ({
  login: (username, password) => {
    dispatch(Actions.Auth.Login.login(username, password))
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(LoginContainer);
