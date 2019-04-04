package main

import "strings"

func transform(song string, size int64) music {
	var mu music
	strs := strings.Split(song, "-")
	if len(strs) > 1 {
		mu.SongName = strings.TrimSpace(strings.Split(strs[1], ".")[0])
		mu.SongAuthor = strings.TrimSpace(strs[0])
		mu.AllTime = 0
		mu.SongSize = size
		mu.URL = song
	}
	return mu
}
