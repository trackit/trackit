import Constants from '../../constants';

export default {
	getFilters: () => ({type: Constants.EVENTS_GET_FILTERS}),
	clearGetFilters: () => ({type: Constants.EVENTS_GET_FILTERS_CLEAR}),
	setFilters: (filters) => ({type: Constants.EVENTS_SET_FILTERS, filters}),
	clearSetFilters: () => ({type: Constants.EVENTS_SET_FILTERS_CLEAR}),
};
