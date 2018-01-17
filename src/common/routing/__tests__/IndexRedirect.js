import React from 'react';
import { shallow } from 'enzyme';
import { Redirect } from 'react-router-dom';
import IndexRedirect from '../IndexRedirect';

describe('<IndexRedirect/>', () => {

  it('renders a <IndexRedirect /> component', () => {
    const wrapper = shallow(<IndexRedirect/>);
    expect(wrapper.length).toEqual(1);
  });

  it('renders a <Redirect /> component inside', () => {
    const wrapper = shallow(<IndexRedirect/>);
    expect(wrapper.find(Redirect)).toHaveLength(1);
  });

  it('redirects to /app', () => {
    const wrapper = shallow(<IndexRedirect/>);
    expect(wrapper.contains(<Redirect to={{pathname: '/app'}}/>)).toBe(true);
  });

});
