package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:,.Rfb8848@/mydb")
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(45)
	db.SetMaxOpenConns(90)
}

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
	// clearTable()
	// pm := getMostListen()
	// nm := getNewMusic()
	// insertID("popular", pm)
	// insertID("new", nm)
	// log.Println("更新排行榜了")
	// go updateList()
	http.HandleFunc("/user/login", loginIn)
	http.HandleFunc("/login_in", comeWabSite)
	http.HandleFunc("/login_out", loginOut)
	http.HandleFunc("/main", mainPage)
	http.HandleFunc("/song/upload", uploadFile)
	http.HandleFunc("/song/uploadPage", toUpload)
	http.HandleFunc("/search", searchSong)
	http.HandleFunc("/searchBy", searchSongByAllWord)
	http.HandleFunc("/song", getSong)
	http.HandleFunc("/lrc", getLrc)
	http.HandleFunc("/uploadApp", dealAppUpdate)
	http.HandleFunc("/checkUpdate", checkUpdate)
	http.HandleFunc("/downloadApp", downloadApp)
	http.HandleFunc("/song/popular", returnPopularMusic)
	http.HandleFunc("/song/new", returnNewMusic)
	http.HandleFunc("/music/user/register", userRegister)
	http.HandleFunc("/music/user/activation", userActivation)
	http.HandleFunc("/music/user/login", userLogin)
	http.HandleFunc("/music/user/getSongs", syncGetMusic)
	http.HandleFunc("/music/user/syncAddMusic", syncAddMusic)
	http.HandleFunc("/music/user/syncDelMusic", syncDelMusic)
	http.HandleFunc("/music/user/syncDelMusicById", syncDelMusicByID)
	http.HandleFunc("/music/user/addSongList", httpAddSongList)
	http.HandleFunc("/music/user/deleteSongList", httpDeleteSongList)
	fmt.Println("服务已启动在port:8000")
	http.ListenAndServe(":8000", nil)
	// http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)
}
