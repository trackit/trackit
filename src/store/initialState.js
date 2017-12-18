import moment from 'moment';

export default {
  aws: {
    accounts: {
      all: [],
      external: null,
    },
    s3: {
      view: {
        startDate: moment().startOf('month'),
        endDate: moment()
      }
    },
    costs: null
  },
  gcp: {},
  auth: {
    token: null,
  },
};
