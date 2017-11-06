import { put, call, all } from 'redux-saga/effects';
import { getS3DataSaga } from '../s3Saga';
import API from '../../../api';
import Constants from '../../../constants';

const token = "42";

describe("S2 Saga", () => {

  describe("Get S3 Data", () => {

    const s3Data = ["s3data1", "s3data2"];
    const validResponse = { success: true, data: s3Data };
    const invalidResponse = { success: true, s3Data };
    const noResponse = { success: false };

    it("handles saga with valid data", () => {

      let saga = getS3DataSaga();

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.S3.getS3Data));

      expect(saga.next(validResponse).value)
        .toEqual(all([
          put({ type: Constants.AWS_GET_S3_DATA_SUCCESS, s3Data })
        ]));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with invalid data", () => {

      let saga = getS3DataSaga();

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.S3.getS3Data));

      expect(saga.next(invalidResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_S3_DATA_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

    it("handles saga with no response", () => {

      let saga = getS3DataSaga();

      expect(saga.next(token).value)
        .toEqual(call(API.AWS.S3.getS3Data));

      expect(saga.next(noResponse).value)
        .toEqual(put({ type: Constants.AWS_GET_S3_DATA_ERROR, error: Error("Error with request") }));

      expect(saga.next().done).toBe(true);

    });

  });

});
