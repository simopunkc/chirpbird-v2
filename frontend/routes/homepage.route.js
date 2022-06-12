const router = require("express").Router();

const {
  middlewareCheckGoogleProfile,
  middlewareCheckEmail,
} = require('../controllers/oauth-google.controller');

const {
  getHomepage,
  getStyleCss,
  getScriptJs,
  getFavicon,
} = require('../controllers/homepage.controller');

router.get("/", middlewareCheckGoogleProfile, middlewareCheckEmail, getHomepage);
router.get("/style.css", getStyleCss);
router.get("/script.js", getScriptJs);
router.get("/favicon.ico", getFavicon);

module.exports = router;