import React, { Component } from 'react';

import Containers from './containers';

class App extends Component {
  render() {
    return (
      <div>
        <Containers.Main>
          <Containers.Page />
        </Containers.Main>
      </div>
    );
  }
}

export default App;
