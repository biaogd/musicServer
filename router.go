package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"musicServer/session"
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

//根据关键字搜索音乐并且返回json,使用模糊搜索
func searchSong(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	words := r.Form["word"]
	var word string
	if len(words) > 0 {
		word = words[0]
	}
	word = strings.TrimSpace(word)
	mList := findMusicByWord(word, 0)
	bytes, _ := json.Marshal(mList)
	w.Write(bytes)
}

//根据关键字搜索音乐，使用mysql全文检索
func searchSongByAllWord(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	words := r.Form["word"]
	var word string
	if len(words) > 0 {
		word = words[0]
	}
	word = strings.TrimSpace(word)
	var mList []music
	if (strings.Count(word, "") - 1) == 1 {
		mList = findMusicByWord(word, 0)
	} else {
		mList = findMusicByWord(word, 1)
	}
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
		fi, _ := os.Stat(path)
		w.Header().Set("Content-Disposition", "attachment;filename="+fi.Name())
		w.Header().Set("Content-Length", string(fi.Size()))
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
		fi, _ := os.Stat(path)
		w.Header().Set("Content-Disposition", "attachment;filename="+fi.Name())
		w.Header().Set("Content-Length", string(fi.Size()))
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
		w.Header().Set("Content-Disposition", "attachment;filename="+name)
		fi, _ := os.Stat(apkdir + name)
		w.Header().Set("Content-Length", string(fi.Size()))
		http.ServeFile(w, r, apkdir+name)
	}
}

//用户注册处理
func userRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	v := r.PostForm
	userName := v["userName"][0]
	email := v["email"][0]
	pw := v["pw"][0]
	// pw2 := v["pw2"][0]
	//查找数据库，查看这个邮箱是否已经被注册，被注册提示被注册，
	//否则把信息加入到数据库中，提示根据邮件去激活这个账号
	if findUser(email) {
		//这个邮箱已经被注册
		w.Write([]byte("failed"))
	} else {
		//没有被注册，发送邮件提示激活账号，返回注册成功
		content := "<p>您已经注册成功，要正常使用，请先激活您的账号</p><br>" +
			"<a href='http://www.mybiao.top:8000/music/user/activation?email=" + email + "'>点击链接激活账号</a>"
		go mySendMail(email, content)
		addUser(user{0, userName, email, pw})
		w.Write([]byte("successed"))
	}
}

//用户激活,并进行初始化
func userActivation(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form["email"][0]
	var out string
	if email != "" {
		if j := checkFlag(email); j == 1 {
			//如果已经激活了
			out = "该用户已经激活，无需反复激活"
		} else {
			i := activation(email)
			if i > 0 {
				id := searchIDByEmail(email)
				if id != 0 {
					if i := insertSongList(id, "我喜欢"); i > 0 {
						out = "账号:" + email + " 激活成功"
					} else {
						out = "激活失败"
					}
				} else {
					out = "激活失败"
				}
			} else {
				out = "激活失败"
			}
		}
	} else {
		out = "激活失败"
	}
	w.Write([]byte(out))
}

//处理用户登录
//unActivation 账号已注册，但未激活
//success 登录成功
//noPassword 密码错误
//unRegister 用户未注册
func userLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.PostForm
	email := value["email"][0]
	password := value["password"][0]
	var status string
	if findUser(email) { //该用户已经注册了
		log.Println("有用户登录")
		if flag := checkFlag(email); flag == 0 { //flag等于0，未激活
			status = "unActivation"
		} else {
			if u := checkLogin(email, password); u.UserName != "" {
				list := selectSongList(u.ID)
				status = "success" + "-" + strconv.Itoa(u.ID) + "-" + u.UserName
				for i, v := range list {
					str := strconv.Itoa(v.ID) + "*" + v.SongListName + "*" + strconv.Itoa(v.Count)
					if i == 0 {
						status = status + "-" + str
						continue
					}
					status = status + ";" + str
				}
			} else {
				status = "noPassword"
			}
		}

	} else {
		status = "unRegister"
	}
	log.Println(status)
	w.Write([]byte(status))
}

//用户歌单，歌曲同步
func syncAddMusic(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.PostForm
	// userID := value["user_id"][0]
	songListID := value["song_list_id"][0]
	musicID := value["music_id"][0]
	musicName := value["music_name"][0]
	musicAuthor := value["music_author"][0]
	musicPath := value["music_path"][0]
	log.Println("得到用户歌曲上传请求")
	count := insertListMusic(songListID, musicID, musicName, musicAuthor, musicPath)
	if count > 0 {
		id, _ := strconv.Atoi(songListID)
		updateSongCount(id, 1)
		log.Println("插入成功")
	} else {
		log.Println("插入失败")
	}
}

//通过歌单id获取这个歌单的所有歌曲
func syncGetMusic(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := r.Form["list_id"][0]
	id, _ := strconv.Atoi(ids)
	if id != 0 {
		songs := selectListBySongListID(id)
		js, err := json.Marshal(songs)
		checkErr(err)
		w.Write(js)
	}
}

//从歌单中删除这个歌曲
func syncDelMusic(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.PostForm
	listID := value["listId"][0]
	songName := value["songName"][0]
	songAuthor := value["songAuthor"][0]
	id, _ := strconv.Atoi(listID)
	if i := deleteSongFromList(id, songName, songAuthor); i > 0 {
		updateSongCount(id, -1) //歌单歌曲数量减1
		log.Println(listID, songName, songAuthor, "删除成功")
	} else {
		log.Println(listID, songName, songAuthor, "删除失败")
	}
}

//通过歌单歌曲id删除歌曲
func syncDelMusicByID(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ids := r.Form["songId"][0]
	id, _ := strconv.Atoi(ids)
	if i := deleteSongByID(id); i > 0 {
		log.Println("删除成功")
		updateSongCount(id, -1)
	} else {
		log.Println("删除失败")
	}
}

//处理新建歌单并且返回歌单信息
func httpAddSongList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.PostForm
	id := value["userId"][0]
	userID, err := strconv.Atoi(id)
	checkErr(err)
	listName := value["listName"][0]
	i := addSongList(userID, listName)
	var bytes []byte
	var sList songList
	if i > 0 {
		//歌单已加入到数据库
		listID := selectSongListID(userID, listName)
		sList = songList{listID, listName, 0}
		log.Println("歌单", listName, "添加成功")
	} else {
		log.Println("歌单", listName, "添加失败")
	}
	bytes, err = json.Marshal(sList)
	checkErr(err)
	w.Write(bytes)
}

//处理歌单的删除
func httpDeleteSongList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.PostForm
	listID, err := strconv.Atoi(value["listId"][0])
	checkErr(err)
	if i := deleteSongList(listID); i > 0 {
		if j := deleteAllByListID(listID); j > 0 {
			log.Println(listID, "歌单歌曲已删除", j)
		} else {
			log.Println(listID, "歌单歌曲删除失败")
		}
		log.Println(listID, "歌单已删除", i)
	} else {
		log.Println(listID, "歌单删除失败")
	}
	w.Write([]byte("success"))
}

//处理错误报告
func httpInsertError(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	value := r.PostForm
	userEmail := value["userEmail"][0]
	now := value["time"][0]
	text := value["text"][0]
	i := insertErrorReport(userEmail, text, now)
	if i > 0 {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("falied"))
	}
}

//使用redis缓存排行榜popular
func redisGetPopular(w http.ResponseWriter, r *http.Request) {
	pop := getSongs("popular")
	w.Write(pop)
}

//使用redis缓存排行榜new
func redisGetNew(w http.ResponseWriter, r *http.Request) {
	news := getSongs("new")
	w.Write(news)
}
