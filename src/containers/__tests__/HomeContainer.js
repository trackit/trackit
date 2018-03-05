import React from 'react';
import HomeContainer from '../HomeContainer';
import Components from '../../components';
import { shallow } from "enzyme";

const Dashboard = Components.Dashboard.Dashboard;

describe('<HomeContainer />', () => {

  it('renders a <HomeContainer /> component', () => {
    const wrapper = shallow(<HomeContainer/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Dashboard/> component', () => {
    const wrapper = shallow(<HomeContainer/>);
    const costBreakdown = wrapper.find(Dashboard);
    expect(costBreakdown.length).toBe(1);
  });

});
