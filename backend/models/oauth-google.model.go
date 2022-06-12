package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func SaveMember(email string, name string, picture string, verified_email bool) bool {
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
	_, err2 := sess.WithTransaction(
		ctx2,
		func(sessCtx mongo.SessionContext) (interface{}, error) {
			coll := client.Database("chirpbird-v2").Collection("members")
			opts := options.Replace().SetUpsert(true)
			filter := bson.D{{Key: "email", Value: email}}
			replacement := bson.D{
				{Key: "email", Value: email},
				{Key: "name", Value: name},
				{Key: "picture", Value: picture},
				{Key: "verified_email", Value: verified_email},
			}
			res, err := coll.ReplaceOne(sessCtx, filter, replacement, opts)
			if err != nil {
				return nil, err
			} else {
				hasil = true
				DeleteRedisMember(email)
			}
			return res, err
		},
		txnOpts)
	if err2 != nil {
		return false
	}
	return hasil
}
