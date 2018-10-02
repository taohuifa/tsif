package session

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"net/http"
)

const (
	SESSION_MAX_TIME     = int64(10 * time.Second) // session超时时间
	SESSION_UPDATE_COUNT = 10                      // session更新次数
	COOKIE_SESSION_NAME  = "SESSION"               // cookie Session储存名字
	COOKIE_MAX_TIME      = int(10 * time.Second)   // cookie超时时间
)

//----------------------session----------------------------

type Session struct {
	ID         string
	lock       sync.RWMutex                //一把互斥锁
	updateTime int64                       //最后访问时间
	caches     map[interface{}]interface{} //主数据
}

// 设置数据
func (this *Session) Set(key, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.caches[key] = value
}

// 读取数据
func (this *Session) Get(key interface{}, defVal interface{}) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if value := this.caches[key]; value != nil {
		return value
	}
	return defVal
}

// 删除数据
func (this *Session) Remove(key interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if value := this.caches[key]; value != nil {
		delete(this.caches, key)
	}
	return nil
}

// 获取ID
func (this *Session) GetID() string {
	return this.ID
}

//----------------------session manager----------------------------

type SessionManager struct {
	sessionID     uint64
	updateCounter uint64
	sessions      map[string]*Session
	lock          sync.RWMutex
}

func (this *SessionManager) newSessionID() string {
	sessionId := atomic.AddUint64(&this.sessionID, 1)
	sessionStr := fmt.Sprintf("%09d", sessionId)
	encStr := base64.URLEncoding.EncodeToString([]byte(sessionStr))
	return string(encStr)
}

func (this *SessionManager) New() (*Session, error) {
	sessionId := this.newSessionID()
	session := Session{ID: sessionId, updateTime: time.Now().UnixNano(), caches: make(map[interface{}]interface{})}
	// 上锁
	this.lock.Lock()
	defer this.lock.Unlock()
	defer this.deferByAction()
	// 设置session
	this.sessions[sessionId] = &session
	return &session, nil
}

func (this *SessionManager) Get(sessionId string) (*Session, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	defer this.deferByAction()
	// 读取session
	session, ok := this.sessions[sessionId]
	if !ok {
		return nil, errors.New("no exist")
	}
	session.updateTime = time.Now().UnixNano()
	return session, nil
}

func (this *SessionManager) Remove(sessionId string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.sessions, sessionId)
}

// get/set 更新操作
func (this *SessionManager) deferByAction() {
	go func() {
		updateCounter := atomic.AddUint64(&this.updateCounter, 1)
		if updateCounter%SESSION_UPDATE_COUNT != 0 {
			return
		}
		this.update()
	}()
}

// 更新session
func (this *SessionManager) update() error {
	this.lock.RLock()

	// 遍历
	nowtime := time.Now().UnixNano()
	removes := make([]*Session, 0)
	for _, session := range this.sessions {
		if session == nil {
			continue
		}
		dt := nowtime - session.updateTime
		// Log.Debug("session dt "+fmt.Sprint(dt), false)
		if dt <= SESSION_MAX_TIME {
			continue // 尚未超时
		}
		// 添加到移除列表
		removes = append(removes, session)
	}
	this.lock.RUnlock()

	// 执行移除
	removeNum := len(removes)
	if removeNum > 0 {
		for _, session := range removes {
			if session == nil {
				continue
			}
			// Log.Debug("remove session: "+fmt.Sprint(session.GID), false)
			this.Remove(session.ID)
		}
	}

	return nil
}

// 使用cookie处理session
func (this *SessionManager) Start(w http.ResponseWriter, r *http.Request) (*Session, error) {
	var session *Session
	// 获取储存的session
	cookie, err := r.Cookie(COOKIE_SESSION_NAME)
	if err == nil {
		sessionId := cookie.Value
		session, err = this.Get(sessionId)
		if err == nil {
			return session, nil
		}
	}
	// 新建session
	session, err = this.New()
	if err != nil {
		return nil, err // 错误
	}

	//让浏览器cookie设置过期时间
	cookie = &http.Cookie{Name: COOKIE_SESSION_NAME, Value: session.ID, Path: "/", HttpOnly: true, MaxAge: COOKIE_MAX_TIME}
	http.SetCookie(w, cookie)
	return session, nil
}

//--------------------instance-------------------------

var Instance SessionManager = SessionManager{sessions: make(map[string]*Session)}

func New() (*Session, error) {
	return Instance.New()
}

func Get(sessionId string) (*Session, error) {
	return Instance.Get(sessionId)
}

func GetValue(sessionId string, key interface{}, defVal interface{}) interface{} {
	session, err := Instance.Get(sessionId)
	if err != nil {
		return defVal
	}
	return session.Get(key, defVal)
}

func SetValue(sessionId string, key interface{}, value interface{}) {
	session, err := Instance.Get(sessionId)
	if err != nil {
		return
	}
	session.Set(key, value)
}

func RemoveValue(sessionId string, key interface{}, value interface{}) error {
	session, err := Instance.Get(sessionId)
	if err != nil {
		return err
	}
	return session.Remove(key)
}

func init() {
	fmt.Printf("session init")
	// Log.Debug("session init", false)
}
