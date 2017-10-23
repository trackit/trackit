import React, { Component } from 'react';
import { Route } from 'react-router-dom';

import Containers from './containers';

class App extends Component {
  render() {
    return (
      <div>
        <Containers.Main>
          <Route path={this.props.match.url + "/setup"} component={Containers.Setup.AWS.AccessManagement}/>
        </Containers.Main>
      </div>
    );
  }
}

export default App;
