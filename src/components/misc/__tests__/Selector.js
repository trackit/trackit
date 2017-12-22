import React from 'react';
import Selector  from '../Selector';
import { shallow } from "enzyme";

const props = {
  values: {
    value: "Value",
    otherValue: "Other value"
  },
  selected: "value",
  selectValue: jest.fn()
};

describe('<Selector />', () => {

  it('renders a <Selector /> component', () => {
    const wrapper = shallow(<Selector {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <select /> component', () => {
    const wrapper = shallow(<Selector {...props}/>);
    const children = wrapper.find("select");
    expect(children.length).toBe(1);
  });

  it('renders <option /> components', () => {
    const wrapper = shallow(<Selector {...props}/>);
    const children = wrapper.find("option");
    expect(children.length).toBe(Object.keys(props.values).length);
  });

  it('can select value', () => {
    const wrapper = shallow(<Selector {...props}/>);
    expect(props.selectValue).not.toHaveBeenCalled();
    wrapper.instance().handleValueSelection({ target: { value: "value" } });
    expect(props.selectValue).toHaveBeenCalled();
  });

});
