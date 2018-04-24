import React from 'react';
import { ReportsContainer } from '../ReportsContainer';
import { shallow } from 'enzyme';

describe('<ReportsContainer />', () => {
  it('renders a <ReportsContainer /> component', () => {
    const wrapper = shallow(<ReportsContainer/>);
    expect(wrapper.length).toBe(1);
  });
});
