import React from 'react';
import { Route } from 'react-router-dom';
import { shallow } from 'enzyme';
import App from '../App';
import Containers from '../containers';

const props = {
  match: {
    url: null
  }
};

describe('<App />', () => {

  it('renders a <App /> component', () => {
    const wrapper = shallow(<App {...props} />);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Containers.Main /> component', () => {
    const wrapper = shallow(<App {...props} />);
    const main = wrapper.find(Containers.Main);
    expect(main.length).toBe(1);
  });

  it('renders a <Route /> component', () => {
    const wrapper = shallow(<App {...props} />);
    const route = wrapper.find(Route);
    expect(route.length).toBe(2);
  });

});
