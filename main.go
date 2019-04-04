package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/user/login", loginIn)
	http.HandleFunc("/login_in", comeWabSite)
	http.HandleFunc("/login_out", loginOut)
	http.HandleFunc("/main", mainPage)
	fmt.Println("服务已启动在port:8000")
	http.ListenAndServe(":8000", nil)
}
