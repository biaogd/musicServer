package session

import (
	"errors"
	"log"
	"time"
)

// Session : 是实体
type Session struct {
	Value    interface{}
	KeepTime int
}

// SessionManager : 用于管理保存session的map
type SessionManager struct {
	values map[string]Session
}

//Manager : 管理的借口
type Manager interface {
	Set(key string, value Session)
	Get(key string) (Session, error)
	Remove(key string)
	Contains(key string) bool
}

var maps = make(map[string]Session)
var manager = &SessionManager{maps}

func GetManager() *SessionManager {
	if manager == nil {
		manager = &SessionManager{maps}
	}
	return manager
}

//Contains : 检查是不是包含这个session
func (sm *SessionManager) Contains(key string) bool {
	if _, ok := sm.values[key]; ok == true {
		return true
	}
	return false
}

//Set : 设置一个session
func (sm *SessionManager) Set(key string, value Session) {
	if !sm.Contains(key) {
		sm.values[key] = value
		//定时器，session到时间会自动移除,单位是秒
		time.AfterFunc(time.Duration(value.KeepTime)*time.Second, func() {
			sm.Remove(key)
			log.Println("user已被移除")
		})
		log.Println("添加" + key + "完成")
	} else {
		log.Println("添加" + key + "失败")
	}

}

//Get : 获取到一个Session
func (sm *SessionManager) Get(key string) (Session, error) {
	if sm.Contains(key) {
		return sm.values[key], nil
	}
	return Session{}, errors.New("没有这个session")
}

//Remove : 删除对话
func (sm *SessionManager) Remove(key string) bool {
	if sm.Contains(key) {
		delete(sm.values, key)
		return true
	}
	return false
}

func begin() {
	for {

	}
}
