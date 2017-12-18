import Accounts from './accountsTypes';
import S3 from './s3Types';
import Costs from './costsTypes';

export default {
	...Accounts,
	...S3,
	...Costs
};
