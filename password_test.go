package dbase2

import "testing"

func Test_pw(t *testing.T) {
	pw, _ := NewPassword("Hello")

	a := pw.Check("Hello")

	if !a {
		t.Log("a should be true")
		t.Fail()
	}
	b := pw.Check("Helli")
	if b {
		t.Log("b should be false")
		t.Fail()
	}

}
