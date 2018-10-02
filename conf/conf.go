package conf

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mickep76/encoding"
	_ "github.com/mickep76/encoding/toml"
)

func load(fn string, c interface{}) error {
	if strings.HasPrefix(fn, "~") {
		fn = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(fn, "~"))
	}

	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}

		codec, _ := encoding.NewCodec("toml")
		if err := codec.Decode(b, c); err != nil {
			return err
		}
	}

	return nil
}

func Load(files []string, c interface{}) error {
	for _, f := range files {
		if err := load(f, c); err != nil {
			return err
		}
	}
	return nil
}

func ParseFlags(fl *flag.FlagSet, args []string, c interface{}) {
	printConf := fl.Bool("print-conf", false, "Print config.")
	fl.Parse(args)

	codec, _ := encoding.NewCodec("toml")
	if *printConf {
		b, _ := codec.Encode(c)
		fmt.Print(string(b))
		os.Exit(0)
	}
}
