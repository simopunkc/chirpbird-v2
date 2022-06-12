const app = require('../../server');
const request = require("supertest");
const agent = request.agent(app);
const sinon = require("sinon");
const setup = require('../rollback');
const network = require('../../modules/axios.module');
const path = require("path");

const validAccToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZXdvZ0lDSnBaQ0k2SUNJeE1UWTNNRE0zTlRneE1UWTFNVFl4TVRRMk5qTWlMQW9nSUNKbGJXRnBiQ0k2SUNKd1pXMWlaV3hoTG1Gc2JHRm9RR2R0WVdsc0xtTnZiU0lzQ2lBZ0luWmxjbWxtYVdWa1gyVnRZV2xzSWpvZ2RISjFaU3dLSUNBaWJtRnRaU0k2SUNKaWRXUnBlV0Z1ZEc4Z2MybHRieUlzQ2lBZ0ltZHBkbVZ1WDI1aGJXVWlPaUFpWW5Wa2FYbGhiblJ2SWl3S0lDQWlabUZ0YVd4NVgyNWhiV1VpT2lBaWMybHRieUlzQ2lBZ0luQnBZM1IxY21VaU9pQWlhSFIwY0hNNkx5OXNhRE11WjI5dloyeGxkWE5sY21OdmJuUmxiblF1WTI5dEwyRXRMMEZQYURFMFIycHBaaTFzVkhGSlVuZHdiMjAzTkd4ck1uVnhWM1F0YjB0cGFISmZWM0JEVTNOS1RreERUa0U5Y3prMkxXTWlMQW9nSUNKc2IyTmhiR1VpT2lBaWFXUWlDbjBLIn0.AreQinIo31287fAciC_jGsiWoNEJp8-5-H1D7MnBn1I"
const validProfile = {
  data: {
    Status: 200,
    Message: {
      verified_email: true
    }
  }
}

describe("Integration Test", () => {
  beforeAll(async () => {
    await setup();
  });

  afterEach(() => {
    sinon.restore();
  });

  describe("GET /undefined", () => {
    it("should get 404", async () => {
      await agent.get("/undefined").expect(404);
    });
  });

  describe("GET static file", () => {
    it("GET /favicon.ico", async () => {
      await agent.get("/favicon.ico").expect(200);
    });

    it("GET /script.js", async () => {
      await agent.get("/script.js").expect(200);
    });

    it("GET /style.css", async () => {
      await agent.get("/style.css").expect(200);
    });
  });

  describe("GET /", () => {
    it("should redirect to login page", async () => {
      const response = await agent.get("/");
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("invalid jwt cookie acc_token", async () => {
      let cookie = "acc_token=blank;"
      await agent.get("/").set('Cookie', cookie).expect(400);
    });

    it("valid cookie jwt acc_token", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      await agent.get("/").set('Cookie', cookie).expect(200);
    });
  });

  describe("should catch error", () => {
    it("error homepage", async () => {
      let mockAxios = sinon.mock(network);
      let mockPath = sinon.mock(path);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockPath.expects("join").once().returns(new Error("file"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockPath.verify();
      mockAxios.restore();
      mockPath.restore();
    });

    it("error favicon", async () => {
      let mockPath = sinon.mock(path);
      mockPath.expects("join").once().returns(new Error("file"));
      const response = await agent.get("/favicon.ico");
      expect(response.status).toEqual(400);
      mockPath.verify();
      mockPath.restore();
    });

    it("error script", async () => {
      let mockPath = sinon.mock(path);
      mockPath.expects("join").once().returns(new Error("file"));
      const response = await agent.get("/script.js");
      expect(response.status).toEqual(400);
      mockPath.verify();
      mockPath.restore();
    });

    it("error css", async () => {
      let mockPath = sinon.mock(path);
      mockPath.expects("join").once().returns(new Error("file"));
      const response = await agent.get("/style.css");
      expect(response.status).toEqual(400);
      mockPath.verify();
      mockPath.restore();
    });
  });
});