import React from 'react';
import MainContainer from '../MainContainer';
import Components from '../../components';
import { shallow } from "enzyme";

const Navigation = Components.Misc.Navigation;

const props = {
  children: (<div id="children"/>)
};

describe('<MainContainer />', () => {

  it('renders a <MainContainer /> component', () => {
    const wrapper = shallow(<MainContainer/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Navigation/> component', () => {
    const wrapper = shallow(<MainContainer/>);
    const navigation = wrapper.find(Navigation);
    expect(navigation.length).toBe(1);
  });

  it('renders children inside component', () => {
    const wrapper = shallow(<MainContainer {...props}/>);
    const children = wrapper.find('div#children');
    expect(children.length).toBe(1);
  });

});
