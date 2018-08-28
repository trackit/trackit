import React from 'react';
import FormComponent from '../FormComponent';
import { shallow } from 'enzyme';
import Form from 'react-validation/build/form';
import Input from 'react-validation/build/input';
import Button from 'react-validation/build/button';

const props = {
  submit: jest.fn(),
  account: 42,
  clear: jest.fn()
};

const propsWithBill = {
  ...props,
  bill: {
    bucket: "s3://test.test",
    prefix: "/path/to/bill"
  }
};

const form = {
  getValues: () => ({
    bucket: "s3://test.test/path/to/bill"
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

  it('renders a tutorial component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const form = wrapper.find("div.tutorial");
    expect(form.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders 2 <Input /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(2);
  });

  it('renders 1 <Button /> component', () => {
    const wrapper = shallow(<FormComponent {...props}/>);
    const button = wrapper.find(Button);
    expect(button.length).toBe(1);
  });

  it('renders 2 <Input /> component inside with bill data', () => {
    const wrapper = shallow(<FormComponent {...propsWithBill}/>);
    const inputs = wrapper.find(Input);
    expect(inputs.length).toBe(2);
  });

  it('can submit info to update bill', () => {
    const wrapper = shallow(<FormComponent {...propsWithBill}/>);
    const instance = wrapper.instance();
    instance.form = form;
    expect(props.submit).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault(){} });
    wrapper.instance().submit({ preventDefault(){} });
    expect(props.submit).toHaveBeenCalled();
  });

  it('can open and close dialog', () => {
    const wrapper = shallow(<FormComponent {...propsWithBill}/>);
    expect(wrapper.state('open')).toBe(false);
    wrapper.instance().openDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(true);
    wrapper.instance().closeDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(false);
  });

});
