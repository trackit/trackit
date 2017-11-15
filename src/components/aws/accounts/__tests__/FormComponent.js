import React from 'react';
import FormComponent from '../FormComponent';
import { shallow } from 'enzyme';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';

const props = {
  submit: jest.fn()
};

const propsWithExternal = {
  ...props,
  external: "external_test"
};

const propsWithAccount = {
  ...props,
  account: {
    id: 42,
    roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE",
    pretty: "pretty",
    bills: []
  }
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
  it('renders 3 <Input /> components inside with accounts data', () => {
    const wrapper = shallow(<FormComponent {...propsWithAccount}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(3);
  });
*/
});
