import React from 'react';
import HomeContainer from '../HomeContainer';
import Containers from '../../containers';
import { shallow } from "enzyme";

const CostBreakdown = Containers.AWS.CostBreakdown;

describe('<HomeContainer />', () => {

  it('renders a <HomeContainer /> component', () => {
    const wrapper = shallow(<HomeContainer/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <CostBreakdown/> component', () => {
    const wrapper = shallow(<HomeContainer/>);
    const costBreakdown = wrapper.find(CostBreakdown);
    expect(costBreakdown.length).toBe(1);
  });

});
