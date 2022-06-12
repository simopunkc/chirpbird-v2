const sinon = require("sinon");
const axios = require("axios");
const network = require('../../../modules/axios.module');

describe("Unit Test module axios", () => {
  afterEach(() => {
    sinon.restore();
  });

  describe("GET", () => {
    it("Should valid", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("get").once().resolves({
        data: 1
      });
      await network.getApi("", {});
      mockAxios.verify();
      mockAxios.restore();
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("get").once().rejects(new Error("network"));
      await network.getApi("", {});
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("POST", () => {
    it("Should valid", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("post").once().resolves({
        data: 1
      });
      await network.postApi("", {}, {});
      mockAxios.verify();
      mockAxios.restore();
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("post").once().rejects(new Error("network"));
      await network.postApi("", {}, {});
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT", () => {
    it("Should valid", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("put").once().resolves({
        data: 1
      });
      await network.putApi("", {}, {});
      mockAxios.verify();
      mockAxios.restore();
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("put").once().rejects(new Error("network"));
      await network.putApi("", {}, {});
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("DELETE", () => {
    it("Should valid", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("delete").once().resolves({
        data: 1
      });
      await network.deleteApi("", {});
      mockAxios.verify();
      mockAxios.restore();
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(axios);
      mockAxios.expects("delete").once().rejects(new Error("network"));
      await network.deleteApi("", {});
      mockAxios.verify();
      mockAxios.restore();
    });
  });
});