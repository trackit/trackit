import React from 'react';
import Panel from '../Panel';
import { shallow } from "enzyme";

const child = <div id="child"/>;

const props = {
  title: "Title",
  children: child
};

const propsCollapsible = {
  ...props,
  collapsible: true
};

const propsCollapsed = {
  ...propsCollapsible,
  defaultCollapse: true
};

describe('<Panel />', () => {

  it('renders a <Panel /> component', () => {
    const wrapper = shallow(<Panel {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a child component', () => {
    const wrapper = shallow(<Panel {...props}/>);
    const children = wrapper.find('div#child');
    expect(children.length).toBe(1);
  });

  it('renders a <Panel /> component with default props', () => {
    const wrapper = shallow(<Panel {...props}/>);
    expect(wrapper.length).toBe(1);
    expect(wrapper.instance().props.collapsible).toBe(false);
    expect(wrapper.instance().props.defaultCollapse).toBe(false);
  });

  it('renders a <Panel /> component with custom props', () => {
    const wrapperCollapsible = shallow(<Panel {...propsCollapsible}/>);
    expect(wrapperCollapsible.length).toBe(1);
    expect(wrapperCollapsible.instance().props.collapsible).toBe(true);
    expect(wrapperCollapsible.instance().props.defaultCollapse).toBe(false);

    const wrapperCollapsed = shallow(<Panel {...propsCollapsed}/>);
    expect(wrapperCollapsed.length).toBe(1);
    expect(wrapperCollapsed.instance().props.collapsible).toBe(true);
    expect(wrapperCollapsed.instance().props.defaultCollapse).toBe(true);
  });

  it('renders no toggle when not collapsible', () => {
    const wrapper = shallow(<Panel {...props}/>);
    const toggle = wrapper.find('span.glyphicon');
    expect(toggle.length).toBe(0);
  });

  it('renders a toggle when when collapsible', () => {
    const wrapper = shallow(<Panel {...propsCollapsible}/>);
    const toggle = wrapper.find('span.glyphicon');
    expect(toggle.length).toBe(1);
  });

  it('can expand body', () => {
    const wrapper = shallow(<Panel {...propsCollapsible}/>);
    expect(wrapper.state('collapsed')).toBe(false);
    wrapper.find('div.panel-heading').prop('onClick')({ preventDefault() {} });
    expect(wrapper.state('collapsed')).toBe(true);
  });

  it('can collapse body', () => {
    const wrapperCollapsible = shallow(<Panel {...propsCollapsed}/>);
    expect(wrapperCollapsible.state('collapsed')).toBe(true);
    wrapperCollapsible.find('div.panel-heading').prop('onClick')({ preventDefault() {} });
    expect(wrapperCollapsible.state('collapsed')).toBe(false);
    const wrapperNonCollapsible = shallow(<Panel {...props}/>);
    expect(wrapperNonCollapsible.state('collapsed')).toBe(false);
    wrapperNonCollapsible.find('div.panel-heading').prop('onClick')({ preventDefault() {} });
    expect(wrapperNonCollapsible.state('collapsed')).toBe(false);
  });

});
