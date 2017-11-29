import React from 'react';
import DeleteConfirmation  from '../DeleteConfirmation';
import Dialog from '../Dialog';
import { shallow } from "enzyme";

const props = {
  entity: "entity",
  confirm: jest.fn()
};

const propsWithoutEntity = {
  ...props,
  entity: undefined
};

describe('<DeleteConfirmation />', () => {

  it('renders a <DeleteConfirmation /> component', () => {
    const wrapper = shallow(<DeleteConfirmation {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Dialog /> component', () => {
    const wrapper = shallow(<DeleteConfirmation {...props}/>);
    const children = wrapper.find(Dialog);
    expect(children.length).toBe(1);
  });

  it('renders a <DeleteConfirmation /> component with default entity', () => {
    const wrapper = shallow(<DeleteConfirmation {...propsWithoutEntity}/>);
    expect(wrapper.length).toBe(1);
  });

});
