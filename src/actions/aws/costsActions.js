import Constants from '../../constants';

export default {
	getCosts: (begin, end, filters, accounts=undefined) => ({
		type: Constants.AWS_GET_COSTS,
    begin,
    end,
    filters,
    accounts
	})
};
