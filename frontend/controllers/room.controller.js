require('dotenv').config();
const { env } = process;
const network = require('../modules/axios.module');

const getSingleRoom = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0];
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

const putAddMember = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/addMember";
    const body = {
      id_target: req.body.id_target,
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

const putExitRoom = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/exit";
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.putApi(urlApi, {}, config);
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

const putRenameRoom = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/rename";
    const body = {
      name: req.body.name,
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

const putMemberToModerator = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/memberToModerator";
    const body = {
      id_target: req.body.id_target,
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

const putModeratorToMember = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/ModeratorToMember";
    const body = {
      id_target: req.body.id_target,
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

const putKickMember = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/kickMember";
    const body = {
      id_target: req.body.id_target,
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

const postNewChat = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/newChat";
    const body = {
      id_parent: req.body.id_parent,
      message: req.body.message,
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

const putEnableNotif = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/enableNotif";
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.putApi(urlApi, {}, config);
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

const putDisableNotif = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/disableNotif";
    const config = {
      withCredentials: true,
      headers: {
        acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
      },
    }
    const response = await network.putApi(urlApi, {}, config);
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

const deleteRoom = async (req, res) => {
  try {
    const urlApi = env.BACKEND_URI + "/room/" + req.params[0] + "/deleteRoom";
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
  getSingleRoom,
  putAddMember,
  putExitRoom,
  putRenameRoom,
  putMemberToModerator,
  putModeratorToMember,
  putKickMember,
  postNewChat,
  putEnableNotif,
  putDisableNotif,
  deleteRoom,
}