import React from 'react';
import DialogComponent  from '../Dialog';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import { shallow } from "enzyme";

const child = <div id="child"/>;

const props = {
  buttonType: "info",
  buttonName: "Button",
  title: "Title",
  actionName: "Action",
  actionFunction: jest.fn(),
  secondActionName: "Action BIS",
  children: child
};

const propsWithoutTitle = {
  ...props,
  title: undefined
};

const propsWithoutChild = {
  ...props,
  children: undefined
};

const propsWithoutAction = {
  ...props,
  actionName: undefined,
  actionFunction: undefined
};

const propsWithCallback = {
  ...props,
  onOpen: jest.fn(),
  onClose: jest.fn(),
};

describe('<DialogComponent />', () => {

  it('renders a <DialogComponent /> component', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Dialog /> component', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    const children = wrapper.find(Dialog);
    expect(children.length).toBe(1);
  });

  it('renders title in a <DialogTitle /> component', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    const children = wrapper.find(DialogTitle);
    expect(children.length).toBe(1);
  });

  it('renders no <DialogTitle /> component if there is no title', () => {
    const wrapper = shallow(<DialogComponent {...propsWithoutTitle}/>);
    const children = wrapper.find(DialogTitle);
    expect(children.length).toBe(0);
  });

  it('renders child in a <DialogContent /> component', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    const children = wrapper.find(DialogContent);
    expect(children.length).toBe(1);
  });

  it('renders no <DialogContent /> component if there is no child', () => {
    const wrapper = shallow(<DialogComponent {...propsWithoutChild}/>);
    const children = wrapper.find(DialogContent);
    expect(children.length).toBe(0);
  });

  it('renders a <DialogActions /> component', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    const children = wrapper.find(DialogActions);
    expect(children.length).toBe(1);
  });

  it('renders 3 <button /> components', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    const children = wrapper.find('button');
    expect(children.length).toBe(3);
  });

  it('renders 2 <button /> component if there is no action', () => {
    const wrapper = shallow(<DialogComponent {...propsWithoutAction}/>);
    const children = wrapper.find('button');
    expect(children.length).toBe(2);
  });

  it('can open and close dialog', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    expect(wrapper.state('open')).toBe(false);
    wrapper.instance().openDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(true);
    wrapper.instance().closeDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(false);
  });

  it('can open and close dialog with callback', () => {
    const wrapper = shallow(<DialogComponent {...propsWithCallback}/>);
    expect(propsWithCallback.onOpen).not.toHaveBeenCalled();
    expect(propsWithCallback.onClose).not.toHaveBeenCalled();
    expect(wrapper.state('open')).toBe(false);
    wrapper.instance().openDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(true);
    expect(propsWithCallback.onOpen).toHaveBeenCalled();
    wrapper.instance().closeDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(false);
    expect(propsWithCallback.onClose).toHaveBeenCalled();
  });

  it('can execute action', () => {
    const wrapper = shallow(<DialogComponent {...props}/>);
    expect(props.actionFunction).not.toHaveBeenCalled();
    wrapper.instance().executeAction({ preventDefault(){} });
    expect(props.actionFunction).toHaveBeenCalled();
  });

});
