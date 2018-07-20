import Login from './loginActions';
import Token from './tokenActions';
import Logout from './logoutActions';
import Registration from './registrationActions';
import ForgotPassword from './forgotPasswordActions';

export default {
	...Login,
	...Token,
	...Logout,
	...Registration,
	...ForgotPassword
};
