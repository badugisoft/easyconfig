package easyconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/badugisoft/xson"
)

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

func MustLoadDir(v interface{}, mode, dir string) {
	err := LoadDir(v, mode, dir)
	if err != nil {
		panic(err)
	}
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

func MustLoadAsset(v interface{}, mode string, f AssetFunc, pathPrefix string) {
	err := LoadAsset(v, mode, f, pathPrefix)
	if err != nil {
		panic(err)
	}
}

func LoadEnv(v interface{}, prefix string) error {
	if !strings.HasSuffix(prefix, "_") {
		prefix = prefix + "_"
	}

	prefixLen := len(prefix)
	buf := bytes.NewBuffer([]byte{})

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		tokens := strings.SplitN(env[prefixLen:], "=", 2)
		if len(tokens) != 2 {
			continue
		}

		fmt.Fprintf(buf, "%v: %v\n", strings.Replace(tokens[0], "_", ".", -1), tokens[1])
	}
	return xson.Unmarshal(xson.FLAT_YAML, buf.Bytes(), v)
}

func MustLoadEnv(v interface{}, prefix string) {
	err := LoadEnv(v, prefix)
	if err != nil {
		panic(err)
	}
}

func LoadArg(v interface{}, flag string) error {
	args := os.Args
	flagStr := "--" + flag
	buf := bytes.NewBuffer([]byte{})

	for i, e := 0, len(args)-1; i < e; i++ {
		if args[i] == flagStr {
			tokens := strings.SplitN(args[i+1], "=", 2)
			if len(tokens) != 2 {
				continue
			}

			fmt.Fprintf(buf, "%v: %v\n", tokens[0], tokens[1])
			i++
		}
	}
	return xson.Unmarshal(xson.FLAT_YAML, buf.Bytes(), v)
}

func MustLoadArg(v interface{}, prefix string) {
	err := LoadArg(v, prefix)
	if err != nil {
		panic(err)
	}
}
