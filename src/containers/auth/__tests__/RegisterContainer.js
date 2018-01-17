import React from 'react';
import { RegisterContainer } from '../RegisterContainer';
import { shallow } from "enzyme";
import {Redirect} from "react-router-dom";
import Components from '../../../components';

const Form = Components.Auth.Form;

const props = {
  register: jest.fn(),
  clear: jest.fn()
};

const propsWithRegistration = {
  ...props,
  registrationStatus: { status: true }
};

describe('<RegisterContainer />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <RegisterContainer /> component', () => {
    const wrapper = shallow(<RegisterContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <Form/> component if token is not available', () => {
    const wrapper = shallow(<RegisterContainer {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
    expect(form.props().registration).toBe(true);
  });

  it('renders <Redirect/> component if registration is done', () => {
    const wrapper = shallow(<RegisterContainer {...propsWithRegistration}/>);
    const redirect = wrapper.find(Redirect);
    expect(redirect.length).toBe(1);
  });

  it('clear registration when component unmount', () => {
    const wrapper = shallow(<RegisterContainer {...propsWithRegistration}/>);
    expect(propsWithRegistration.clear).not.toHaveBeenCalled();
    wrapper.instance().componentWillUnmount();
    expect(propsWithRegistration.clear).toHaveBeenCalled();
  });

});
