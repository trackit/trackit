import { all, put, call } from 'redux-saga/effects';
import moment from 'moment';
import { getDataSaga, saveDashboardSaga, loadDashboardSaga, initDashboardSaga } from '../dashboardSaga';
import { getDashboard as getDashboardLS } from '../../common/localStorage';
import {getToken, getAWSAccounts, getDashboard, initialDashboard} from '../misc';
import API from '../../api';
import Constants from '../../constants';

const token = "42";
const begin = moment().startOf('month');
const end = moment();
const filters = ["product", "day"];
const accounts = ["account1", "account2"];

describe("Dashboard Saga", () => {

  describe("Get Data", () => {

    const id = "id";
    const data = "data";
    const validResponse = { success: true, data };
    const errorResponse = { success: true, data: { error: "Error" } };
    const invalidResponse = { success: true, values: data };
    const noResponse = { success: false };

    describe("Cost Breakdown", () => {

      const itemType = "costbreakdown";

      it("handles saga with valid data", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

        expect(saga.next(validResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_SUCCESS, id, data }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with valid data and without accounts", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next([]).value)
          .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, []));

        expect(saga.next(validResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_SUCCESS, id, data }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with invalid data", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

        expect(saga.next(invalidResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error with request") }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with error in data", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

        expect(saga.next(errorResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error") }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with no response", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

        expect(saga.next(noResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error with request") }));

        expect(saga.next().done).toBe(true);

      });

    });

    describe("S3 Analytics", () => {

      const itemType = "s3";

      it("handles saga with valid data", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.S3.getData, token, begin, end, accounts));

        expect(saga.next(validResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_SUCCESS, id, data }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with valid data and without accounts", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next([]).value)
          .toEqual(call(API.AWS.S3.getData, token, begin, end, []));

        expect(saga.next(validResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_SUCCESS, id, data }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with invalid data", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.S3.getData, token, begin, end, accounts));

        expect(saga.next(invalidResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error with request") }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with error in data", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.S3.getData, token, begin, end, accounts));

        expect(saga.next(errorResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error") }));

        expect(saga.next().done).toBe(true);

      });

      it("handles saga with no response", () => {

        let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

        expect(saga.next().value)
          .toEqual(getToken());

        expect(saga.next(token).value)
          .toEqual(getAWSAccounts());

        expect(saga.next(accounts).value)
          .toEqual(call(API.AWS.S3.getData, token, begin, end, accounts));

        expect(saga.next(noResponse).value)
          .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error with request") }));

        expect(saga.next().done).toBe(true);

      });

    });

    it("handles saga with invalid itemType", () => {

      const itemType = null;

      let saga = getDataSaga({id, itemType, begin, end, filters, accounts});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(put({ type: Constants.DASHBOARD_GET_VALUES_ERROR, id, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Save Dashboard", () => {

    it("handles saga", () => {

      let saga = saveDashboardSaga();

      expect(saga.next().value)
        .toEqual(getDashboard());

      expect(saga.next({}).done).toBe(true);

    });

  });

  describe("Load Dashboard", () => {

    const data = {
      items: {},
      dates: {},
      intervals: {},
      filters: {}
    };

    const invalidData = {
      invalid: {}
    };

    it("handles saga with valid data", () => {

      let saga = loadDashboardSaga();

      expect(saga.next().value)
        .toEqual(call(getDashboardLS));

      expect(saga.next(data).value)
        .toEqual(all([
          put({type: Constants.DASHBOARD_INSERT_ITEMS, items: data.items}),
          put({type: Constants.DASHBOARD_INSERT_ITEMS_DATES, dates: data.dates}),
          put({type: Constants.DASHBOARD_INSERT_ITEMS_INTERVAL, intervals: data.intervals}),
          put({type: Constants.DASHBOARD_INSERT_ITEMS_FILTER, filters: data.filters})
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = loadDashboardSaga();

      expect(saga.next().value)
        .toEqual(call(getDashboardLS));

      expect(saga.next(invalidData).value)
        .toEqual(put({ type: Constants.DASHBOARD_INIT_ITEMS_ERROR, error: Error("Invalid data for dashboard") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = loadDashboardSaga();

      expect(saga.next().value)
        .toEqual(call(getDashboardLS));

      expect(saga.next(null).value)
        .toEqual(put({ type: Constants.DASHBOARD_INIT_ITEMS_ERROR, error: Error("No dashboard available") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Init Dashboard", () => {

    const validResponse = initialDashboard();
    const invalidResponse = {};

    it("handles saga with valid data", () => {

      let saga = initDashboardSaga();

      expect(saga.next().value)
        .toEqual(call(initialDashboard));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({type: Constants.DASHBOARD_INSERT_ITEMS, items: validResponse.items}),
          put({type: Constants.DASHBOARD_INSERT_ITEMS_DATES, dates: validResponse.dates}),
          put({type: Constants.DASHBOARD_INSERT_ITEMS_INTERVAL, intervals: validResponse.intervals}),
          put({type: Constants.DASHBOARD_INSERT_ITEMS_FILTER, filters: validResponse.filters})
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = initDashboardSaga();

      expect(saga.next().value)
        .toEqual(call(initialDashboard));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.DASHBOARD_INIT_ITEMS_ERROR, error: Error("Invalid data for dashboard") }));

      expect(saga.next().done).toBe(true);

    });

  });

});