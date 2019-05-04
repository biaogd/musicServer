package main

func redisAddSongs(key string, value []byte) {
	err := client.Set(key, value, 0).Err()
	checkErr(err)
}

func getSongs(key string) []byte {
	value, err := client.Get(key).Bytes()
	checkErr(err)
	return value
}
