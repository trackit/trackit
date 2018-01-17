import moment from "moment";

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
    costs: {
      values: {},
      dates: {},
      interval: {},
      filter: {}
    }
  },
  gcp: {},
  auth: {
    token: null,
    mail: null
  },
};
