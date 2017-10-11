import React, { Component } from 'react';
import Header from './common/Header';

import PageContainer from './containers/PageContainer';

class App extends Component {
  render() {
    return (
      <div>
        <Header />
        <PageContainer />
      </div>
    );
  }
}


export default App;
