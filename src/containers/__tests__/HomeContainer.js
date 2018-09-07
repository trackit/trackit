import React from 'react';
import HomeContainer from '../HomeContainer';
import Components from '../../components';
import { shallow } from "enzyme";

const HighLevel = Components.HighLevel.HighLevel;

describe('<HomeContainer />', () => {

  it('renders a <HomeContainer /> component', () => {
    const wrapper = shallow(<HomeContainer/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <HighLevel/> component', () => {
    const wrapper = shallow(<HomeContainer/>);
    const costBreakdown = wrapper.find(HighLevel);
    expect(costBreakdown.length).toBe(1);
  });

});
