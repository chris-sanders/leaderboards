package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	Log.SetLevel(Log.TraceLevel)
	Log.Info("Info logged")
	Log.Warn("Warn logged")
	Log.Close()
	byteSlice, err := ioutil.ReadFile("leaderboards.log")
	if err != nil {
		t.Error(err)
	}
	expect := "Info logged"
	if !strings.Contains(string(byteSlice), expect) {
		t.Errorf("Couldn't find log mesage: %v", expect)
	}
	expect = "Warn logged"
	if !strings.Contains(string(byteSlice), expect) {
		t.Errorf("Couldn't find log mesage: %v", expect)
	}
	os.Remove("leaderboards.log")
}
