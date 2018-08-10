import React from 'react';
import Popover from '../Popover';
import { OverlayTrigger } from 'react-bootstrap';
import { shallow } from "enzyme";

const child = <div id="child"/>;
const tooltip = <div id="tooltip"/>;

const defaultProps = {
  tooltip
};

const propsWithChild = {
  ...defaultProps,
  children: child
};

const propsWithInfo = {
  ...defaultProps,
  info: true
};

describe('<Popover />', () => {

  it('renders a <Popover /> component', () => {
    const wrapper = shallow(<Popover {...propsWithChild}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <OverlayTrigger /> component', () => {
    const wrapper = shallow(<Popover {...propsWithChild}/>);
    const children = wrapper.find(OverlayTrigger);
    expect(children.length).toBe(1);
  });

  it('renders child component', () => {
    const wrapper = shallow(<Popover {...propsWithChild}/>);
    const children = wrapper.find('div#child');
    expect(children.length).toBe(1);
  });

  it('renders info icon component', () => {
    const wrapper = shallow(<Popover {...propsWithInfo}/>);
    const children = wrapper.find('i.fa-info-circle');
    expect(children.length).toBe(1);
  });

  it('can show and hide tooltip', () => {
    const wrapper = shallow(<Popover {...propsWithChild}/>);
    expect(wrapper.state('showPopOver')).toBe(false);
    wrapper.instance().handlePopoverOpen({ preventDefault(){} });
    expect(wrapper.state('showPopOver')).toBe(true);
    wrapper.instance().handlePopoverClose({ preventDefault(){} });
    expect(wrapper.state('showPopOver')).toBe(false);
  });

});
