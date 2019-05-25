package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var client *redis.Client

func init() {
	var err error
	db, err = sql.Open("mysql", "root:,.Rfb8848@/mydb")
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(200)
	db.SetMaxOpenConns(500)

	//初始化redis连接对象
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

//用于定时更新两个排行榜的数据
func updateList() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(next))
		<-t.C
		pm := getMostListen()
		nm := getNewMusic()
		pmjson, _ := json.Marshal(pm)
		redisAddSongs("popular", pmjson)
		nmjson, _ := json.Marshal(nm)
		redisAddSongs("new", nmjson)
	}
}

func main() {
	// clearTable()
	pm := getMostListen()
	nm := getNewMusic()
	pmjson, _ := json.Marshal(pm)
	redisAddSongs("popular", pmjson)
	nmjson, _ := json.Marshal(nm)
	redisAddSongs("new", nmjson)
	// insertID("popular", pm)
	// insertID("new", nm)
	// log.Println("更新排行榜了")
	go updateList()
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
	http.HandleFunc("/song/popular", redisGetPopular)
	http.HandleFunc("/song/new", redisGetNew)
	http.HandleFunc("/music/user/register", userRegister)
	http.HandleFunc("/music/user/activation", userActivation)
	http.HandleFunc("/music/user/login", userLogin)
	http.HandleFunc("/music/user/getSongs", syncGetMusic)
	http.HandleFunc("/music/user/syncAddMusic", syncAddMusic)
	http.HandleFunc("/music/user/syncDelMusic", syncDelMusic)
	http.HandleFunc("/music/user/syncDelMusicById", syncDelMusicByID)
	http.HandleFunc("/music/user/addSongList", httpAddSongList)
	http.HandleFunc("/music/user/deleteSongList", httpDeleteSongList)
	http.HandleFunc("/errorReport", httpInsertError)
	fmt.Println("服务已启动在port:8000")
	http.ListenAndServe(":8000", nil)
}
