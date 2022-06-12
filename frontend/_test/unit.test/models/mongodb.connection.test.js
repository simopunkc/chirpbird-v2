const sinon = require("sinon");
const dbConnection = require('../../../models/mongodb.database');
const database = require('../../../models/mongodb.connection');
const mongoose = require('mongoose');
const mongo = {
  Schema: {},
  model: {},
}

describe("Get database connection", () => {
  afterEach(() => {
    sinon.restore();
  });

  it("Should not reuse connection", async () => {
    let mockDB1 = sinon.mock(dbConnection);
    mockDB1.expects("getDb").once().returns(null);
    const model = await database.getModel();
    expect(model).toHaveProperty("memberModel");
    expect(model).toHaveProperty("roomModel");
    expect(model).toHaveProperty("roomActivityModel");
    mockDB1.verify();
    mockDB1.restore();
  });

  it("Should reuse connection", async () => {
    let mockDB1 = sinon.mock(mongoose);
    mockDB1.expects("connect").once().resolves(mongo);
    await dbConnection.connectDb();
    const model = await dbConnection.getDb();
    expect(model).not.toEqual(null);
    mockDB1.verify();
    mockDB1.restore();
  });
});