import * as AWS from './aws';
import * as GCP from './gcp';

export default {
  ...AWS,
  ...GCP
};
