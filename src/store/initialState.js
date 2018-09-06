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
        status: true,
        value: null
      },
      billsStatus: {
        status: false
      },
      billEdition: {
        status: true,
        value: null
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
    map: {
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
    tags: {
      charts: {},
      dates: {},
      interval: {},
      keys: {},
      values: {}
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
    resources: {
      account: '',
      EC2: {
        status: true,
        value: null
      },
      RDS: {
        status: true,
        value: null
      }
    },
  },
  events: {
    dates: {
      startDate: moment().subtract(30, 'days'),
      endDate : moment(),
    },
    values: {}
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
