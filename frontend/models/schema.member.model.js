module.exports = (mongo) => {
  let model
  try {
    model = mongo.model('members')
  } catch {
    const Schema = mongo.Schema;
    const modelMember = new Schema({
      email: {
        type: String,
        required: true,
        unique: true,
        lowercase: true,
        index: true,
      },
      name: {
        type: String,
        required: true,
      },
      picture: {
        type: String,
        required: true,
      },
      verified_email: {
        type: Boolean,
        required: true,
        default: true,
      },
    });
    model = mongo.model('members', modelMember);
  }
  return model;
};