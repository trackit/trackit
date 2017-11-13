import React, { Component } from 'react';

import Containers from './containers';

class App extends Component {
  render() {
    return (
      <div>
        <Containers.Main>
          <div className="app-container" style={{paddingLeft: '60px'}}>
            <Route
              path={this.props.match.url + '/s3'}
              component={Containers.S3Analytics}
            />
            <Route
              path={this.props.match.url + "/setup"}
              component={Containers.Setup.Main}
            />
          </div>
        </Containers.Main>
      </div>
    );
  }
}

export default App;
