import { put, call } from 'redux-saga/effects';
import moment from 'moment';
import { getCostsSaga } from '../costsSaga';
import { getToken } from '../../misc';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";
const begin = moment().startOf('month');
const end = moment();
const filters = ["product", "day"];
const accounts = [];

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
        .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, undefined));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_COSTS_SUCCESS, costs }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getCostsSaga({begin, end, filters, accounts});

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
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
        .toEqual(call(API.AWS.Costs.getCosts, token, begin, end, filters, accounts));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_COSTS_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

});