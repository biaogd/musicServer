package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"music/session"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const (
	musicDir = "/home/admin/mymusic/"
	apkdir   = "/home/admin/apkdir/"
)

var man = session.GetManager()

func loginIn(w http.ResponseWriter, r *http.Request) {
	if man.Contains("user") {
		http.Redirect(w, r, "/main", http.StatusFound)
	}
	log.Println("函数loginIn执行了")
	//使用模板解析
	t := template.Must(template.ParseFiles("static/login.html"))
	t.Execute(w, nil)
	// http.ServeFile(w, r, "static/login.html")
}
func mainPage(w http.ResponseWriter, r *http.Request) {
	many := make(map[string]interface{})
	many["size"] = getMusicCount()
	vCode, _, name := findMaxVCode()
	many["vCode"] = vCode
	many["name"] = name
	if man.Contains("user") {
		log.Println("函数mainPage执行了")
		t := template.Must(template.ParseFiles("static/main.html"))
		t.Execute(w, many)
		// http.ServeFile(w, r, "static/main.html")
	} else {
		http.Redirect(w, r, "/user/login", http.StatusFound)
	}
}

//登陆网站，输入密码，正确后把user用户保存到session中
func comeWabSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	if !man.Contains("user") {
		pw := r.PostForm["password"][0]
		if pw == "123456" {
			man.Set("user", session.Session{pw, 60 * 30})
			log.Println("登陆成功")
			http.Redirect(w, r, "/main", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/main", http.StatusFound)
	}
}

//退出登陆，删除保存的session
func loginOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if man.Contains("user") {
		man.Remove("user")
	}
}

//用于处理歌曲和歌词的上传
func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	reqSongFile, hander, err := r.FormFile("songFile")
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, "上传文件失败", err)
	}
	reqLrcFile, hander1, err := r.FormFile("lrcFile")
	if err != nil {
		println(err)
		fmt.Fprintln(w, "上传文件失败", err)
	}
	//这首歌没有被加入到数据库当中
	if findMusicBySongName(hander.Filename) == 0 {

		songFile, err := os.Create(musicDir + hander.Filename)
		if err != nil {
			println(err)
		}
		io.Copy(songFile, reqSongFile)
		lrcFile, err := os.Create(musicDir + hander1.Filename)
		if err != nil {
			println(err)
		}
		io.Copy(lrcFile, reqLrcFile)
		log.Println("歌曲和歌词文件已写入")
		m := transform(hander.Filename, hander.Size)
		insertMusic(m)
		log.Println("歌曲信息已插入到数据库当中")
		fmt.Fprintln(w, "上传并插入到数据库成功")
	} else {
		fmt.Fprintln(w, "该歌曲已存在")
	}

}

func toUpload(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("static/upload.html"))
	t.Execute(w, nil)
}

//根据关键字搜索音乐并且返回json
func searchSong(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	words := r.Form["word"]
	var word string
	if len(words) > 0 {
		word = words[0]
	}
	mList := findMusicByWord(word)
	bytes, _ := json.Marshal(mList)
	w.Write(bytes)
}

//根据歌曲id下载歌曲
func getSong(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]
	if id != "" {
		url := findMusicById(id)
		path := musicDir + url
		//收听的次数加1
		ids, _ := strconv.Atoi(id)
		addCount(ids)
		http.ServeFile(w, r, path)
	}
}

//根据id下载歌词
func getLrc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]
	if id != "" {
		url := findMusicById(id)
		str := strings.Split(url, ".")[0] + ".lrc"
		path := musicDir + str
		http.ServeFile(w, r, path)
	}
}

//处理软件上传的
func dealAppUpdate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100 * 1024 * 1024)
	file, hander, _ := r.FormFile("updateApp")
	apkFile, err := os.Create(apkdir + hander.Filename)
	if err != nil {
		log.Panicln(err)
	}
	io.Copy(apkFile, file)
	code := r.PostForm["vCode"][0]
	content := r.PostForm["content"][0]
	insertApp(code, content, hander.Filename)
	fmt.Fprintln(w, "上传成功")
}

//检查更新,传入版本号
func checkUpdate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vCode := r.Form["code"][0]
	vCode1, content, name := findMaxVCode()
	vCode2, _ := strconv.Atoi(vCode)
	if vCode1 > vCode2 {
		var app myApp
		app.Name = name
		app.Content = content
		app.Status = "ok"
		bytes, _ := json.Marshal(app)
		w.Write(bytes)
	} else {
		var app myApp
		app.Status = "no"
		bytes, _ := json.Marshal(app)
		w.Write(bytes)
	}
}

//直接下载最新版软件包
func downloadApp(w http.ResponseWriter, r *http.Request) {
	_, _, name := findMaxVCode()
	if name != "" {
		http.ServeFile(w, r, apkdir+name)
	}
}

//返回热歌20个
func returnPopularMusic(w http.ResponseWriter, r *http.Request) {
	mList := getPopularByIds()
	js, err := json.Marshal(mList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(js)
}

//返回最新的20个歌曲
func returnNewMusic(w http.ResponseWriter, r *http.Request) {
	mList := getNewByIds()
	js, err := json.Marshal(mList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(js)
}
