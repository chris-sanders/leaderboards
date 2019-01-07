package cfg

import (
	_ "fmt"
	"testing"
)

func TestNewConfg(t *testing.T) {
	config, err := (&Config{}).New()
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
}
