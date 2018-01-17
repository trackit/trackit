import React from 'react';
import Panel, { PanelItem } from '../Panel';
import { shallow } from "enzyme";

const child = <div id="child"/>;
const childWithClasses = <div id="child" className="child-with-class"/>;

const props = {
  children: child
};

const propsWithChilds = {
  children: [child, child]
};

const propsWithNullChild = {
  children: [child, null]
};

const propsWithClasses = {
  children: childWithClasses
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

  it('handles null child', () => {
    const wrapper = shallow(<Panel {...propsWithNullChild}/>);
    const children = wrapper.find(PanelItem);
    expect(children.length).toBe(propsWithNullChild.children.length - 1);
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

  it('renders a child component and support classes', () => {
    const wrapper = shallow(<PanelItem {...propsWithClasses}/>);
    const children = wrapper.find('div.white-box.child-with-class');
    expect(children.length).toBe(1);
  });

});
