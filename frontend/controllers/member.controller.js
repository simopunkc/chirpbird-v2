require('dotenv').config();
const { env } = process;
const network = require('../modules/axios.module');

const getMemberRoom = async (req, res) => {
  try {
    const page = parseInt(req.params[0]) > 1 ? req.params[0] : "1";
    const urlApi = env.BACKEND_URI + "/member/room/page" + page;
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

const postCreateRoom = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/member/room/create";
    const body = {
      name: req.body.name,
    }
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.postApi(urlApi, body, config);
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

const putJoinRoom = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/member/room/join";
    const body = {
      token: req.body.token,
    }
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.putApi(urlApi, body, config);
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
  getMemberRoom,
  postCreateRoom,
  putJoinRoom,
}