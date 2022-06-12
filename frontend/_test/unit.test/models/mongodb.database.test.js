const sinon = require("sinon");
const dbConnection = require('../../../models/mongodb.database');
const database = require('../../../models/mongodb.connection');

describe("Database connection failed", () => {
  afterEach(() => {
    sinon.restore();
  });

  it("Should catch error", async () => {
    let mockDB1 = sinon.mock(dbConnection);
    mockDB1.expects("getDb").once().returns(new Error("type"));
    const model = await database.getModel();
    const blank = {
      memberModel: {},
      roomModel: {},
      roomActivityModel: {},
    }
    expect(model).toEqual(blank);
    mockDB1.verify();
    mockDB1.restore();
  });
});