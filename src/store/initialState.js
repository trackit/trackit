import moment from "moment";

export default {
  aws: {
    accounts: {
      all: {
        status: false
      },
      creation: {
        status: true,
        value: null
      },
      billCreation: {
        status: true
      },
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
      charts: {},
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
