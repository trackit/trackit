import React from 'react';
import SetupContainer from '../SetupContainer';
import Panels from '..';
import { shallow } from 'enzyme';

describe('<SetupContainer />', () => {

  it('renders a <SetupContainer /> component', () => {
    const wrapper = shallow(<SetupContainer/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Panels.AWS.Accounts /> component', () => {
    const wrapper = shallow(<SetupContainer/>);
    const panel = wrapper.find(Panels.AWS.Accounts);
    expect(panel.length).toBe(1);
  });

});
