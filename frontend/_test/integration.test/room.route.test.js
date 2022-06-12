const app = require('../../server');
const request = require("supertest");
const agent = request.agent(app);
const sinon = require("sinon");
const setup = require('../rollback');
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
    await setup();
  });

  afterEach(() => {
    sinon.restore();
  });

  describe("GET /room/:id", () => {
    it("should blank acc_token", async () => {
      const response = await agent.get("/room/GRt4eASqkK8")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.get("/room/GRt4eASqkK8").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid get room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/room/GRt4eASqkK8").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("getApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.get("/room/GRt4eASqkK8").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/addMember", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/addMember")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/addMember").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should failed add member", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/addMember").set('Cookie', cookie).send(body);
      expect(response.status).toEqual(400);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/addMember").set('Cookie', cookie).send(body);
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/exit", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/exit")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/exit").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid exit room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.put("/room/GRt4eASqkK8/exit").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.put("/room/GRt4eASqkK8/exit").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/rename", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/rename")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/rename").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid rename room", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        name: "pemburu mie ayam",
      }
      const response = await agent.put("/room/GRt4eASqkK8/rename").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        name: "pemburu mie ayam",
      }
      const response = await agent.put("/room/GRt4eASqkK8/rename").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/memberToModerator", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/memberToModerator")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/memberToModerator").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should invalid id target", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/memberToModerator").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/memberToModerator").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/ModeratorToMember", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/ModeratorToMember")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/ModeratorToMember").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should invalid id target", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/ModeratorToMember").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/ModeratorToMember").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/kickMember", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/kickMember")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/kickMember").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should failed job desk", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/kickMember").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_target: "anonymous@gmail.com",
      }
      const response = await agent.put("/room/GRt4eASqkK8/kickMember").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("POST /room/:id/newChat", () => {
    it("should blank acc_token", async () => {
      const response = await agent.post("/room/GRt4eASqkK8/newChat")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.post("/room/GRt4eASqkK8/newChat").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid post chat", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_parent: "",
        message: "lorem ipsum dolor sir amet"
      }
      const response = await agent.post("/room/GRt4eASqkK8/newChat").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(201);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("postApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      let body = {
        id_parent: "",
        message: "lorem ipsum dolor sir amet"
      }
      const response = await agent.post("/room/GRt4eASqkK8/newChat").set('Cookie', cookie).send(body)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/enableNotif", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/enableNotif")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/enableNotif").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid enable notification", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.put("/room/GRt4eASqkK8/enableNotif").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.put("/room/GRt4eASqkK8/enableNotif").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("PUT /room/:id/disableNotif", () => {
    it("should blank acc_token", async () => {
      const response = await agent.put("/room/GRt4eASqkK8/disableNotif")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.put("/room/GRt4eASqkK8/disableNotif").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid disable notification", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.put("/room/GRt4eASqkK8/disableNotif").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("putApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.put("/room/GRt4eASqkK8/disableNotif").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });

  describe("DELETE /room/:id/deleteRoom", () => {
    it("should blank acc_token", async () => {
      const response = await agent.delete("/room/GRt4eASqkK8/deleteRoom")
      expect(response.headers["location"]).toEqual("/oauth/login");
      expect(response.status).toEqual(302);
    });

    it("should unauthorized user", async () => {
      let cookie = "acc_token=blank;"
      const response = await agent.delete("/room/GRt4eASqkK8/deleteRoom").set('Cookie', cookie)
      expect(response.status).toEqual(400);
    });

    it("should valid delete single chat", async () => {
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.delete("/room/GRt4eASqkK8/deleteRoom").set('Cookie', cookie)
      expect(response.status).toEqual(200);
    });

    it("should catch error", async () => {
      let mockAxios = sinon.mock(network);
      mockAxios.expects("getApi").once().resolves(validProfile);
      mockAxios.expects("deleteApi").once().rejects(new Error("network"));
      let cookie = "acc_token=" + validAccToken + ";"
      const response = await agent.delete("/room/GRt4eASqkK8/deleteRoom").set('Cookie', cookie)
      expect(response.status).toEqual(400);
      mockAxios.verify();
      mockAxios.restore();
    });
  });
});