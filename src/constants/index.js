import AWS from './aws';
import GCP from './gcp';
import Auth from './auth';
import User from './user';
import Dashboard from './dashboardTypes';
import Events from './events';
import Highlevel from './highlevel';

export default {
	...AWS,
	...GCP,
	...Auth,
	...Dashboard,
	...User,
	...Events,
	...Highlevel
};
