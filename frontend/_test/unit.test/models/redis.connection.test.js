const sinon = require("sinon");
const dbConnection = require('../../../models/redis.database');
const database = require('../../../models/redis.connection');

describe("Get database connection", () => {
  afterEach(() => {
    sinon.restore();
  });

  it("Should not reuse connection", async () => {
    let mockDB1 = sinon.mock(dbConnection);
    mockDB1.expects("getRedis").once().returns(null);
    mockDB1.expects("connectRedis").once().returns(1);
    const model = await database.getModel();
    expect(model.redis).toEqual(1);
    mockDB1.verify();
    mockDB1.restore();
  });

  it("Should reuse connection", async () => {
    await dbConnection.connectRedis();
    const model = await database.getModel();
    expect(model.redis).not.toEqual(null);
  });
});