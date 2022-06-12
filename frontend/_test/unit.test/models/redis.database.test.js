const sinon = require("sinon");
const dbConnection = require('../../../models/redis.database');
const database = require('../../../models/redis.connection');

describe("Database connection failed", () => {
  afterEach(() => {
    sinon.restore();
  });

  it("Should catch error", async () => {
    let mockDB1 = sinon.mock(dbConnection);
    mockDB1.expects("getRedis").once().rejects(new Error("type"));
    const model = await database.getModel();
    expect(model.redis).toEqual({});
    mockDB1.verify();
    mockDB1.restore();
  });
});