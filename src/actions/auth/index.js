import Login from './loginActions';
import Token from './tokenActions';
import Logout from './logoutActions';

export default {
	...Login,
	...Token,
	...Logout
};
