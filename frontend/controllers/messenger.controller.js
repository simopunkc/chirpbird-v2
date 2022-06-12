require('dotenv').config();
const { env } = process;
const network = require('../modules/axios.module');

const getMemberChat = async (req, res) => {
  try {
    const page = parseInt(req.params[1]) > 1 ? req.params[1] : "1";
    const urlApi = env.BACKEND_URI + "/messenger/" + req.params[0] + "/page" + page;
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    };
    const response = await network.getApi(urlApi, config);
    res.status(response.data.Status).json({
      status: response.data.Status,
      message: response.data.Message,
    });
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

module.exports = {
  getMemberChat,
}