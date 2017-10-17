import React from 'react';
import { Route, Redirect } from 'react-router-dom'


const PrivateRoute = ({ component: Component, store, ...rest }) => {

  const state = store.getState();

  return(<Route {...rest} render={props => (
    state.auth.token ? (
      <Component {...props}/>
    ) : (
      <Redirect to={{
        pathname: '/login',
        state: { from: props.location }
      }}/>
    )
  )}/>);
}

export default PrivateRoute;
