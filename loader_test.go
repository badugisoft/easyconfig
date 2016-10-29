package easyconfig_test

import (
	"os"
	"testing"

	"github.com/badugisoft/easyconfig"
	"github.com/badugisoft/easyconfig/test"
	"github.com/badugisoft/xson"
)

func TestLoadDir(t *testing.T) {
	cfg := test.ConfigData{}
	err := easyconfig.LoadDir(&cfg, "dev", "test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = xson.Marshal(xson.YAML, cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadAsset(t *testing.T) {
	cfg := test.ConfigData{}
	err := easyconfig.LoadAsset(&cfg, "dev", test.Asset, "test/")
	if err != nil {
		t.Fatal(err)
	}

	_, err = xson.Marshal(xson.YAML, cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("cfg_mode_position", "env")
	os.Setenv("cfg_partial_sub_two", "223344")

	cfg := test.ConfigData{}
	err := easyconfig.LoadEnv(&cfg, "cfg")
	if err != nil {
		t.Fatal(err)
	}

	_, err = xson.Marshal(xson.YAML, cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadArg(t *testing.T) {
	os.Args = []string{
		"load_test",
		"--cfg",
		"mode.position=arg",
		"--cfg",
		"partial.sub.two=223344",
	}
	os.Setenv("cfg_mode_position", "env")
	os.Setenv("cfg_partial_sub_two", "223344")

	cfg := test.ConfigData{}
	err := easyconfig.LoadArg(&cfg, "cfg")
	if err != nil {
		t.Fatal(err)
	}

	_, err = xson.Marshal(xson.YAML, cfg)
	if err != nil {
		t.Fatal(err)
	}
}
