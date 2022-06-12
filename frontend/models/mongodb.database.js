require('dotenv').config();
const mongoose = require('mongoose');
const { env } = process;

let _db = null
const connectDb = async () => {
  _db = await mongoose.connect(env.MONGODB_HOST, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  });
  return _db;
};

const getDb = () => {
  return _db;
}

module.exports = {
  connectDb,
  getDb,
};