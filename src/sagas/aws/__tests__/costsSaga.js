import { all, put, call } from 'redux-saga/effects';
import moment from 'moment';
import { getCostsSaga, saveChartsSaga, loadChartsSaga } from '../costsSaga';
import { getCostBreakdownCharts as getCostBreakdownChartsLS } from '../../../common/localStorage';
import {getToken, getAWSAccounts, getCostBreakdownCharts} from '../../misc';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";
const begin = moment().startOf('month');
const end = moment();
const filters = ["product", "day"];
const accounts = ["account1", "account2"];

describe("Costs Saga", () => {

  describe("Get Costs", () => {

    const costs = ["cost1", "cost2"];
    const validResponse = { success: true, data: costs };
    const invalidResponse = { success: true, costs };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = getCostsSaga({begin, end, filters, accounts});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_COSTS_SUCCESS, costs }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with valid data and without accounts", () => {

      let saga = getCostsSaga({begin, end, filters});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next([]).value)
        .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, []));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_COSTS_SUCCESS, costs }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getCostsSaga({begin, end, filters, accounts});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_COSTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = getCostsSaga({begin, end, filters, accounts});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_COSTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Save Charts", () => {

    it("handles saga", () => {

      let saga = saveChartsSaga();

      expect(saga.next().value)
        .toEqual(getCostBreakdownCharts());

      expect(saga.next({}).done).toBe(true);

    });

  });

  describe("Load Charts", () => {

    const data = {
      charts: [],
      dates: {},
      interval: {},
      filter: {}
    };

    const invalidData = {
      invalid: {}
    };

    it("handles saga with valid data", () => {

      let saga = loadChartsSaga();

      expect(saga.next().value)
        .toEqual(call(getCostBreakdownChartsLS));

      expect(saga.next(data).value)
        .toEqual(all([
          put({type: Constants.AWS_INSERT_CHARTS, charts: data.charts}),
          put({type: Constants.AWS_INSERT_COSTS_DATES, dates: data.dates}),
          put({type: Constants.AWS_INSERT_COSTS_INTERVAL, interval: data.interval}),
          put({type: Constants.AWS_INSERT_COSTS_FILTER, filter: data.filter})
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = loadChartsSaga();

      expect(saga.next().value)
        .toEqual(call(getCostBreakdownChartsLS));

      expect(saga.next(null).value)
        .toEqual(put({ type: Constants.AWS_LOAD_CHARTS_ERROR, error: Error("No cost breakdown chart available") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = loadChartsSaga();

      expect(saga.next().value)
        .toEqual(call(getCostBreakdownChartsLS));

      expect(saga.next(invalidData).value)
        .toEqual(put({ type: Constants.AWS_LOAD_CHARTS_ERROR, error: Error("Invalid data for cost breakdown charts") }));

      expect(saga.next().done).toBe(true);

    });

  });

});