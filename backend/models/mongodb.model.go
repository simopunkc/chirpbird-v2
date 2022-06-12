package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongodbRead(collection string, kolom string, id_primary string) (interface{}, bool) {
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
	res, err := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection(collection)
			opts := options.FindOne().SetSort(bson.D{{Key: kolom, Value: 1}})
			var result bson.M
			err := coll.FindOne(ctx2, bson.D{
				{Key: kolom, Value: id_primary},
			}, opts).Decode(&result)
			if err != nil {
				return nil, err
			} else {
				hasil = true
			}
			return result, err
		},
		txnOpts)
	if err != nil {
		return err, false
	}
	return res, hasil
}
