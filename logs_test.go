package dbase

import "testing"

func Test_Qlog(t *testing.T) {
	QLog("Hello world")
	SetQLogFolder("testdata/loggy")
	QLog("Hello now")
}

type tlogger struct {
	logs []string
}

func (t *tlogger) Log(m string) {
	t.logs = append(t.logs, m)
}

func Test_SetQLogger(t *testing.T) {
	tlog := &tlogger{}
	SetQLogger(tlog)
	QLog("Hello")
	QLog("Goodbye")
	if len(tlog.logs) != 2 {
		t.Log("logs should have 2 messages")
		t.Fail()
	}
}
