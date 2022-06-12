package models

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetListRoomActivity(id_room string, page string, id_user string) (interface{}, bool) {
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
			var results []bson.M
			coll2 := client.Database("chirpbird-v2").Collection("room_activities")
			var itemPerPage int64 = 1000
			pages, _ := strconv.Atoi(page)
			var start int64 = 0
			if pages > 1 {
				start = int64(pages)*itemPerPage - itemPerPage
			}
			opts2 := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetSkip(start).SetLimit(itemPerPage)
			filter := bson.D{
				{Key: "id_room", Value: id_room},
			}
			cursor, _ := coll2.Find(ctx2, filter, opts2)
			if err = cursor.All(ctx2, &results); err != nil {
				return nil, err
			} else {
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
