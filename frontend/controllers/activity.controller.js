require('dotenv').config();
const { env } = process;
const network = require('../modules/axios.module');

const getRoomActivity = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/activity/" + req.params[0];
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.getApi(urlApi, config);
    return res.status(response.data.Status).json({
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

const deleteRoomActivity = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/activity/" + req.params[0] + "/deleteChat";
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.deleteApi(urlApi, config);
    return res.status(response.data.Status).json({
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
  getRoomActivity,
  deleteRoomActivity,
}