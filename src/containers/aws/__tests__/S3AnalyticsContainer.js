import React from 'react';
import { S3AnalyticsContainer } from '../S3AnalyticsContainer';
import Components from '../../../components';
import { shallow } from "enzyme";

const S3Analytics = Components.AWS.S3Analytics;

const props = {
  getS3Data: jest.fn(),
  s3Data: [{
    _id: "id",
    size: 42,
    storage_cost: 42,
    bw_cost: 42,
    total_cost: 42,
    transfer_in: 42,
    transfer_out: 42
  }]
};

describe('<S3AnalyticsContainer />', () => {

  it('renders a <S3AnalyticsContainer /> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders <S3Analytics.Infos/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.Infos);
    expect(navigation.length).toBe(1);
  });

  it('renders <S3Analytics.BarChart/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.BarChart);
    expect(navigation.length).toBe(1);
  });

  it('renders <S3Analytics.Table/> component', () => {
    const wrapper = shallow(<S3AnalyticsContainer {...props}/>);
    const navigation = wrapper.find(S3Analytics.Table);
    expect(navigation.length).toBe(1);
  });

});
