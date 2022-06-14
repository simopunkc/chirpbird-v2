const setup = async () => {
  const dbMongo = require('../models/mongodb.database');
  const mongodbConnection = await dbMongo.connectDb();
  await mongodbConnection.connection.db.dropDatabase();
  const rollback = require('./rollback');
  await rollback.setup();
}

module.exports = async () => {
  await setup();
};