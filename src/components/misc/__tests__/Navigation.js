import React from 'react';
import { Navigation } from '../Navigation';
import { shallow } from "enzyme";

const props = {
  eventsDates: {startDate: "2018-01-01", endDate: "2018-01-30"},
  events: {status: false, values: {}},
  getData: jest.fn(),
};

describe('<Navigation />', () => {

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <Navigation /> component', () => {
    const wrapper = shallow(<Navigation {...props}/>);
    expect(wrapper.length).toEqual(1);
  });

});
