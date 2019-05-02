package main

import (
	"database/sql"
	"errors"
	"fmt"
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

//歌曲搜索,使用模糊搜索
func findMusicByWord(word string, like int) []music {
	var musicList []music
	db := getDB()
	defer db.Close()
	var sql, param string
	//模糊搜索
	if like == 0 {
		sql = "select *from music where url like concat('%',?,'%')"
		param = word
	} else {
		//中文全文检索，这是默认的自然语言检索方式，还有boolean模式
		sql = "select *from music where match(url) against(?)"
		param = "*" + word + "*"
	}
	state, _ := db.Prepare(sql)
	row, err := state.Query(param)
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

//返回歌曲收听次数的前20名
func getMostListen() []music {
	var mlist []music
	db := getDB()
	defer db.Close()
	sql := "select *from music order by count desc limit 20"
	rows, err := db.Query(sql)
	if err != nil {
		println(err)
	}
	for rows.Next() {
		var m music
		rows.Scan(&m.ID, &m.SongName, &m.SongName, &m.AllTime, &m.SongSize, &m.URL, &m.Count)
		mlist = append(mlist, m)
	}
	return mlist
}

//返回歌曲最新上传的前20个
func getNewMusic() []music {
	var mList []music
	db := getDB()
	defer db.Close()
	sql := "select *from music order by id desc limit 20"
	rows, err := db.Query(sql)
	if err != nil {
		println(err)
	}
	for rows.Next() {
		var m music
		rows.Scan(&m.ID, &m.SongName, &m.SongName, &m.AllTime, &m.SongSize, &m.URL, &m.Count)
		mList = append(mList, m)
	}
	return mList

}

//通过id查找歌曲，返回一个music
func selectMusicById(id int) music {
	var m music
	db := getDB()
	defer db.Close()
	sql := "select *from music where id=?"
	state, _ := db.Prepare(sql)
	rows, err := state.Query(id)
	if err != nil {
		println(err)
	}
	for rows.Next() {
		rows.Scan(&m.ID, &m.SongName, &m.SongName, &m.AllTime, &m.SongSize, &m.URL, &m.Count)
		break
	}
	return m
}

//清空两个排行榜的数据库
func clearTable() {
	db := getDB()
	defer db.Close()
	sql1 := "delete from popular"
	sql2 := "delete from new"
	db.Exec(sql1)
	db.Exec(sql2)
}

//把最新的排行列表插入到两个表中
func insertID(table string, mList []music) {
	db := getDB()
	defer db.Close()
	var sql string
	if table == "popular" {
		sql = "insert into popular(song_id) values(?)"
		state, err := db.Prepare(sql)
		if err != nil {
			println(err)
		}
		for _, mu := range mList {
			state.Exec(mu.ID)
		}
	} else {
		sql = "insert into new(song_id) values(?)"
		state, err := db.Prepare(sql)
		if err != nil {
			println(err)
		}
		for _, mu := range mList {
			state.Exec(mu.ID)
		}
	}

}

//先从popular表中找到歌曲id,在从歌曲表中返回歌曲信息
func getPopularByIds() []music {
	db := getDB()
	defer db.Close()
	var ids []int
	sql := "select song_id from popular"
	rows, _ := db.Query(sql)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	var mList []music
	sqls := "select *from music where id=?"
	state, _ := db.Prepare(sqls)
	for _, id := range ids {
		var mu music
		rows, _ := state.Query(id)
		for rows.Next() {
			rows.Scan(&mu.ID, &mu.SongName, &mu.SongAuthor, &mu.AllTime, &mu.SongSize, &mu.URL, &mu.Count)
			break
		}
		mList = append(mList, mu)
	}
	return mList
}

func getNewByIds() []music {
	db := getDB()
	defer db.Close()
	var ids []int
	sql := "select song_id from new"
	rows, _ := db.Query(sql)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	var mList []music
	sqls := "select *from music where id=?"
	state, _ := db.Prepare(sqls)
	for _, id := range ids {
		var mu music
		rows, _ := state.Query(id)
		for rows.Next() {
			rows.Scan(&mu.ID, &mu.SongName, &mu.SongAuthor, &mu.AllTime, &mu.SongSize, &mu.URL, &mu.Count)
			break
		}
		mList = append(mList, mu)
	}
	return mList
}

//查找邮箱是否已经被注册
//return true when
func findUser(email string) bool {
	db := getDB()
	defer db.Close()
	sql := "select *from users where email=?"
	state, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	result, err := state.Query(email)
	if err != nil {
		panic(err)
	}
	if result.Next() {
		return true
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

//添加一个用户
//返回影响的行数
func addUser(u user) int64 {
	db := getDB()
	defer db.Close()
	sql := "insert into users(user_name,email,password,flag) values(?,?,?,?)"
	state, err := db.Prepare(sql)
	checkErr(err)
	result, err := state.Exec(u.UserName, u.Email, u.Password, 0)
	checkErr(err)
	rows, err := result.RowsAffected()
	checkErr(err)
	return rows
}

//根据email查找用户id
func searchIDByEmail(email string) int {
	db := getDB()
	defer db.Close()
	sql := "select id from users where email=?"
	res, err := db.Query(sql, email)
	checkErr(err)
	var id int
	if res.Next() {
		res.Scan(&id)
	}
	return id
}

//用户添加歌单,返回影响的行数
func insertSongList(id int, name string) int64 {
	db := getDB()
	defer db.Close()
	sql := "insert into user_song(user_id,song_list_name,count) values(?,?,0)"
	res, err := db.Exec(sql, id, name)
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

//修改数据库flag,激活用户
func activation(email string) int64 {
	db := getDB()
	defer db.Close()
	sql := "update users set flag=1 where email=?"
	res, err := db.Exec(sql, email)
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

//检查用户名密码是否正确
func checkLogin(email string, password string) user {
	db := getDB()
	defer db.Close()
	sql := "select *from users where email=? and password=?"
	state, err := db.Prepare(sql)
	checkErr(err)
	result, err := state.Query(email, password)
	var u user
	var flag int
	if result.Next() {
		// return true
		result.Scan(&u.ID, &u.UserName, &u.Email, &u.Password, &flag)
		return u
	}
	return user{}
}

//检查是否激活
func checkFlag(email string) int {
	db := getDB()
	defer db.Close()
	sql := "select flag from users where email=?"
	res, err := db.Query(sql, email)
	checkErr(err)
	var flag int
	for {
		if res.Next() {
			res.Scan(&flag)
			break
		}
	}
	return flag
}

//插入一个歌单歌曲信息
//返回影响的行数
func insertListMusic(songListID, musicID, musicName, musicAuthor, musicPath string) int64 {
	db := getDB()
	defer db.Close()
	sql := "insert into list_music(song_list_id,music_id,music_name,music_author,music_path) values(?,?,?,?,?)"
	result, err := db.Exec(sql, songListID, musicID, musicName, musicAuthor, musicPath)
	checkErr(err)
	count, err := result.RowsAffected()
	checkErr(err)
	return count
}

//查找这个用户的所有歌单
func selectSongList(userID int) []songList {
	db := getDB()
	defer db.Close()
	sql := "select id,song_list_name,count from user_song where user_id=?"
	result, err := db.Query(sql, userID)
	checkErr(err)
	var songLists []songList
	for result.Next() {
		var list songList
		result.Scan(&list.ID, &list.SongListName, &list.Count)
		songLists = append(songLists, list)
	}
	return songLists
}

//查找用户的歌单id,通过用户id和歌单名称
//返回歌单的id
func selectSongListID(userID int, listName string) int {
	db := getDB()
	defer db.Close()
	var id int
	fmt.Println(userID, listName)
	sql := "select id from user_song where user_id=? and song_list_name=?"
	result, err := db.Query(sql, userID, listName)
	checkErr(err)
	if result.Next() {
		result.Scan(&id)
	}
	return id
}

//返回这个歌单所有歌曲信息
func selectListBySongListID(id int) []selfSong {
	db := getDB()
	defer db.Close()
	var songs []selfSong
	sql := "select *from list_music where song_list_id = ?"
	result, err := db.Query(sql, id)
	checkErr(err)
	for result.Next() {
		var song selfSong
		result.Scan(&song.ID, &song.ListID, &song.SongID, &song.SongName, &song.SongAuthor, &song.SongPath)
		songs = append(songs, song)
	}
	return songs
}

//更新歌单歌曲数目
func updateSongCount(listID int, count int) {
	db := getDB()
	defer db.Close()
	var sql string
	if count > 0 {
		sql = "update user_song set count=count+? where id=?"
	} else {
		sql = "update user_song set count=count-? where id=?"
	}
	if count < 0 {
		count = -count
	}
	db.Exec(sql, count, listID)
}

//删除id为listID的歌单的名称为songName,作者为songAuthor的歌曲
//return affect rows
func deleteSongFromList(listID int, songName, songAuthor string) int64 {
	db := getDB()
	defer db.Close()
	sql := "delete from list_music where song_list_id=? and music_name=? and music_author=?"
	res, err := db.Exec(sql, listID, songName, songAuthor)
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

//通过主键删除歌曲
//返回影响的行数
func deleteSongByID(songID int) int64 {
	db := getDB()
	defer db.Close()
	sql := "delete from list_music where id=?"
	res, err := db.Exec(sql, songID)
	checkErr(err)
	count, err := res.RowsAffected()
	checkErr(err)
	return count
}

//添加一个歌单
//参数 userId 用户的id;listName 歌单的名称
//返回值 影响的行数
func addSongList(userID int, listName string) int64 {
	db := getDB()
	defer db.Close()
	sql := "insert into user_song(user_id,song_list_name,count) values(?,?,?)"
	result, err := db.Exec(sql, userID, listName, 0)
	checkErr(err)
	count, err := result.RowsAffected()
	checkErr(err)
	return count
}
