import React from 'react';
import { RegisterContainer } from '../RegisterContainer';
import { shallow } from "enzyme";
import Components from '../../../components';

const Form = Components.Auth.Form;

const props = {
  register: jest.fn(),
  clear: jest.fn()
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

});
