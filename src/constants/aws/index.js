import Pricing from './pricingTypes';
import Accounts from './accountsTypes';
import S3 from './s3Types';

export default {
	...Pricing,
	...Accounts,
	...S3,
};
