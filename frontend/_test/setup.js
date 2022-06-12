const setup = async () => {
  try {
    const dbMongo = require('../models/mongodb.database');
    const mongodbConnection = await dbMongo.connectDb();
    await mongodbConnection.connection.db.dropDatabase();
  } catch (e) {
    console.log(e.message);
  }
}

module.exports = async () => {
  await setup();
};