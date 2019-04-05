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
	sql := "insert into music(song_name,song_author,all_time,song_size,url,count) values(?,?,?,?,?,?)"
	state, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return errors.New("sql 语句有语法错误在prepare(sql)")
	}
	row, err := state.Exec(m.SongName, m.SongAuthor, m.AllTime, m.SongSize, m.URL, m.Count)
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

//歌曲搜索
func findMusicByWord(word string) []music {
	var musicList []music
	db := getDB()
	defer db.Close()
	sql := "select *from music where url like concat('%',?,'%')"
	state, _ := db.Prepare(sql)
	row, err := state.Query(word)
	if err != nil {
		log.Panicln(err)
	}
	for row.Next() {
		var mu music
		row.Scan(&mu.ID, &mu.SongName, &mu.SongAuthor, &mu.AllTime, &mu.SongSize, &mu.URL, &mu.Count)
		musicList = append(musicList, mu)
	}
	return musicList
}

//id搜索歌曲
func findMusicById(id string) string {
	var url string
	db := getDB()
	defer db.Close()
	sql := "select url from music where id=?"
	state, _ := db.Prepare(sql)
	row, _ := state.Query(id)
	for row.Next() {
		row.Scan(&url)
	}
	return url
}

//返回最大的版本号和安装包文件名
func findMaxVCode() (int, string, string) {
	db := getDB()
	defer db.Close()
	var path string
	var content string
	var vCode int
	sql := "select max(v_code) from version"
	row, _ := db.Query(sql)
	for row.Next() {
		row.Scan(&vCode)
	}
	sqls := "select content,name from version where v_code = ?"
	state, _ := db.Prepare(sqls)
	rows, _ := state.Query(vCode)
	for rows.Next() {
		rows.Scan(&content, &path)
		break
	}
	return vCode, content, path
}

//把安装包文件名和版本号添加到数据库当中
func insertApp(vCode string, content string, name string) {
	db := getDB()
	defer db.Close()
	if vCode != "" && name != "" {
		sql := "insert into version(v_code,content,name) values(?,?,?)"
		state, _ := db.Prepare(sql)
		state.Exec(vCode, content, name)
	}
}

//把歌曲的收听次数加1,param id
func addCount(id int) {
	db := getDB()
	defer db.Close()
	sql := "update music set count = count+1 where id =?"
	state, _ := db.Prepare(sql)
	state.Exec(id)
}

//清空所有歌曲的count值
func clearCount() {
	db := getDB()
	defer db.Close()
	sql := "update music set count=0"
	db.Exec(sql)
}
