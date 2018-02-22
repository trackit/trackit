import { put, call } from 'redux-saga/effects';
import { getS3DataSaga, saveS3DatesSaga, loadS3DatesSaga } from '../s3Saga';
import {getAWSAccounts, getS3Dates, getToken} from "../../misc";
import API from '../../../api';
import Moment from 'moment';
import Constants from '../../../constants';
import {getS3Dates as getS3DatesLS} from "../../../common/localStorage";

const token = "42";
const accounts = ["account1", "account2"];

describe("S3 Saga", () => {

  describe("Get S3 Data", () => {

    const dates = {
      begin: Moment().startOf('week'),
      end: Moment().endOf('week'),
    };
    const data = {
      bucket1: "data",
      bucket2: "data"
    };
    const validResponse = { success: true, data };
    const invalidResponse = { success: true, s3Data: data };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = getS3DataSaga(dates);

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(call(API.AWS.S3.getData, token, dates.begin, dates.end, accounts));

      expect(saga.next(validResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_S3_DATA_SUCCESS, data }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getS3DataSaga(dates);

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(call(API.AWS.S3.getData, token, dates.begin, dates.end, accounts));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_S3_DATA_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = getS3DataSaga(dates);

      expect(saga.next().value)
        .toEqual(getToken());

      expect(saga.next(token).value)
        .toEqual(getAWSAccounts());

      expect(saga.next(accounts).value)
        .toEqual(call(API.AWS.S3.getData, token, dates.begin, dates.end, accounts));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_S3_DATA_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

  describe("Save S3 Dates", () => {

    it("handles saga", () => {

      let saga = saveS3DatesSaga();

      expect(saga.next().value)
        .toEqual(getS3Dates());

      expect(saga.next({}).done).toBe(true);

    });

  });

  describe("Load S3 Dates", () => {

    const data = {
      startDate: Moment(),
      endDate: Moment()
    };

    const invalidData = [];

    it("handles saga with valid data", () => {

      let saga = loadS3DatesSaga();

      expect(saga.next().value)
        .toEqual(call(getS3DatesLS));

      expect(saga.next(data).value)
        .toEqual(put({type: Constants.AWS_INSERT_S3_DATES, dates: data}));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = loadS3DatesSaga();

      expect(saga.next().value)
        .toEqual(call(getS3DatesLS));

      expect(saga.next(invalidData).value)
        .toEqual(put({ type: Constants.AWS_LOAD_S3_DATES_ERROR, error: Error("Invalid data for S3 Analytics dates") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = loadS3DatesSaga();

      expect(saga.next().value)
        .toEqual(call(getS3DatesLS));

      expect(saga.next(null).value)
        .toEqual(put({ type: Constants.AWS_LOAD_S3_DATES_ERROR, error: Error("No S3 Analytics dates available") }));

      expect(saga.next().done).toBe(true);

    });

  });

});
