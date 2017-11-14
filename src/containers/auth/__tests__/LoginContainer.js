import React from 'react';
import { LoginContainer } from '../LoginContainer';
import { shallow } from "enzyme";
import {Redirect} from "react-router-dom";
import Form from 'react-validation/build/form';

const props = {
  login: jest.fn()
};

const propsWithToken = {
  login: jest.fn(),
  token: "token"
};

describe('<LoginContainer />', () => {

  it('renders a <LoginContainer /> component', () => {
    const wrapper = shallow(<LoginContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Form/> component if token is not available', () => {
    const wrapper = shallow(<LoginContainer {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders <Redirect/> component if token is available', () => {
    const wrapper = shallow(<LoginContainer {...propsWithToken}/>);
    const redirect = wrapper.find(Redirect);
    expect(redirect.length).toBe(1);
  });

});
