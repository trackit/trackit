import React from 'react';
import Picture from '../Picture';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import { shallow } from "enzyme";

const button = <div id="button"/>;

const props = {
  src: "src",
  alt: "alt",
  text: "text"
};

const propsWithPreview = {
  ...props,
  preview: true
};

const propsWithButton = {
  ...props,
  button
};

describe('<Picture />', () => {

  it('renders a <Picture /> component', () => {
    const wrapper = shallow(<Picture {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders an <img /> component', () => {
    const wrapper = shallow(<Picture {...props}/>);
    const image = wrapper.find(`img[src="${props.src}"][alt="${props.alt}"]`);
    expect(image.length).toBe(1);
  });

  it('renders an <img /> component as preview', () => {
    const wrapper = shallow(<Picture {...propsWithPreview}/>);
    const preview = wrapper.find(`img[src="${propsWithPreview.src}"][alt="${propsWithPreview.alt}"]`);
    expect(preview.length).toBe(2);
  });

  it('renders a <button /> component', () => {
    const wrapper = shallow(<Picture {...props}/>);
    const button = wrapper.find(`button`);
    expect(button.length).toBe(1);
  });

  it('renders a custom button component', () => {
    const wrapper = shallow(<Picture {...propsWithButton}/>);
    const customButton = wrapper.find(`div#button`);
    expect(customButton.length).toBe(1);
    const button = wrapper.find(`button`);
    expect(button.length).toBe(0);
  });

  it('renders a <Dialog /> component', () => {
    const wrapper = shallow(<Picture {...props}/>);
    const children = wrapper.find(Dialog);
    expect(children.length).toBe(1);
  });

  it('renders a <DialogContent /> component', () => {
    const wrapper = shallow(<Picture {...props}/>);
    const children = wrapper.find(DialogContent);
    expect(children.length).toBe(1);
  });

  it('can open and close dialog', () => {
    const wrapper = shallow(<Picture {...props}/>);
    expect(wrapper.state('open')).toBe(false);
    wrapper.instance().openDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(true);
    wrapper.instance().closeDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(false);
  });

/*
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
*/
});
