import React from 'react';
import { shallow } from 'enzyme';
import { Route, Redirect } from 'react-router-dom';
import PrivateRoute from '../PrivateRoute';

const emptyDiv = () => (
  <div/>
);

const storeWithoutAuth = {
  getState: () => {
    return {
      auth: {
      }
    }
  }
};

const storeWithAuth = {
  getState: () => {
    return {
      auth: {
        token: "42"
      }
    }
  }
};

describe('<PrivateRoute/>', () => {

  it('renders a <PrivateRoute /> component', () => {
    const wrapper = shallow(<PrivateRoute component={emptyDiv} store={storeWithoutAuth}/>);
    expect(wrapper.length).toEqual(1);
  });

  it('renders a <Route /> component inside', () => {
    const wrapper = shallow(<PrivateRoute component={emptyDiv} store={storeWithoutAuth}/>);
    expect(wrapper.find(Route)).toHaveLength(1);
  });

  it('renders a <Route /> with <Redirect /> as component if not auth', () => {
    const wrapper = shallow(<PrivateRoute component={emptyDiv} store={storeWithoutAuth}/>);
    const route = wrapper.find(Route);
    const component = route.props().render({location: "/app"});
    expect(component.type).toEqual(Redirect);
  });

  it('renders a <Route /> with <div /> as component if auth', () => {
    const wrapper = shallow(<PrivateRoute component={emptyDiv} store={storeWithAuth}/>);
    const route = wrapper.find(Route);
    const component = route.props().render();
    expect(component.type).toBe(emptyDiv);
  });

});
