package dbase2

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type SessionData struct {
	Data       interface{}
	LastAccess time.Time
}

type SessionControl struct {
	sessions map[string]SessionData
	sync.Mutex
	lastSweep   time.Time
	MaxDuration time.Duration
}

func NewSessionControl(md time.Duration) *SessionControl {
	return &SessionControl{
		sessions:    make(map[string]SessionData),
		lastSweep:   time.Now(),
		MaxDuration: md,
	}
}

func (sc *SessionControl) Login(w http.ResponseWriter, data interface{}) {
	dt := SessionData{
		Data:       data,
		LastAccess: time.Now(),
	}
	sc.Lock()
	//Get ID
	id := fmt.Sprintf("%X", rand.Uint64())
	_, ok := sc.sessions[id]
	for ok {
		id := fmt.Sprintf("%X", rand.Uint64())
		_, ok = sc.sessions[id]
	}
	sc.sessions[id] = dt
	//Write about it
	http.SetCookie(w, &http.Cookie{
		Name:  "Session",
		Value: id,
	})
	sc.Unlock()
	go sc.Sweep()
}

const (
	OK = iota
	TIMEOUT
	NOLOGIN
)

func (sc *SessionControl) GetLogin(w http.ResponseWriter, r *http.Request) (SessionData, int) {
	c, err := r.Cookie("Session")
	if err != nil {
		return SessionData{}, NOLOGIN
	}
	sc.Lock()
	defer sc.Unlock()
	dt, ok := sc.sessions[c.Value]
	if !ok {
		return dt, NOLOGIN
	}
	if time.Now().After(dt.LastAccess.Add(sc.MaxDuration)) {
		http.SetCookie(w, &http.Cookie{
			Name:    "Session",
			Value:   "None",
			Expires: time.Now().Add(-time.Second),
		})
		return SessionData{}, TIMEOUT

	}
	dt.LastAccess = time.Now()
	sc.sessions[c.Value] = dt
	return dt, OK
}

func (sc *SessionControl) EditLogin(r *http.Request, data interface{}) err {
	c, err := r.Cookie("Session")
	if err != nil {
		return errors.New("No Login cookie")
	}
	sc.Lock()
	defer sc.Unlock()

	sdat, ok := sc.sessions[c.Value]
	sdat.Data = data
	sdat.LastAccess = time.Now()
	return nil
}

func (sc *SessionControl) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("Session")
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "Session",
		Value:   "None",
		Expires: time.Now().Add(-time.Second),
	})
	sc.Lock()
	delete(sc.sessions, c.Value)
	sc.Unlock()
}

func (sc *SessionControl) Sweep() {
	if sc.lastSweep.Add(sc.MaxDuration * 10).After(time.Now()) {
		return
	}
	sc.lastSweep = time.Now()
	t := time.Now().Add(-sc.MaxDuration)
	sc.Lock()
	for k, v := range sc.sessions {
		if t.After(v.LastAccess) {
			delete(sc.sessions, k)
		}
	}
	sc.Unlock()

}
