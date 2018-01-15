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

const form = {
  getValues: () => ({
    roleArn: "roleArn",
    pretty: "pretty"
  })
};

describe('<FormComponent />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <FormComponent /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders 2 <Input /> components', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(2);
  });

  it('renders 1 <Button /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const button = wrapper.find(Button);
    expect(button.length).toBe(1);
  });

  it('can submit info to create account', () => {
    const wrapper = shallow(<FormComponent {...propsWithExternal}/>);
    const instance = wrapper.instance();
    instance.form = form;
    expect(props.submit).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault(){} });
    expect(props.submit).toHaveBeenCalled();
  });

  it('renders 2 <Input /> components inside with accounts data', () => {
    const wrapper = shallow(<FormComponent {...propsWithAccount}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(2);
  });

  it('can submit info to update account', () => {
    const wrapper = shallow(<FormComponent {...propsWithAccount}/>);
    const instance = wrapper.instance();
    instance.form = form;
    expect(props.submit).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault(){} });
    expect(props.submit).toHaveBeenCalled();
  });

  it('can open and close dialog', () => {
    const wrapper = shallow(<FormComponent {...propsWithAccount}/>);
    expect(wrapper.state('open')).toBe(false);
    wrapper.instance().openDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(true);
    wrapper.instance().closeDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(false);
  });

});
