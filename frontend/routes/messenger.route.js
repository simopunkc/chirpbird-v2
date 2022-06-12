const router = require("express").Router();

const {
  middlewareCheckGoogleProfile,
  middlewareCheckEmail,
} = require('../controllers/oauth-google.controller');

const {
  getMemberChat,
} = require('../controllers/messenger.controller');

router.get(/\/([A-Za-z0-9]+)\/page(\d+)/, middlewareCheckGoogleProfile, middlewareCheckEmail, getMemberChat);

module.exports = router;