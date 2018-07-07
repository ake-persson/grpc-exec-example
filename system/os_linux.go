// +build linux

package system

import (
	"fmt"
	"strings"

	pb_info "github.com/mickep76/runshit/info"
)

func getOS(s *pb_info.System) error {
	descr, err := readFile("/etc/redhat-release")
	if err != nil {
		return err
	}

	s.OsDescription = strings.TrimSpace(string(descr))

	a := strings.SplitN(s.OsDescription, " release ", 2)
	if len(a) != 2 {
		return fmt.Errorf("unknown format in [%s]: %s", "/etc/redhat-release", s.OsDescription)
	}

	s.Os = strings.Replace(a[0], " Linux", "", 1)

	a = strings.SplitN(a[1], " ", 2)
	if len(a) != 2 {
		return fmt.Errorf("unknown format in [%s]: %s", "/etc/redhat-release", s.OsDescription)
	}

	s.OsVersion = a[0]
	s.OsBuild = a[1][1 : len(a[1])-1]

	return nil
}
