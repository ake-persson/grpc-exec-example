package conf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mickep76/encdec"
	_ "github.com/mickep76/encdec/toml"
)

func load(fn string, c interface{}) error {
	if strings.HasPrefix(fn, "~") {
		fn = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(fn, "~"))
	}

	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		if err := encdec.FromFile("toml", fn, c); err != nil {
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

	if *printConf {
		b, _ := encdec.ToBytes("toml", c)
		fmt.Print(string(b))
		os.Exit(0)
	}
}
