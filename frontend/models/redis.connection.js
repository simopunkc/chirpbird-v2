const dbConnection = require('./redis.database');

const getModel = async () => {
  try {
    let redis = await dbConnection.getRedis();
    if (redis == null) {
      redis = await dbConnection.connectRedis();
    }
    return {
      redis,
    }
  } catch (err) {
    return {
      redis: {},
    }
  }
};

module.exports = {
  getModel
};