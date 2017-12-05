import React from 'react';
import { NavbarHeader } from '../NavbarHeader';
import { shallow } from "enzyme";

const props = {
  signOut: jest.fn()
};

describe('<NavbarHeader />', () => {

  beforeEach(() => {
    props.signOut.mockReset();
  });

  it('renders a <NavbarHeader /> component', () => {
    const wrapper = shallow(<NavbarHeader {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('dispatches a logout action', () => {
    const wrapper = shallow(<NavbarHeader {...props}/>);
    const logout = wrapper.find('a');
    expect(props.signOut).not.toHaveBeenCalled();
    logout.prop('onClick')();
    expect(props.signOut).toHaveBeenCalled();
  });

  it('renders without user menu', () => {
    const wrapper = shallow(<NavbarHeader {...props}/>);
    expect(wrapper.state('userMenuExpanded')).toBe(false);
  });

  it('can expand user menu', () => {
    const wrapper = shallow(<NavbarHeader {...props}/>);
    expect(wrapper.state('userMenuExpanded')).toBe(false);
    wrapper.find('button.navbar-user-dropdown-toggle').prop('onClick')();
    expect(wrapper.state('userMenuExpanded')).toBe(true);
  });

});
