import React from 'react';
import { Navigation } from '../Navigation';
import { shallow } from "enzyme";

const props = {
  signOut: jest.fn()
};

describe('<Navigation />', () => {

  beforeEach(() => {
    props.signOut.mockReset();
  });

  it('renders a <Navigation /> component', () => {
    const wrapper = shallow(<Navigation {...props}/>);
    expect(wrapper.length).toEqual(1);
  });

  it('dispatches a logout action', () => {
    const wrapper = shallow(<Navigation {...props}/>);
    wrapper.setState({userMenuExpanded: true});
    const logout = wrapper.find('a');
    expect(props.signOut).not.toHaveBeenCalled();
    logout.prop('onClick')();
    expect(props.signOut).toHaveBeenCalled();
  });

  it('renders without user menu', () => {
    const wrapper = shallow(<Navigation {...props}/>);
    expect(wrapper.state('userMenuExpanded')).toBe(false);
  });

  it('can expand user menu', () => {
    const wrapper = shallow(<Navigation {...props}/>);
    expect(wrapper.state('userMenuExpanded')).toBe(false);
    wrapper.find('button').prop('onClick')();
    expect(wrapper.state('userMenuExpanded')).toBe(true);
  });

});
