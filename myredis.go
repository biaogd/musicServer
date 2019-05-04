package main

func redisAddSongs(key, value string) {
	err := client.Set(key, value, 0).Err()
	checkErr(err)
}

func getSongs(key string) string {
	value, err := client.Get(key).Result()
	checkErr(err)
	return value
}
