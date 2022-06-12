const router = require("express").Router();

const {
  middlewareCheckGoogleProfile,
  middlewareCheckEmail,
} = require('../controllers/oauth-google.controller');

const {
  getRoomActivity,
  deleteRoomActivity,
} = require('../controllers/activity.controller');

router.get(/\/([A-Za-z0-9]+)/, middlewareCheckGoogleProfile, middlewareCheckEmail, getRoomActivity);
router.delete(/\/([A-Za-z0-9]+)\/deleteChat/, middlewareCheckGoogleProfile, middlewareCheckEmail, deleteRoomActivity);

module.exports = router;