package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/badugisoft/xson"
)

func loadValue(v interface{}, names []string, str string) error {
	buf := bytes.NewBuffer([]byte{})

	for l, name := range names {
		fmt.Fprint(buf, "\n")
		for i := 0; i < l; i++ {
			fmt.Fprint(buf, "  ")
		}
		fmt.Fprintf(buf, "%v:", name)
	}

	fmt.Fprintf(buf, " %v", str)

	err := xson.Unmarshal(xson.YAML, buf.Bytes(), v)
	if err != nil {
		return err
	}

	return nil
}

func LoadDir(v interface{}, mode, dir string) error {
	for _, name := range []string{"default", mode, "local"} {
		for _, t := range xson.GetTypes() {
			for _, ext := range xson.GetExtensions(t) {
				d, err := ioutil.ReadFile(dir + "/" + name + "." + ext)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return err
				}
				err = xson.Unmarshal(t, d, v)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type AssetFunc func(name string) ([]byte, error)

func LoadAsset(v interface{}, mode string, f AssetFunc, pathPrefix string) error {
	for _, name := range []string{"default", mode, "local"} {
		for _, t := range xson.GetTypes() {
			for _, ext := range xson.GetExtensions(t) {
				d, err := f(pathPrefix + name + "." + ext)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						continue
					}
					return err
				}

				err = xson.Unmarshal(t, d, v)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func LoadEnv(v interface{}, prefix string) error {
	if !strings.HasSuffix(prefix, "_") {
		prefix = prefix + "_"
	}

	prefixLen := len(prefix)

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		tokens := strings.SplitN(env[prefixLen:], "=", 2)
		if len(tokens) != 2 {
			continue
		}

		err := loadValue(v, strings.Split(tokens[0], "_"), tokens[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadArg(v interface{}, flag string) error {
	args := os.Args
	flagStr := "--" + flag

	for i, e := 0, len(args)-1; i < e; i++ {
		if args[i] == flagStr {
			tokens := strings.SplitN(args[i+1], "=", 2)
			if len(tokens) != 2 {
				continue
			}

			err := loadValue(v, strings.Split(tokens[0], "."), tokens[1])
			if err != nil {
				return err
			}
			i++
		}
	}
	return nil
}
