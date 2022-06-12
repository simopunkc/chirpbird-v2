package models

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateRoom(id_room string, id_member_creator string, name string, list_id_member []string, list_id_member_moderator []string, list_id_member_banned []string, list_id_member_enable_notification []string, date_created string, date_last_activity string, link_join string, id_log string, log_message string, log_type string) bool {
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
			opts := options.Replace().SetUpsert(true)
			filter := bson.D{{Key: "id_primary", Value: id_room}}
			replacement := bson.D{
				{Key: "id_primary", Value: id_room},
				{Key: "id_member_creator", Value: id_member_creator},
				{Key: "name", Value: name},
				{Key: "list_id_member", Value: list_id_member},
				{Key: "list_id_member_moderator", Value: list_id_member_moderator},
				{Key: "list_id_member_banned", Value: list_id_member_banned},
				{Key: "list_id_member_enable_notification", Value: list_id_member_enable_notification},
				{Key: "date_created", Value: date_created},
				{Key: "date_last_activity", Value: date_last_activity},
				{Key: "link_join", Value: link_join},
			}
			res, err := coll.ReplaceOne(sessCtx, filter, replacement, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				CreateRoomActivity(id_log, "", id_room, id_member_creator, "", log_type, log_message, date_created, []string{})
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}

func GetListRoom(page string, id string) (interface{}, bool) {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return err, false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	result, err := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			var itemPerPage int64 = 1000
			pages, _ := strconv.Atoi(page)
			var start int64 = 0
			if pages > 1 {
				start = int64(pages)*itemPerPage - itemPerPage
			}
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Find().SetSort(bson.D{{Key: "date_last_activity", Value: 1}}).SetSkip(start).SetLimit(itemPerPage)
			list := []string{id}
			filter := bson.D{
				{Key: "list_id_member", Value: bson.D{
					{Key: "$in", Value: list},
				}},
			}
			cursor, err := coll.Find(ctx2, filter, opts)
			if err != nil {
				return nil, err
			}
			var results []bson.M
			if err = cursor.All(ctx2, &results); err != nil {
				return nil, err
			}
			if len(results) > 0 {
				hasil = true
			}
			return results, err
		},
		txnOpts)
	if err != nil {
		return err, false
	}
	return result, hasil
}

func CreateRoomActivity(id_chat string, id_parent string, id_room string, id_member_actor string, id_member_target string, type_activity string, message string, date_created string, list_id_member_unread []string) bool {
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
	_, err1 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("room_activities")
			opts := options.Replace().SetUpsert(true)
			filter := bson.D{
				{Key: "id_primary", Value: id_chat},
			}
			replacement := bson.D{
				{Key: "id_primary", Value: id_chat},
				{Key: "id_parent", Value: id_parent},
				{Key: "id_room", Value: id_room},
				{Key: "id_member_actor", Value: id_member_actor},
				{Key: "id_member_target", Value: id_member_target},
				{Key: "type_activity", Value: type_activity},
				{Key: "message", Value: message},
				{Key: "date_created", Value: date_created},
				{Key: "list_id_member_unread", Value: list_id_member_unread},
			}
			res, err := coll.ReplaceOne(sessCtx, filter, replacement, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				UpdateRoomDateActivity(id_room, date_created)
			}
			return res, err
		},
		txnOpts)
	if err1 != nil {
		return false
	}
	return hasil
}

func CheckTokenGroupExist(token string) (interface{}, bool) {
	client := GetMongoModel()
	defer client.Disconnect(ctx2)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := client.StartSession(opts)
	hasil := false
	if err != nil {
		return err, false
	}
	defer sess.EndSession(ctx2)
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	res, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.FindOne()
			var result bson.M
			err := coll.FindOne(ctx2, bson.D{
				{Key: "link_join", Value: token},
			}, opts).Decode(&result)
			if err != nil {
				return nil, err
			} else {
				hasil = true
			}
			return result, err
		},
		txnOpts)
	if err2 != nil {
		return err2, false
	}
	return res, hasil
}

func JoinRoom(id_group string, user string, id_log string, date_log string, list_unread []string, log_message string, log_type string) bool {
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
	_, err1 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			list := []string{user}
			filter := bson.D{
				{Key: "id_primary", Value: id_group},
				{Key: "list_id_member", Value: bson.D{
					{Key: "$nin", Value: list},
				}},
				{Key: "list_id_member_banned", Value: bson.D{
					{Key: "$nin", Value: list},
				}},
			}
			update := bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "list_id_member", Value: user},
				}},
			}
			res, err2 := coll.UpdateOne(sessCtx, filter, update, opts)
			if err2 != nil {
				return nil, err2
			} else {
				hasil = true
				DeleteRedisRoom(id_group)
				CreateRoomActivity(id_log, "", id_group, user, "", log_type, log_message, date_log, list_unread)
			}
			return res, err2
		},
		txnOpts)
	if err1 != nil {
		return false
	}
	return hasil
}

func UpdateRoomDateActivity(id_group string, date string) bool {
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
	_, err1 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("rooms")
			opts := options.Update().SetUpsert(false)
			filter := bson.D{
				{Key: "id_primary", Value: id_group},
			}
			update := bson.D{
				{Key: "date_last_activity", Value: date},
			}
			res, err2 := coll.UpdateOne(sessCtx, filter, update, opts)
			if err2 != nil {
				return nil, err2
			} else {
				hasil = true
				DeleteRedisRoom(id_group)
			}
			return res, err2
		},
		txnOpts)
	if err1 != nil {
		return false
	}
	return hasil
}
