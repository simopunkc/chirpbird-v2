const setup = async () => {
  try {
    const cache = require('../models/redis.connection');
    const redisConnection = await cache.getModel();
    const mongodb = require('../models/mongodb.connection');
    const mongodbModel = await mongodb.getModel();
    const options1 = { upsert: true, new: true, setDefaultsOnInsert: true };
    const queryMember1 = { email: "pembela.allah@gmail.com" };
    const updateMember1 = {
      email: 'pembela.allah@gmail.com',
      name: 'budiyanto simo',
      picture: 'https://lh3.googleusercontent.com/a-/AOh14Gjif-lTqIRwpom74lk2uqWt-oKihr_WpCSsJNLCNA=s96-c',
      verified_email: true
    };
    const queryMember2 = { email: "budiyanto.simo@gmail.com" };
    const updateMember2 = {
      email: 'budiyanto.simo@gmail.com',
      name: 'Budiyanto No Last Name',
      picture: 'https://lh3.googleusercontent.com/a-/AOh14GgJIRFxftSlzPiNmHa2v9b959IWf5MBOSVYCiSY=s96-c',
      verified_email: true
    };
    const queryRoom1 = { id_primary: "GRt4eASqkK8" };
    const updateRoom1 = {
      id_primary: 'GRt4eASqkK8',
      id_member_creator: 'pembela.allah@gmail.com',
      name: 'pemburu jamur barat',
      list_id_member: ['pembela.allah@gmail.com'],
      list_id_member_moderator: ['pembela.allah@gmail.com'],
      list_id_member_banned: [],
      list_id_member_enable_notification: [],
      date_created: '2022-01-30T14:40:23.133120073+07:00',
      date_last_activity: '2022-01-30T14:40:23.133120073+07:00',
      link_join: 'LJMbKMSqH6cw'
    };
    const queryRoom2 = { id_primary: "GxtteWSpOMn" };
    const updateRoom2 = {
      id_primary: 'GxtteWSpOMn',
      id_member_creator: 'budiyanto.simo@gmail.com',
      name: 'pemburu jamur barat',
      list_id_member: ['budiyanto.simo@gmail.com'],
      list_id_member_moderator: ['budiyanto.simo@gmail.com'],
      list_id_member_banned: [],
      list_id_member_enable_notification: [],
      date_created: '2022-01-30T14:40:23.133120073+07:00',
      date_last_activity: '2022-01-30T14:40:23.133120073+07:00',
      link_join: 'LJbkKYrZnp9q'
    };
    const queryRoomActivity1 = { id_primary: "RAvKCc61kmk2" };
    const updateRoomActivity1 = {
      id_primary: 'RAvKCc61kmk2',
      id_parent: '',
      id_room: 'GRt4eASqkK8',
      id_member_actor: 'pembela.allah@gmail.com',
      id_member_target: '',
      type_activity: 'group_created',
      message: 'budiyanto simo membuat group baru',
      date_created: '2022-01-30T14:40:23.133120073+07:00',
      list_id_member_unread: []
    };
    const queryRoomActivity2 = { id_primary: "RAbtWxa4klk4" };
    const updateRoomActivity2 = {
      id_primary: 'RAbtWxa4klk4',
      id_parent: '',
      id_room: 'GxtteWSpOMn',
      id_member_actor: 'budiyanto.simo@gmail.com',
      id_member_target: '',
      type_activity: 'group_created',
      message: 'Budiyanto No Last Name membuat group baru',
      date_created: '2022-01-30T14:40:23.133120073+07:00',
      list_id_member_unread: []
    };
    await Promise.all([
      redisConnection.redis.flushall(),
      mongodbModel.memberModel.findOneAndUpdate(queryMember1, updateMember1, options1).exec(),
      mongodbModel.memberModel.findOneAndUpdate(queryMember2, updateMember2, options1).exec(),
      mongodbModel.roomModel.findOneAndUpdate(queryRoom1, updateRoom1, options1).exec(),
      mongodbModel.roomModel.findOneAndUpdate(queryRoom2, updateRoom2, options1).exec(),
      mongodbModel.roomActivityModel.findOneAndUpdate(queryRoomActivity1, updateRoomActivity1, options1).exec(),
      mongodbModel.roomActivityModel.findOneAndUpdate(queryRoomActivity2, updateRoomActivity2, options1).exec()
    ]);
  } catch (e) {
    console.log(e.message)
  }
}

module.exports = async () => {
  await setup();
};