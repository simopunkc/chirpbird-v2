module.exports = (mongo) => {
  let model
  try {
    model = mongo.model('room_activities');
  } catch {
    const Schema = mongo.Schema;
    const modelRoomActivity = new Schema({
      id_primary: {
        type: String,
        required: true,
        unique: true,
        index: true,
      },
      id_parent: {
        type: String,
        required: true,
      },
      id_room: {
        type: String,
        required: true,
      },
      id_member_actor: {
        type: String,
        required: true,
      },
      id_member_target: {
        type: String,
        required: true,
      },
      type_activity: {
        type: String,
        required: true,
      },
      message: {
        type: String,
        required: true,
      },
      list_id_member_unread: {
        type: [String],
        required: true,
      },
      date_created: { type: Date, default: Date.now },
    });
    model = mongo.model('room_activities', modelRoomActivity);
  }
  return model;
};