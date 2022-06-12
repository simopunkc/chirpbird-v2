const app = require('../../server');
const request = require("supertest");
const agent = request.agent(app);
const sinon = require("sinon");
const network = require('../../modules/axios.module');
const setup = require('../rollback');

const validAccToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZXdvZ0lDSnBaQ0k2SUNJeE1UWTNNRE0zTlRneE1UWTFNVFl4TVRRMk5qTWlMQW9nSUNKbGJXRnBiQ0k2SUNKd1pXMWlaV3hoTG1Gc2JHRm9RR2R0WVdsc0xtTnZiU0lzQ2lBZ0luWmxjbWxtYVdWa1gyVnRZV2xzSWpvZ2RISjFaU3dLSUNBaWJtRnRaU0k2SUNKaWRXUnBlV0Z1ZEc4Z2MybHRieUlzQ2lBZ0ltZHBkbVZ1WDI1aGJXVWlPaUFpWW5Wa2FYbGhiblJ2SWl3S0lDQWlabUZ0YVd4NVgyNWhiV1VpT2lBaWMybHRieUlzQ2lBZ0luQnBZM1IxY21VaU9pQWlhSFIwY0hNNkx5OXNhRE11WjI5dloyeGxkWE5sY21OdmJuUmxiblF1WTI5dEwyRXRMMEZQYURFMFIycHBaaTFzVkhGSlVuZHdiMjAzTkd4ck1uVnhWM1F0YjB0cGFISmZWM0JEVTNOS1RreERUa0U5Y3prMkxXTWlMQW9nSUNKc2IyTmhiR1VpT2lBaWFXUWlDbjBLIn0.AreQinIo31287fAciC_jGsiWoNEJp8-5-H1D7MnBn1I"
const validProfile = {
  data: {
    Status: 200,
    Message: {
      verified_email: true
    }
  }
}

describe("Integration Test /admin", () => {
  beforeAll(async () => {
    await setup();
  });

  afterEach(() => {
    sinon.restore();
  });

  describe("GET /messenger/:id/page:pid", () => {
    it("should blank acc_token", async () => {
      const response = await agent.get("/messenger/GRt4eASqkK8/page1")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.get("/messenger/GRt4eASqkK8/page1").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid get list chat", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/messenger/GRt4eASqkK8/page1").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should not get list chat because unavailable page", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/messenger/GRt4eASqkK8/page2").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("getApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/messenger/GRt4eASqkK8/page1").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });
});