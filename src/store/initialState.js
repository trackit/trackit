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
      values: null,
      dates: {
        startDate: moment().subtract(1, 'month').startOf('month'),
        endDate: moment().subtract(1, 'month').endOf('month')
      },
      interval: "day",
      filter: "product"
    }
  },
  gcp: {},
  auth: {
    token: null,
  },
};
