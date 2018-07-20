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
      billsStatus: {
        status: false
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
    },
    reports: {
      account: '',
      download: {
        failed: false,
      },
      reportList: {
        status: false,
        values: []
      }
    },
  },
  gcp: {},
  user: {
    viewers: {
      all: {status: true, value: null},
      creation: {status: true, value: null}
    },
  },
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
    mail: null,
    recoverStatus: {status: true, value: null},
    renewStatus: {status: true, value: null}
  },
};
