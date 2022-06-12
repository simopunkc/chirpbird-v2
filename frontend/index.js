require('dotenv').config();
const { env } = process;
const app = require("./server");
const hostname = '0.0.0.0';
const port = env.FRONTEND_PORT;
module.exports = app.listen(port, hostname);