import React, { Component } from 'react';
import { Route } from 'react-router-dom';
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
          </div>
        </Containers.Main>
      </div>
    );
  }
}

export default App;
