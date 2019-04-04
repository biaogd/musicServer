package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//获得连接对象
func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:,.Rfb8848@/mydb")
	if err != nil {
		log.Panicln(err)
	}
	return db
}

//插入单个歌曲信息
func insertMusic(m music) error {
	db := getDB()
	defer db.Close()
	sql := "insert into music(song_name,song_author,all_time,song_size,url) values(?,?,?,?,?)"
	state, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return errors.New("sql 语句有语法错误在prepare(sql)")
	}
	row, err := state.Exec(m.SongName, m.SongAuthor, m.AllTime, m.SongSize, m.URL)
	if err != nil {
		log.Println(err)
		return errors.New("插入歌曲信息错误")
	}
	i, _ := row.RowsAffected()
	log.Println("插入了" + string(i) + "条")
	return nil
}

//获取歌曲的总数目
func getMusicCount() int {
	var size int
	db := getDB()
	defer db.Close()
	sql := "select count(*) as size from music"
	res, err := db.Query(sql)
	if err != nil {
		log.Println(err)
	}
	for res.Next() {
		res.Scan(&size)
	}
	log.Printf("歌曲的总数,%d", size)
	return size
}

//查找歌曲，根据歌曲名
//返回值,查找到的个数
func findMusicBySongName(name string) int {
	var size int
	db := getDB()
	defer db.Close()
	sql := "select count(*) as size from music where url=?"
	state, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
	}
	res, _ := state.Query(name)
	for res.Next() {
		res.Scan(&size)
	}
	log.Printf("查找到歌曲的总数,%d", size)
	return size
}
