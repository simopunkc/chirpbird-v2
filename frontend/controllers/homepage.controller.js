const path = require("path");

const getHomepage = async (_, res) => {
  try {
    res.sendFile(path.join(__dirname, "../views/static/messenger.html"));
  } catch (e) {
    res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getStyleCss = async (_, res) => {
  try {
    res.sendFile(path.join(__dirname, "../views/static/style.css"));
  } catch (e) {
    res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getScriptJs = async (_, res) => {
  try {
    res.sendFile(path.join(__dirname, "../views/static/script.js"));
  } catch (e) {
    res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

const getFavicon = async (_, res) => {
  try {
    res.sendFile(path.join(__dirname, "../views/static/favicon.ico"));
  } catch (e) {
    res.status(400).json({
      status: 400,
      message: e.message,
    });
  }
}

module.exports = {
  getHomepage,
  getStyleCss,
  getScriptJs,
  getFavicon,
}