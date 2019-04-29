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
func checkLogin(email string, password string) bool {
	db := getDB()
	defer db.Close()
	sql := "select *from users where email=? and password=?"
	state, err := db.Prepare(sql)
	checkErr(err)
	result, err := state.Query(email, password)
	if result.Next() {
		return true
	}
	return false
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
