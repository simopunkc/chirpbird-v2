const router = require("express").Router();

const {
  middlewareCheckGoogleProfile,
  middlewareCheckEmail,
} = require('../controllers/oauth-google.controller');

const {
  getMemberRoom,
  postCreateRoom,
  putJoinRoom,
} = require('../controllers/member.controller');

router.get(/\/room\/page(\d+)/, middlewareCheckGoogleProfile, middlewareCheckEmail, getMemberRoom);
router.post("/room/create", middlewareCheckGoogleProfile, middlewareCheckEmail, postCreateRoom);
router.put("/room/join", middlewareCheckGoogleProfile, middlewareCheckEmail, putJoinRoom);

module.exports = router;