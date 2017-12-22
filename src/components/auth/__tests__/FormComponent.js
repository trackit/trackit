import React from 'react';
import { FormComponent } from '../FormComponent';
import { shallow } from "enzyme";
import Form from 'react-validation/build/form';

const props = {
  submit: jest.fn()
};

const form = {
  getValues: () => ({
    email: "email",
    password: "password"
  })
};

const propsForRegistration = {
  ...props,
  registration: true
};

const propsWithError = {
  ...props,
  loginStatus: {
    status: false,
    error: "error"
  }
};

describe('<FormComponent />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <FormComponent /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Form/> component if token is not available', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders <Form/> component if in registration mode', () => {
    const wrapper = shallow(<FormComponent {...propsForRegistration}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('can submit credentials to log user in', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const instance = wrapper.instance();
    instance.form = form;
    expect(props.submit).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault(){} });
    expect(props.submit).toHaveBeenCalled();
  });

  it('renders <div/> component if ther is a login error', () => {
    const wrapper = shallow(<FormComponent {...propsWithError}/>);
    const form = wrapper.find("div.alert");
    expect(form.length).toBe(1);
  });

});