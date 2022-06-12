require('dotenv').config();
var Redis = require("ioredis");
const { env } = process;

let _redis = null
const connectRedis = async () => {
  _redis = new Redis({
    sentinels: [
      {
        port: env.REDIS_PORT,
        host: env.REDIS_SENTINEL1,
      },
      {
        port: env.REDIS_PORT,
        host: env.REDIS_SENTINEL2,
      },
      {
        port: env.REDIS_PORT,
        host: env.REDIS_SENTINEL3,
      },
    ],
    name: "masterredis1",
  });
  return _redis;
};

const getRedis = () => {
  return _redis;
}

module.exports = {
  connectRedis,
  getRedis,
};