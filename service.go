package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func transform(song string, size int64) music {
	var mu music
	strs := strings.Split(song, "-")
	if len(strs) > 1 {
		mu.SongName = strings.TrimSpace(strings.Split(strs[1], ".")[0])
		mu.SongAuthor = strings.TrimSpace(strs[0])
		mu.AllTime = 0
		mu.SongSize = size
		mu.URL = song
		mu.Count = 0
	}
	return mu
}

func getVerifyCode() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var num string
	buf := bytes.NewBufferString(num)
	for i := 0; i < 6; i++ {
		n := r.Intn(10)
		buf.WriteString(strconv.Itoa(n))
	}
	return buf.String()
}
