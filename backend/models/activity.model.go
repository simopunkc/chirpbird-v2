package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DeleteRoomActivity(id_activity string, id_moderator string) bool {
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
			coll := client.Database("chirpbird-v2").Collection("room_activities")
			opts := options.Delete()
			filter := bson.D{
				{Key: "id_primary", Value: id_activity},
				{Key: "id_member_actor", Value: id_moderator},
			}
			res, err := coll.DeleteOne(sessCtx, filter, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisRoomActivity(id_activity)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}
