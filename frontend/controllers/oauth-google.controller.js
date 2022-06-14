require('dotenv').config();
const { env } = process;
const network = require('../modules/axios.module');

const middlewareCheckGoogleProfile = async (req, res, next) => {
  try {
    if (req.cookies[env.COOKIE_ACCESS_TOKEN]) {
      const urlAPI = env.BACKEND_URI + "/oauth/google/profile";
      const config = {
        withCredentials: true,
        headers: {
          acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
        },
      }
      const response = await network.getApi(urlAPI, config);
      res.locals.profile = response.data.Message;
      next();
    } else {
      return res.redirect("/oauth/login");
    }
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const middlewareCheckEmail = async (_, res, next) => {
  if (res.locals.profile.verified_email) {
    next();
  } else {
    return res.status(400).json({
      status: 400,
      message: "user belum diverifikasi",
    })
  }
}

const getLoginPage = async (req, res) => {
  try {
    if (req.cookies.ref_token) {
      const urlAPI = env.BACKEND_URI + "/token/csrf";
      const response = await network.getApi(urlAPI, {});
      response.data.Message.map((data) => {
        res.cookie(data.cookie, data.value, data.config);
      });
      return res.redirect("/oauth/google/refresh");
    } else {
      const urlAPI = env.BACKEND_URI + "/oauth/google/url";
      const response = await network.getApi(urlAPI, {});
      let fixUrl = ""
      response.data.Message.map((data) => {
        res.cookie(data.state.cookie, data.state.value, data.state.config);
        fixUrl = decodeURIComponent(Buffer.from(data.url, 'base64').toString('ascii'));
      });
      res.render('../views/login.ejs', { printUrl: fixUrl });
    }
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getVerifyLogin = async (req, res) => {
  try {
    if (!req.query.code) {
      return res.status(400).json({
        status: 400,
        message: "query parameter code tidak ditemukan",
      });
    } else if (!req.query.state) {
      return res.status(400).json({
        status: 400,
        message: "query parameter state tidak ditemukan",
      });
    } else if (!req.cookies.xsrf_token) {
      return res.status(400).json({
        status: 400,
        message: "cookie xsrf expired",
      });
    } else {
      const urlAPI = env.BACKEND_URI + "/oauth/google/verify";
      const body = {
        state: req.query.state,
        code: req.query.code,
      }
      const config = {
        withCredentials: true,
        headers: {
          xsrf_token: req.cookies.xsrf_token,
        },
      }
      const response = await network.postApi(urlAPI, body, config);
      response.data.Message.map((data) => {
        res.cookie(data.cookie, data.value, data.config);
      });
      return res.redirect("/");
    }
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getRefreshLogin = async (req, res) => {
  try {
    if (!req.cookies.xsrf_token) {
      return res.status(400).json({
        status: 400,
        message: "cookie xsrf tidak ditemukan",
      });
    } else if (!req.cookies.ref_token) {
      return res.status(400).json({
        status: 400,
        message: "cookie refresh token tidak ditemukan",
      });
    } else {
      const urlAPI = env.BACKEND_URI + "/oauth/google/refresh";
      const body = {
        ref_token: req.cookies.ref_token,
      }
      const config = {
        withCredentials: true,
        headers: {
          xsrf_token: req.cookies.xsrf_token,
        },
      }
      const response = await network.postApi(urlAPI, body, config);
      response.data.Message.map((data) => {
        res.cookie(data.cookie, data.value, data.config);
      });
      return res.redirect("/");
    }
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getLogOut = async (req, res) => {
  try {
    if (req.cookies.ref_token) {
      const urlAPI = env.BACKEND_URI + "/oauth/logout";
      const body = {
        ref_token: req.cookies.ref_token,
      }
      const response = await network.postApi(urlAPI, body, {});
      response.data.Message.map((data) => {
        res.cookie(data.cookie, data.value, data.config);
      });
      return res.redirect("/oauth/login");
    } else {
      return res.status(401).json({
        status: 401,
        message: "cookie refresh token tidak ditemukan",
      });
    }
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getProfileLogin = async (req, res) => {
  try {
    if (req.cookies[env.COOKIE_ACCESS_TOKEN]) {
      const urlAPI = env.BACKEND_URI + "/oauth/google/profile";
      const config = {
        withCredentials: true,
        headers: {
          acc_token: req.cookies[env.COOKIE_ACCESS_TOKEN],
        },
      }
      const response = await network.getApi(urlAPI, config);
      return res.status(response.data.Status).json({
        status: response.data.Status,
        message: response.data.Message,
      });
    } else {
      return res.status(400).json({
        status: 400,
        message: "cookie acc_token tidak ditemukan",
      })
    }
  } catch (e) {
    return res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

module.exports = {
  middlewareCheckGoogleProfile,
  middlewareCheckEmail,
  getLoginPage,
  getVerifyLogin,
  getRefreshLogin,
  getLogOut,
  getProfileLogin,
}