import moment from "moment";

export default {
  aws: {
    accounts: {
      all: [],
      external: null,
    },
    s3: {
      dates: {
        startDate: moment().startOf('month'),
        endDate: moment()
      },
      values: {}
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
