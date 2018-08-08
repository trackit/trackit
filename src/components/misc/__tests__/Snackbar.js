import React from 'react';
import SnackBar from '../Snackbar';
import { createShallow } from '@material-ui/core/test-utils';
import { mount } from 'enzyme';

describe('<Snackbar />', () => {

    const propsNoAction = {
        variant: 'success',
        message: 'this is a test',
    };

    const propsErrorNoAction = {
        variant: 'error',
        message: 'this is a test',
    };

    const propsAction = {
        variant: 'success',
        message: 'this is a test',
        action: <button className="test-action" key="test">test</button>
    };

    let shallow;

    beforeEach(() => {
      shallow = createShallow();
    });

    it('renders a <SnackBar /> component', () => {
        const wrapper = shallow(<SnackBar {...propsNoAction}/>);
        expect(wrapper.length).toBe(1);
    });

    it('displays the message properly', () => {
        const wrapper = mount(<SnackBar {...propsNoAction}/>);
        expect(wrapper.find('span#client-snackbar').text()).toBe(propsNoAction.message);
    });

    it('displays the message type properly', () => {
        const wrapperSuccess = mount(<SnackBar {...propsNoAction}/>);
        expect(wrapperSuccess.find('.SimpleSnackbar-success-1').length).toBeGreaterThan(0);
        expect(wrapperSuccess.find('.SimpleSnackbar-error-2').length).toBe(0);
        const wrapperError = mount(<SnackBar {...propsErrorNoAction}/>);
        expect(wrapperError.find('.SimpleSnackbar-error-2').length).toBeGreaterThan(0);
        expect(wrapperError.find('.SimpleSnackbar-success-1').length).toBe(0);
    });

    it('displays an action if one is specified', () => {
        const wrapper = mount(<SnackBar {...propsAction}/>);
        expect(wrapper.find('.test-action').length).toEqual(1);
    });
    
});