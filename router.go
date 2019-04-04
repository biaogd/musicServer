package main

import (
	"fmt"
	"io"
	"log"
	"music/session"
	"net/http"
	"os"
	"text/template"
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
			man.Set("user", session.Session{pw, 20})
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
	}
	reqLrcFile, hander1, err := r.FormFile("lrcFile")
	if err != nil {
		println(err)
	}
	//这首歌没有被加入到数据库当中
	if findMusicBySongName(hander.Filename) == 0 {

		songFile, err := os.Create("/home/biao/music/" + hander.Filename)
		if err != nil {
			println(err)
		}
		io.Copy(songFile, reqSongFile)
		lrcFile, err := os.Create("/home/biao/music/" + hander1.Filename)
		if err != nil {
			println(err)
		}
		io.Copy(lrcFile, reqLrcFile)
		log.Println("歌曲和歌词文件已写入")
		m := transform(hander.Filename, hander.Size)
		insertMusic(m)
		log.Println("歌曲信息已插入到数据库当中")
		fmt.Fprintln(w, "上传并插入到数据库成功")
	}

}

func toUpload(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("static/upload.html"))
	t.Execute(w, nil)
}
