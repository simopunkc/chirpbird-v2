const dbConnection = require('./mongodb.database');

const getModel = async () => {
  try {
    let dbMongo = await dbConnection.getDb();
    if (dbMongo == null) {
      dbMongo = await dbConnection.connectDb();
    }
    const memberModel = require('./schema.member.model')(dbMongo);
    const roomModel = require('./schema.room.model')(dbMongo);
    const roomActivityModel = require('./schema.room-activity.model')(dbMongo);
    return {
      memberModel,
      roomModel,
      roomActivityModel,
    }
  } catch (err) {
    return {
      memberModel: {},
      roomModel: {},
      roomActivityModel: {},
    }
  }
};

module.exports = {
  getModel
};