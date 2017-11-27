import moment from 'moment';

export default {
  aws: {
    s3: {
      view: {
        startDate: moment().startOf('month'),
        endDate: moment()
      }
    }
  },
  gcp: {},
  auth: {
    token: null,
  },
};
