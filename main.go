package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//用于定时更新两个排行榜的数据
func updateList() {
	for {
		now := time.Now()
		hour, minute, second := now.Clock()
		if hour == 0 && minute == 0 && second == 0 {
			clearTable()
			pm := getMostListen()
			nm := getNewMusic()
			insertID("popular", pm)
			insertID("new", nm)
			log.Println("更新排行榜了")
		}
	}
}

func main() {
	// mux := http.NewServeMux()
	clearTable()
	pm := getMostListen()
	nm := getNewMusic()
	insertID("popular", pm)
	insertID("new", nm)
	log.Println("更新排行榜了")
	go updateList()
	http.HandleFunc("/user/login", loginIn)
	http.HandleFunc("/login_in", comeWabSite)
	http.HandleFunc("/login_out", loginOut)
	http.HandleFunc("/main", mainPage)
	http.HandleFunc("/song/upload", uploadFile)
	http.HandleFunc("/song/uploadPage", toUpload)
	http.HandleFunc("/search", searchSong)
	http.HandleFunc("/song", getSong)
	http.HandleFunc("/lrc", getLrc)
	http.HandleFunc("/uploadApp", dealAppUpdate)
	http.HandleFunc("/checkUpdate", checkUpdate)
	http.HandleFunc("/downloadApp", downloadApp)
	http.HandleFunc("/song/popular", returnPopularMusic)
	http.HandleFunc("/song/new", returnNewMusic)
	fmt.Println("服务已启动在port:8000")
	http.ListenAndServe(":8000", nil)
	// http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)
}