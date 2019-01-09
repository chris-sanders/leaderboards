package cfg

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var config *Config

func TestNewConfg(t *testing.T) {
	logFile, err := os.OpenFile("leaderboards.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	config, err = (&Config{}).New()
	if err != nil {
		config, err = (&Config{}).New()
		if err != nil {
			t.Error(err)
		}
	}
	expect := 10
	got := config.Import.Local_limit
	if got != expect {
		t.Errorf("Wrong default value %v, expected %v", got, expect)
	}
	expect = 1
	got = config.Import.Remote_limit
	if got != expect {
		t.Errorf("Wrong default value %v, expected %v", got, expect)
	}
	expects := "all"
	gots := config.Sync.Accounts
	if gots != expects {
		t.Errorf("Wrong default value %v, expected %v", gots, expects)
	}
	expects = "sync"
	gots = config.Sync.Folder
	if gots != expects {
		t.Errorf("Wrong default value %v, expected %v", gots, expects)
	}
	expects = "INFO"
	gots = config.Global.Logging
	if gots != expects {
		t.Errorf("Wrong default value %v, expected %v", gots, expects)
	}
	expects = "APPDATA"
	gots = config.Global.Local_folder
	if gots != expects {
		t.Errorf("Wrong default value %v, expected %v", gots, expects)
	}
	expect = 0
	info, err := os.Stat("leaderboards.log")
	got = int(info.Size())
	if got != expect {
		t.Errorf("Log file exected %v bytes, got %v bytes", expect, got)
	}
	err = os.Remove("leaderboards.log")
	if err != nil {
		t.Errorf("Failed to remove log file: %v", err)
	}
	err = os.Remove("settings.yaml")
	if err != nil {
		t.Errorf("Failed to remove log file: %v", err)
	}
}

func TestPostProcessing(t *testing.T) {
	expect := "LocalLeaderboards.dat"
	got := config.Global.Local_file
	if got != expect {
		t.Errorf("Expected local_file %v, got %v", expect, got)
	}

}
