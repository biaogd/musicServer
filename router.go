package main

import (
	"log"
	"music/session"
	"net/http"
	"text/template"
)

var man = session.GetManager()

func loginIn(w http.ResponseWriter, r *http.Request) {
	if man.Contains("user") {
		http.Redirect(w, r, "/main", http.StatusFound)
	}
	log.Println("函数loginIn执行了")
	t := template.Must(template.ParseFiles("static/login.html"))
	t.Execute(w, nil)
	// http.ServeFile(w, r, "static/login.html")
}
func mainPage(w http.ResponseWriter, r *http.Request) {
	if man.Contains("user") {
		log.Println("函数mainPage执行了")
		http.ServeFile(w, r, "static/main.html")
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
			man.Set("user", session.Session{pw, 10})
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
