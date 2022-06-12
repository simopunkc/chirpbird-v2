package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CheckUserIsExist(id_user string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return hasil
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err = sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("members")
			opts := options.Count()
			count, err := coll.CountDocuments(ctx2, bson.D{
				{Key: "email", Value: id_user},
			}, opts)
			if err != nil {
				return nil, err
			}
			if count > 0 {
				hasil = true
			}
			return count, err
		},
		txnOpts)
	if err != nil {
		return false
	}
	return hasil
}

func CreateChat(id_log string, id_parent string, id_group string, user string, target string, message string, date_log string, list_unread []string, log_type string) bool {
	temp := CreateRoomActivity(id_log, id_parent, id_group, user, target, log_type, message, date_log, list_unread)
	return temp
}

func CheckUserIsJoinGroup(id_room string, id_user string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err = sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Count()
			list := []string{id_user}
			count, err := coll.CountDocuments(ctx2, bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list},
				}},
			}, opts)
			if err != nil {
				return nil, err
			}
			if count > 0 {
				hasil = true
			}
			return count, err
		},
		txnOpts)
	if err != nil {
		return false
	}
	return hasil
}

func ExitRoom(id_room string, id_member string, id_log string, date_log string, list_unread []string, log_type string, log_message string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list := []string{id_member}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: bson.D{
					{Key: "$ne", Value: id_member},
				}},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list},
				}},
			}
			update := bson.D{
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member", Value: id_member},
				}},
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_moderator", Value: id_member},
				}},
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_enable_notification", Value: id_member},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
				CreateRoomActivity(id_log, "", id_room, id_member, "", log_type, log_message, date_log, list_unread)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func RenameRoom(id_room string, id_moderator string, new_name string, id_log string, date_log string, list_unread []string, log_type string, log_message string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list1 := []string{id_moderator}
			list2 := []bson.D{
				{
					{Key: "list_id_member_moderator", Value: bson.D{
						{Key: "$in", Value: list1},
					}},
				},
				{
					{Key: "id_member_actor", Value: id_moderator},
				},
			}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "$or", Value: list2},
			}
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "name", Value: new_name},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
				CreateRoomActivity(id_log, "", id_room, id_moderator, "", log_type, log_message, date_log, list_unread)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func KickMember(id_room string, id_moderator string, id_target string, id_log string, date_log string, list_unread []string, log_type string, log_message string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list1 := []string{id_target}
			list2 := []string{id_moderator}
			list3 := []bson.D{
				{
					{Key: "list_id_member_moderator", Value: bson.D{
						{Key: "$in", Value: list2},
					}},
				},
				{
					{Key: "id_member_actor", Value: id_moderator},
				},
			}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: bson.D{
					{Key: "$ne", Value: id_target},
				}},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list1},
				}},
				{Key: "$or", Value: list3},
			}
			update := bson.D{
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member", Value: id_target},
				}},
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_moderator", Value: id_target},
				}},
				{Key: "$addToSet", Value: bson.D{
					{Key: "list_id_member_banned", Value: id_target},
				}},
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_enable_notification", Value: id_target},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
				CreateRoomActivity(id_log, "", id_room, id_moderator, id_target, log_type, log_message, date_log, list_unread)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func AddMember(id_room string, id_moderator string, id_target string, id_log string, date_log string, list_unread []string, log_type string, log_message string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list1 := []string{id_target}
			list2 := []string{id_moderator}
			list3 := []bson.D{
				{
					{Key: "list_id_member_moderator", Value: bson.D{
						{Key: "$in", Value: list2},
					}},
				},
				{
					{Key: "id_member_actor", Value: id_moderator},
				},
			}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: bson.D{
					{Key: "$ne", Value: id_target},
				}},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$nin", Value: list1},
				}},
				{Key: "$or", Value: list3},
			}
			update := bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "list_id_member", Value: id_target},
				}},
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_moderator", Value: id_target},
				}},
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_banned", Value: id_target},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
				CreateRoomActivity(id_log, "", id_room, id_moderator, id_target, log_type, log_message, date_log, list_unread)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func MemberBecomeModerator(id_room string, id_moderator string, id_target string, id_log string, date_log string, list_unread []string, log_type string, log_message string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list1 := []string{id_target}
			list2 := []string{id_moderator}
			list3 := []bson.D{
				{
					{Key: "list_id_member_moderator", Value: bson.D{
						{Key: "$in", Value: list2},
					}},
				},
				{
					{Key: "id_member_actor", Value: id_moderator},
				},
			}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: bson.D{
					{Key: "$ne", Value: id_target},
				}},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list1},
				}},
				{Key: "$or", Value: list3},
			}
			update := bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "list_id_member_moderator", Value: id_target},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
				CreateRoomActivity(id_log, "", id_room, id_moderator, id_target, log_type, log_message, date_log, list_unread)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func ModeratorBecomeMember(id_room string, id_moderator string, id_target string, id_log string, date_log string, list_unread []string, log_type string, log_message string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list1 := []string{id_target}
			list2 := []string{id_moderator}
			list3 := []bson.D{
				{
					{Key: "list_id_member_moderator", Value: bson.D{
						{Key: "$in", Value: list2},
					}},
				},
				{
					{Key: "id_member_actor", Value: id_moderator},
				},
			}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: bson.D{
					{Key: "$ne", Value: id_target},
				}},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list1},
				}},
				{Key: "$or", Value: list3},
			}
			update := bson.D{
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_moderator", Value: id_target},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
				CreateRoomActivity(id_log, "", id_room, id_moderator, id_target, log_type, log_message, date_log, list_unread)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func EnableRoomNotification(id_room string, id_actor string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list := []string{id_actor}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list},
				}},
				{Key: "list_id_member_banned", Value: bson.D{
					{Key: "$nin", Value: list},
				}},
			}
			update := bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "list_id_member_enable_notification", Value: id_actor},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func DisableRoomNotification(id_room string, id_actor string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list := []string{id_actor}
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list},
				}},
				{Key: "list_id_member_banned", Value: bson.D{
					{Key: "$nin", Value: list},
				}},
			}
			update := bson.D{
				{Key: "$pull", Value: bson.D{
					{Key: "list_id_member_enable_notification", Value: id_actor},
				}},
			}
			res, err := coll.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoom(id_room)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func DeleteRoom(id_room string, id_moderator string) bool {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	_, err3 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Delete()
			filter := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: id_moderator},
			}
			_, err1 := coll.DeleteOne(sessCtx, filter, opts)
			if err != nil {
				return nil, err1
			} else {
				DeleteRedisRoom(id_room)
			}
			coll2 := client.Database("chirpbird-v2").Collection("room_activities")
			opts2 := options.Delete()
			filter2 := bson.D{
				{Key: "id_room", Value: id_room},
			}
			res, err2 := coll2.DeleteMany(sessCtx, filter2, opts2)
			if err2 != nil {
				return nil, err2
			} else {
				hasil = true
			}
			return res, err2
		},
		txnOpts)
	if err3 != nil {
		return false
	}
	return hasil
}
