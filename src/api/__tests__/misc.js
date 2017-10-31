import fetchMock from 'fetch-mock';
import { call } from '../misc';

const token = "42";

describe("API Calls", () => {

  afterEach(() => {
    fetchMock.restore();
  });

  describe("GET method", () => {

    it("handless valid response without token", () => {

      const response = { data: "data" };
      fetchMock.get('*', response);

      call('/test', 'GET', null, null)
        .then((result) => {
          expect(result).toEqual({success: true, data: response});
        });

    });

    it("handless valid response with token", () => {

      fetchMock.get('*', (url, data) => ({ token: data.headers['Authorization'] }));

      call('/test', 'GET', null, token)
        .then((result) => {
          expect(result).toEqual({success: true, data: { token }});
        });

    });
/*
    it("handless body without token", () => {

      const response = { data: "data" };
      fetchMock.get('*', (url, data) => (data.body));

      call('/test', 'GET', response, null)
        .then((result) => {
          expect(result).toEqual({success: true, data: response});
        });

    });

    it("handless body with token", () => {

      const response = { data: "data" };
      fetchMock.get('*',
        (url, data) => ({ ...JSON.parse(data.body), token: data.headers['Authorization'] })
      );

      call('/test', 'GET', response, token)
        .then((result) => {
          expect(result).toEqual({success: true, data: { ...response, token }});
        });

    });
*/
  });

  describe("POST method", () => {

    it("handless valid response", () => {

      const response = { data: "data" };
      fetchMock.post('*', (url, data) => (data.body));

      call('/test', 'POST', response, null)
        .then((result) => {
          expect(result).toEqual({success: true, data: response});
        });

    });

    it("handless valid response with token", () => {

      const response = { data: "data" };
      fetchMock.post('*',
        (url, data) => ({ ...JSON.parse(data.body), token: data.headers['Authorization'] })
      );

      call('/test', 'POST', response, token)
        .then((result) => {
          expect(result).toEqual({success: true, data: { ...response, token }});
        });

    });

/*
    it("handless no body without token", () => {

      const response = { data: "data" };
      fetchMock.post('*', response);

      call('/test', 'POST', null, null)
        .then((result) => {
          expect(result).toEqual({success: true, data: response});
        });

    });

    it("handless no body with token", () => {

      fetchMock.post('*', (url, data) => ({ token: data.headers['Authorization'] }));

      call('/test', 'POST', null, token)
        .then((result) => {
          expect(result).toEqual({success: true, data: { token }});
        });

    });
*/
  });

});
