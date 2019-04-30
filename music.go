package main

import "fmt"

type music struct {
	ID         int
	SongName   string
	SongAuthor string
	AllTime    int
	SongSize   int64
	URL        string
	Count      int
}

func (m music) String() string {
	return fmt.Sprintf("ID=%d,SongName=%s,SongAuthor=%s,AllTime=%d,SongSize=%d,URL=%s", m.ID, m.SongName, m.SongAuthor, m.AllTime, m.SongSize, m.URL)
}

type myApp struct {
	Content string
	Name    string
	Status  string
}

type songList struct {
	ID           int
	SongListName string
	Count        int
}

//用户登录后从服务器获取歌单的结构体
type selfSong struct {
	ID         int    //编号
	ListID     int    //歌单id
	SongID     int    //歌曲的id
	SongName   string //歌曲名
	SongAuthor string //歌曲作者
	SongPath   string //歌曲的播放地址
}
