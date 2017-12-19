import React, { Component } from 'react';
import { Route } from 'react-router-dom';
import Containers from './containers';

class App extends Component {
  render() {
    return (
      <div>
        <Containers.Main>
          <div className="app-container">
            <Route
              path={this.props.match.url} exact
              component={Containers.AWS.CostBreakdown}
            />
            <Route
              path={this.props.match.url + '/s3'}
              component={Containers.AWS.S3Analytics}
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
