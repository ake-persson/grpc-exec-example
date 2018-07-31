package auth

import (
	"fmt"
	"strings"

	"github.com/mickep76/grpc-exec-example/color"
	"github.com/mickep76/grpc-exec-example/ts"
)

func (t *Token) FmtStringColor() string {
	f := fmt.Sprintf("%s%%-24s%s : %s%%v%s\n", color.Cyan, color.Reset, color.Yellow, color.Reset)

	txt := fmt.Sprintf(f, "UUID", t.Uuid)
	txt += fmt.Sprintf(f, "Issued At", ts.Timestamp(*t.IssuedAt))
	txt += fmt.Sprintf(f, "Expires At", ts.Timestamp(*t.ExpiresAt))
	txt += fmt.Sprintf(f, "Renewed", t.Renewed)
	txt += fmt.Sprintf(f, "Username", t.Username)
	txt += fmt.Sprintf(f, "Name", t.Name)
	txt += fmt.Sprintf(f, "E-Mail", t.Mail)
	if len(t.Roles) > 0 {
		txt += fmt.Sprintf(f, "Roles", strings.Join(t.Roles, ", "))
	}
	txt += "\n"
	return txt
}
