package dbase2

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type Logger struct {
	Folder string
	sync.Mutex
}

var dLog = Logger{
	Folder: "logs",
}

func DLog(m string) {
	dLog.Log(m)
}

func (l *Logger) Log(m string) {
	l.Lock()
	defer l.Unlock()
	now := time.Now()
	fname := now.Format("060102")
	p := path.Join(l.Folder, fname+".log")

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("message not logged : ", err, "::", m)
		return
	}
	defer f.Close()

	line := now.Format("15:04:05") + "::" + m + "\n"
	_, err = f.WriteString(line)
	if err != nil {
		fmt.Println("message not logged: ", err, "::", m)
	}

}

type LogGroup struct {
	Main Logger
	logs map[string]Logger
	sync.Mutex
}
