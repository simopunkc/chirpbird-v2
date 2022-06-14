const app = require('../../server');
const request = require("supertest");
const agent = request.agent(app);
const sinon = require("sinon");
const rollback = require('../rollback');
const network = require('../../modules/axios.module');

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
    await rollback.setup();
  });

  afterEach(() => {
    sinon.restore();
  });

  describe("GET /member/room/page:pid", () => {
    it("should blank acc_token", async () => {
      const response = await agent.get("/member/room/page1")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.get("/member/room/page1").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid get list room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/member/room/page1").set('Cookie', cookie);
      expect(response.status).toEqual(200);
    });

    it("should not get unavailable page", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/member/room/page2").set('Cookie', cookie);
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("getApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/member/room/page1").set('Cookie', cookie);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("POST /member/room/create", () => {
    it("should blank acc_token", async () => {
      const response = await agent.post("/member/room/create")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.post("/member/room/create").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid create room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        name: "pemburu jamur barat",
      }
      const response = await agent.post("/member/room/create").set('Cookie', cookie).send(body);
      expect(response.status).toEqual(201);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("postApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        name: "pemburu jamur barat",
      }
      const response = await agent.post("/member/room/create").set('Cookie', cookie).send(body);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /member/room/join", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/member/room/join")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/member/room/join").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should unknown room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        token: "AABBCCDDEE",
      }
      const response = await agent.put("/member/room/join").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
    });

    it("should valid join room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        token: "LJbkKYrZnp9q",
      }
      const response = await agent.put("/member/room/join").set('Cookie', cookie).send(body);
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        token: "LJbkKYrZnp9q",
      }
      const response = await agent.put("/member/room/join").set('Cookie', cookie).send(body);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });
});