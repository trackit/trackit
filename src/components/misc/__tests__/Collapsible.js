import React from 'react';
import Collapsible from '../Collapsible';
import Collapse from "@material-ui/core/Collapse/Collapse";
import { shallow } from "enzyme";

const header = <div id="header"/>;
const child = <div id="child"/>;

const defaultProps = {
  header,
  children: child
};

describe('<Collapsible />', () => {

  it('renders a <Collapsible /> component', () => {
    const wrapper = shallow(<Collapsible {...defaultProps}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Collapse /> component', () => {
    const wrapper = shallow(<Collapsible {...defaultProps}/>);
    const collapse = wrapper.find(Collapse);
    expect(collapse.length).toBe(1);
  });

  it('can show and hide tooltip', () => {
    const wrapper = shallow(<Collapsible {...defaultProps}/>);
    expect(wrapper.state('expanded')).toBe(false);
    wrapper.instance().toggleCollapse({ preventDefault(){} });
    expect(wrapper.state('expanded')).toBe(true);
    wrapper.instance().toggleCollapse({ preventDefault(){} });
    expect(wrapper.state('expanded')).toBe(false);
  });

});
