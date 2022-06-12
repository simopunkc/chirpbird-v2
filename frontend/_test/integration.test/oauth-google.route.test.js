const app = require('../../server');
const request = require("supertest");
const agent = request.agent(app);
const sinon = require("sinon");
const setup = require('../rollback');
const network = require('../../modules/axios.module');

const validAccToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZXdvZ0lDSnBaQ0k2SUNJeE1UWTNNRE0zTlRneE1UWTFNVFl4TVRRMk5qTWlMQW9nSUNKbGJXRnBiQ0k2SUNKd1pXMWlaV3hoTG1Gc2JHRm9RR2R0WVdsc0xtTnZiU0lzQ2lBZ0luWmxjbWxtYVdWa1gyVnRZV2xzSWpvZ2RISjFaU3dLSUNBaWJtRnRaU0k2SUNKaWRXUnBlV0Z1ZEc4Z2MybHRieUlzQ2lBZ0ltZHBkbVZ1WDI1aGJXVWlPaUFpWW5Wa2FYbGhiblJ2SWl3S0lDQWlabUZ0YVd4NVgyNWhiV1VpT2lBaWMybHRieUlzQ2lBZ0luQnBZM1IxY21VaU9pQWlhSFIwY0hNNkx5OXNhRE11WjI5dloyeGxkWE5sY21OdmJuUmxiblF1WTI5dEwyRXRMMEZQYURFMFIycHBaaTFzVkhGSlVuZHdiMjAzTkd4ck1uVnhWM1F0YjB0cGFISmZWM0JEVTNOS1RreERUa0U5Y3prMkxXTWlMQW9nSUNKc2IyTmhiR1VpT2lBaWFXUWlDbjBLIn0.AreQinIo31287fAciC_jGsiWoNEJp8-5-H1D7MnBn1I"
const validRefToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZXlKV1lXeDFaU0k2SWpFdkx6Qm5XVjlqVVMxUE5WaElZWFJEWjFsSlFWSkJRVWRDUVZOT2QwWXRURGxKY2pObVVVbDZRekJNWTJaeWRXbERjM1ZUTlUxVWJFODRWUzFEUm5CVlFVYzJNR3hTWlY4emFFOTBTbEZ0TnpOeVJFTlVMVTV1VURsMGNsaGZVMEpOTlRKWVVtOGlMQ0pGZUhCcGNtVmtJam96TXpFNE1ERTBNekExT1gwPSJ9.q5INOFq7XotHnG_lBtOu4qBOtOW8rYi_Au0yUsYe5UQ"
const ValidXsrfToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiZXlKU1lXNWtiMjBpT2lKVU4wbzBaV2xwUzA5dmVTSXNJa1Y0Y0dseVpXUWlPakUyTkRVNU9EUTVORE45In0.5gaXBdaDGAqCqo_fUVL2nHL3R20Hut1h_FDyvbAn4sE"

describe("Integration Test /admin", () => {
  beforeAll(async () => {
    await setup();
  });

  afterEach(() => {
    sinon.restore();
  });

  describe("GET /oauth/login", () => {
    it("should show link login google", async () => {
      const response = await agent.get("/oauth/login");
      expect(response.status).toEqual(200);
    });

    it("should redirect to refresh token login google", async () => {
      let cookie = "ref_token=" + validRefToken + ";"
      const response = await agent.get("/oauth/login").set('Cookie', cookie);
      expect(response.headers["location"]).toEqual("/oauth/google/refresh");
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/oauth/login").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("GET /oauth/google/verify", () => {
    it("should xsrf_token not found", async () => {
      await agent.get("/oauth/google/verify").query({
        code: "test",
        state: "test",
      }).expect(400);
    });

    it("should query code not found", async () => {
      let cookie = "xsrf_token=" + ValidXsrfToken + ";"
      await agent.get("/oauth/google/verify").set('Cookie', cookie).query({
        state: "test",
      }).expect(400);
    });

    it("should query state not found", async () => {
      let cookie = "xsrf_token=" + ValidXsrfToken + ";"
      await agent.get("/oauth/google/verify").set('Cookie', cookie).query({
        code: "test",
      }).expect(400);
    });

    it("should get 400", async () => {
      let cookie = "xsrf_token=" + ValidXsrfToken + ";"
      await agent.get("/oauth/google/verify").set('Cookie', cookie).query({
        code: "test",
        state: "test",
      }).expect(400);
    });

    it("should valid verify login", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("postApi").once().resolves({
        data: {
          Message: [{
            cookie: "test",
            value: "test",
            config: {}
          }]
        }
      });
      let cookie = "xsrf_token=" + ValidXsrfToken + ";"
      await agent.get("/oauth/google/verify").set('Cookie', cookie).query({
        code: "test",
        state: "test",
      }).expect(302);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("GET /oauth/google/refresh", () => {
    it("should xsrf_token not found", async () => {
      await agent.get("/oauth/google/refresh").expect(400);
    });

    it("should ref_token not found", async () => {
      let cookie = "xsrf_token=" + ValidXsrfToken + ";"
      await agent.get("/oauth/google/refresh").set('Cookie', cookie).expect(400);
    });

    it("should valid refresh acc_token google", async () => {
      let cookie = "xsrf_token=" + ValidXsrfToken + ";ref_token=" + validRefToken + ";"
      const response = await agent.get("/oauth/google/refresh").set('Cookie', cookie)
      expect(response.headers["location"]).toEqual("/");
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("postApi").once().rejects(new Error("network"));
      let cookie = "xsrf_token=" + ValidXsrfToken + ";ref_token=" + validRefToken + ";"
      const response = await agent.get("/oauth/google/refresh").set('Cookie', cookie);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("GET /oauth/logout", () => {
    it("should not authorized", async () => {
      await agent.get("/oauth/logout").expect(401);
    });

    it("should valid logout", async () => {
      let cookie = "ref_token=" + validRefToken + ";"
      const response = await agent.get("/oauth/logout").set('Cookie', cookie)
      expect(response.headers["location"]).toEqual("/oauth/login");
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("postApi").once().rejects(new Error("network"));
      let cookie = "ref_token=" + validRefToken + ";"
      const response = await agent.get("/oauth/logout").set('Cookie', cookie);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("GET /oauth/google/profile", () => {
    it("should acc_token not found", async () => {
      await agent.get("/oauth/google/profile").expect(400);
    });

    it("should valid get profile google", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/oauth/google/profile").set('Cookie', cookie);
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/oauth/google/profile").set('Cookie', cookie);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });
});