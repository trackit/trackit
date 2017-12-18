import React from 'react';
import Panel, { PanelItem } from '../Panel';
import { shallow } from "enzyme";

const child = <div id="child"/>;

const props = {
  title: "Title",
  children: child
};

const propsWithChilds = {
  ...props,
  children: [child, child]
};

describe('<Panel />', () => {

  it('renders a <Panel /> component', () => {
    const wrapper = shallow(<Panel {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <PanelItem /> component', () => {
    const wrapper = shallow(<Panel {...props}/>);
    const children = wrapper.find(PanelItem);
    expect(children.length).toBe(1);
  });

  it('renders multiple <PanelItem /> components', () => {
    const wrapper = shallow(<Panel {...propsWithChilds}/>);
    const children = wrapper.find(PanelItem);
    expect(children.length).toBe(propsWithChilds.children.length);
  });

});

describe('<PanelItem />', () => {

  it('renders a <PanelItem /> component', () => {
    const wrapper = shallow(<PanelItem {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a child component', () => {
    const wrapper = shallow(<PanelItem {...props}/>);
    const children = wrapper.find('div#child');
    expect(children.length).toBe(1);
  });

});
