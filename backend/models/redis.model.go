package models

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx1 = context.Background()

func RedisDelete(key string) {
	rdb := GetRedisModel()
	defer rdb.Close()
	err := rdb.Del(ctx1, key).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func RedisSet(key string, val string) {
	rdb := GetRedisModel()
	defer rdb.Close()
	err := rdb.Set(ctx1, key, val, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func RedisGet(key string) (string, bool) {
	rdb := GetRedisModel()
	defer rdb.Close()
	val, err := rdb.Get(ctx1, key).Result()
	var status bool
	var hasil string
	if err == redis.Nil {
		status = false
	} else if err != nil {
		status = false
	} else {
		status = true
	}
	hasil = val
	return hasil, status
}

func RedisRead(key string, tableDB string, kolomDB string, idDB string) string {
	rawRedis, checkRedis := RedisGet(key)
	if checkRedis {
		return rawRedis
	} else {
		rawMongodb, checkMongodb := MongodbRead(tableDB, kolomDB, idDB)
		if checkMongodb {
			final, _ := json.Marshal(rawMongodb)
			temp := base64.StdEncoding.EncodeToString(final)
			RedisSet(key, temp)
			return temp
		} else {
			return ""
		}
	}
}

func RedisReadMember(idDB string) string {
	key := "M_" + idDB
	rawData := RedisRead(key, "members", "email", idDB)
	return rawData
}

func RedisReadRoom(idDB string) string {
	key := "G_" + idDB
	rawData := RedisRead(key, "rooms", "id_primary", idDB)
	return rawData
}

func RedisReadRoomActivity(idDB string) string {
	key := "RA_" + idDB
	rawData := RedisRead(key, "room_activities", "id_primary", idDB)
	return rawData
}

func DeleteRedisMember(idDB string) {
	key := "M_" + idDB
	RedisDelete(key)
}

func DeleteRedisRoom(idDB string) {
	key := "G_" + idDB
	RedisDelete(key)
}

func DeleteRedisRoomActivity(idDB string) {
	key := "RA_" + idDB
	RedisDelete(key)
}
