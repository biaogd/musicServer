package main

import (
	"time"
)

func redisAddSongs(key string, value []byte) {
	err := client.Set(key, value, 0).Err()
	checkErr(err)
}

func getSongs(key string) []byte {
	value, err := client.Get(key).Bytes()
	checkErr(err)
	return value
}

func redisSet(key, value string) {
	err := client.Set(key, value, time.Minute*30).Err()
	checkErr(err)
}

func redisGet(key string) string {
	value := client.Get(key).Val()
	return value
}
