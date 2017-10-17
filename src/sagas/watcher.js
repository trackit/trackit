import * as AWS from './aws';
import * as GCP from './gcp';
import * as Auth from './auth';


export default {
  ...AWS,
  ...GCP,
  ...Auth,
};
