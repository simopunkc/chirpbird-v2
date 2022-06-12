const router = require("express").Router();

const {
  middlewareCheckGoogleProfile,
  middlewareCheckEmail,
} = require('../controllers/oauth-google.controller');

const {
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
} = require('../controllers/room.controller');

router.get(/\/([A-Za-z0-9]+)/, middlewareCheckGoogleProfile, middlewareCheckEmail, getSingleRoom);
router.put(/\/([A-Za-z0-9]+)\/addMember/, middlewareCheckGoogleProfile, middlewareCheckEmail, putAddMember);
router.put(/\/([A-Za-z0-9]+)\/exit/, middlewareCheckGoogleProfile, middlewareCheckEmail, putExitRoom);
router.put(/\/([A-Za-z0-9]+)\/rename/, middlewareCheckGoogleProfile, middlewareCheckEmail, putRenameRoom);
router.put(/\/([A-Za-z0-9]+)\/memberToModerator/, middlewareCheckGoogleProfile, middlewareCheckEmail, putMemberToModerator);
router.put(/\/([A-Za-z0-9]+)\/ModeratorToMember/, middlewareCheckGoogleProfile, middlewareCheckEmail, putModeratorToMember);
router.put(/\/([A-Za-z0-9]+)\/kickMember/, middlewareCheckGoogleProfile, middlewareCheckEmail, putKickMember);
router.post(/\/([A-Za-z0-9]+)\/newChat/, middlewareCheckGoogleProfile, middlewareCheckEmail, postNewChat);
router.put(/\/([A-Za-z0-9]+)\/enableNotif/, middlewareCheckGoogleProfile, middlewareCheckEmail, putEnableNotif);
router.put(/\/([A-Za-z0-9]+)\/disableNotif/, middlewareCheckGoogleProfile, middlewareCheckEmail, putDisableNotif);
router.delete(/\/([A-Za-z0-9]+)\/deleteRoom/, middlewareCheckGoogleProfile, middlewareCheckEmail, deleteRoom);

module.exports = router;