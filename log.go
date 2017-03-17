package dbase2

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Logger interface {
	Log(string)
}

type FmtLog struct{}

func (fl FmtLog) Log(m string) {
	fmt.Println(m)
}

type FolderLog struct {
	Folder string
	sync.Mutex
}

var dlog Logger

// QLog provides a quick log method, that requieres no setup to allow simple logging
// If you want more fancy stuff, Use a LogGroup
func QLog(m string) error {
	if dlog == nil {
		var err error
		dlog, err = NewLogGroup("logs")
		if err != nil {
			return errors.Wrap(err, "Could not log")
		}
	}
	dlog.Log(m)
	return nil
}

// SetQLogFolder changes the default folder for QLog default "logs"
// returns an error if it cannot prepare the directory
func SetQLogFolder(f string) error {
	nlog, err := NewLogGroup(f)
	if err != nil {
		return errors.Wrap(err, "Could not set Log Folder")
	}
	dlog = nlog
	return nil
}

// SetQLogger sets the default QLog method to use anything that fits the Log method, useful for quick debog options, etc
func SetQLogger(l Logger) {
	dlog = l
}

func NewFolderLog(f string) (*FolderLog, error) {

	err := os.MkdirAll(f, 0774)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create Logger")
	}

	return &FolderLog{Folder: f}, nil

}

func (l *FolderLog) Log(m string) {
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
	main Logger
	logs map[string]Logger
	sync.Mutex
}

func NewLogGroup(f string) (*LogGroup, error) {
	var main Logger
	main, err := NewFolderLog(f)

	if err != nil {
		main = FmtLog{}
	}

	return &LogGroup{
		main: main,
		logs: make(map[string]Logger),
	}, nil
}
func (lg *LogGroup) SetMain(l Logger) {
	lg.Lock()
	lg.main = l
	lg.Unlock()
}

func (lg *LogGroup) AddLogger(k string, l Logger) {
	lg.Lock()
	lg.logs[k] = l
	lg.Unlock()

}

func (lg *LogGroup) AddFolderLog(k, fol string) error {
	lg.Lock()
	fl, err := NewFolderLog(fol)
	if err != nil {
		return errors.Wrap(err, "Could not Add folder log")
	}
	lg.logs[k] = fl
	lg.Unlock()
	return nil
}

func (lg *LogGroup) Log(m string) {
	lg.main.Log(m)
}

func (lg *LogGroup) LogTo(k, m string) {
	lg.Lock()
	defer lg.Unlock()

	go lg.main.Log(k + m)
	mini, ok := lg.logs[k]
	if ok {
		go mini.Log(m)
	}
}
