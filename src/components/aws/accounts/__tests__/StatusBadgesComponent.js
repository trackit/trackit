import React from 'react';
import { StatusBadgeComponent } from '../StatusBadgesComponent';
import {
    shallow,
} from 'enzyme';

const account1 = {
    id: 42,
    roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE",
    pretty: "pretty",
    payer: true,
    billRepositories: [{
        awsAccountId: 1,
        bucket: "trackit-billing-report",
        error: "",
        id: 4,
        lastImportedManifest: "2018-08-01T13:12:25Z",
        nextPending: false,
        nextUpdate: "2018-08-02T16:12:26Z",
        prefix: "usagecost/AllHourlyToS3/"
    }]
};

const account2 = {
    id: 84,
    roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE_BIS",
    pretty: "pretty_bis",
    payer: true,
    billRepositories: [{
        awsAccountId: 1,
        bucket: "trackit-billing-report",
        error: "this is an error",
        id: 4,
        lastImportedManifest: "2018-08-01T13:12:25Z",
        nextPending: false,
        nextUpdate: "2018-08-02T16:12:26Z",
        prefix: "usagecost/AllHourlyToS3/"
    }]
};

describe('<StatusBadgeComponent />', () => {

    const propsNoSelectionWithValues = {
        accounts: {
            status: true,
            values: [account1, account2],
        },
        selected: [],
        values: { value: true }
    };
    const propsNoSelectionWithoutValues = {
        accounts: {
            status: true,
            values: [account1, account2],
        },
        selected: [],
        values: {}
    };
    const propsSelectionWithValues = {
        accounts: {
            status: true,
            values: [account1, account2],
        },
        selected: [account1],
        values: { value: true }
    };

    it('renders a <StatusBadges /> component', () => {
        const wrapper = shallow(<StatusBadgeComponent {...propsNoSelectionWithValues}/>);
        expect(wrapper.length).toBe(1);
    });


    it('renders a badge for each account', () => {
        const wrapper = shallow(<StatusBadgeComponent {...propsNoSelectionWithValues}/>);
        const badges = wrapper.find('span.account-status-badge');
        expect(badges.length).toBe(2);
    });

    it('renders account badges colors properly', () => {
        const wrapperValues = shallow(<StatusBadgeComponent {...propsNoSelectionWithValues}/>);
        const wrapperNoValues = shallow(<StatusBadgeComponent {...propsNoSelectionWithoutValues}/>);

        expect(wrapperValues.find('.green-color').length).toBe(1);
        expect(wrapperValues.find('.red-color').length).toBe(1);
        expect(wrapperNoValues.find('.orange-color').length).toBe(1);
        expect(wrapperNoValues.find('.red-color').length).toBe(1);
    });

    it('renders selection when a selection is set', () => {
        const wrapper = shallow(<StatusBadgeComponent {...propsSelectionWithValues}/>);
        const badges = wrapper.find('span.account-status-badge');
        expect(badges.length).toBe(1);
    }); 
});