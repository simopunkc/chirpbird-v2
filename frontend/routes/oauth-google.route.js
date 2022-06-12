const router = require("express").Router();

const {
  getLoginPage,
  getVerifyLogin,
  getRefreshLogin,
  getLogOut,
  getProfileLogin,
} = require('../controllers/oauth-google.controller');

router.get("/login", getLoginPage);
router.get("/google/verify", getVerifyLogin);
router.get("/google/refresh", getRefreshLogin);
router.get("/logout", getLogOut);
router.get("/google/profile", getProfileLogin);

module.exports = router;