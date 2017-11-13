import React from 'react';
import FormComponent from '../FormComponent';
import { shallow } from "enzyme/build/index";
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';
import Validations from '../../../../common/forms';

const Validation = Validations.AWSAccount;

const props = {
  submit: jest.fn()
};

const propsWithExternal = {
  submit: jest.fn(),
  external: "external_test"
};

describe('<FormComponent />', () => {

  it('renders a <FormComponent /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders 3 <Input /> components', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(3);
  });

  it('renders 1 <Button /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const button = wrapper.find(Button);
    expect(button.length).toBe(1);
  });

  it('renders external value in a disabled dedicated <Input /> component', () => {
    const wrapper = shallow(<FormComponent {...propsWithExternal}/>);
    const input = wrapper.find(Input).first();
    expect(input.prop("disabled")).toBe(true);
    expect(input.prop("value")).toBe(propsWithExternal.external);
  });

  it('renders 3 <Input /> components inside', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(3);
  });
/*
  it('dispatches a submit action', () => {
    const wrapper = shallow(<FormComponent {...propsWithExternal}/>);
    const inputs = wrapper.find(Input);
    console.log(inputs);
    console.log(inputs[1]);
    console.log(inputs[2]);
    const button = wrapper.find(Button);
    expect(props.submit.mock.calls.length).toBe(0);
    button.simulate('click');
    expect(props.submit.mock.calls.length).toBe(1);
  });
*/
/*
  it('renders without user menu', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    expect(wrapper.state('userMenuExpanded')).toBe(false);
  });
*/
/*
  it('can expand user menu', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    expect(wrapper.state('userMenuExpanded')).toBe(false);
    wrapper.find('button.navbar-user-dropdown-toggle').prop('onClick')();
    expect(wrapper.state('userMenuExpanded')).toBe(true);
  });
*/
});
