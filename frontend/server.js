const express = require('express');
const cors = require('cors');
const cookieParser = require("cookie-parser");
const homepageRoute = require("./routes/homepage.route");
const oauthGoogleRoute = require("./routes/oauth-google.route");
const roomRoute = require("./routes/room.route");
const activityRoute = require("./routes/activity.route");
const messengerRoute = require("./routes/messenger.route");
const memberRoute = require("./routes/member.route");
const app = express();
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(cors());
app.use(cookieParser());
app.use("", homepageRoute);
app.use("/oauth", oauthGoogleRoute);
app.use("/member", memberRoute);
app.use("/room", roomRoute);
app.use("/activity", activityRoute);
app.use("/messenger", messengerRoute);
app.use((_, res) => {
  res.sendStatus(404);
});
module.exports = app