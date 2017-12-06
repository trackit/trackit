import Login from './loginActions';
import Token from './tokenActions';
import Logout from './logoutActions';
import Registration from './registrationActions';

export default {
	...Login,
	...Token,
	...Logout,
	...Registration
};
