import moment from 'moment';

export default {
  aws: {
    pricing: [],
    accounts: {
      all: [],
      external: null
    },
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
