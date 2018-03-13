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
        startDate: moment().subtract(1, 'months').startOf('month'),
        endDate: moment().subtract(1, 'months').endOf('month')
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
  dashboard: {
    items: {},
    values: {},
    intervals: {},
    filters: {},
    dates: {
      startDate: moment().subtract(1, 'month').startOf('month'),
      endDate: moment().subtract(1, 'month').endOf('month')
    }
  },
  auth: {
    token: null,
    mail: null
  },
};
