import React from 'react';
import { Redirect } from 'react-router-dom';

const IndexRedirect = () => (
      <Redirect to={{
        pathname: '/app',
      }}/>
);

export default IndexRedirect;
