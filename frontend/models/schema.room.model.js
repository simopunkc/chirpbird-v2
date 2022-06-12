module.exports = (mongo) => {
  let model
  try {
    model = mongo.model('rooms');
  } catch {
    const Schema = mongo.Schema;
    const modelRoom = new Schema({
      id_primary: {
        type: String,
        required: true,
        unique: true,
        index: true,
      },
      id_member_creator: {
        type: String,
        required: true,
        index: true,
      },
      name: {
        type: String,
        required: true,
      },
      list_id_member: {
        type: [String],
        required: true,
      },
      list_id_member_moderator: {
        type: [String],
        required: true,
      },
      list_id_member_banned: {
        type: [String],
        required: true,
      },
      list_id_member_enable_notification: {
        type: [String],
        required: true,
      },
      date_created: { type: Date, default: Date.now },
      date_last_activity: { type: Date, default: Date.now },
      link_join: {
        type: String,
        required: true,
        unique: true,
        index: true,
      },
    });
    model = mongo.model('rooms', modelRoom);
  }
  return model;
};