package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "github.com/mickep76/runshit/auth"
	"github.com/mickep76/runshit/color"
	"github.com/mickep76/runshit/ts"
)

type Token pb.Token

func (t *Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		IssuedAt  time.Time `json:"issuedAt"`
		ExpiresAt time.Time `json:"expiresAt"`
		*pb.Token
	}{
		IssuedAt:  ts.Timestamp(*t.IssuedAt).Time,
		ExpiresAt: ts.Timestamp(*t.ExpiresAt).Time,
		Token:     (*pb.Token)(t),
	})
}

func (t *Token) String() string {
	f := fmt.Sprintf("%s%%-10s%s : %s%%v%s", color.Cyan, color.Reset, color.Yellow, color.Reset)
	s := fmt.Sprintf(f, "UUID", t.Uuid)
	s += fmt.Sprintf("\n"+f, "Issued At", ts.Timestamp(*t.IssuedAt).String())
	s += fmt.Sprintf("\n"+f, "Expires At", ts.Timestamp(*t.ExpiresAt).String())
	s += fmt.Sprintf("\n"+f, "Renewed", t.Renewed)
	s += fmt.Sprintf("\n"+f, "Username", t.Username)
	s += fmt.Sprintf("\n"+f, "Name", t.Name)
	s += fmt.Sprintf("\n"+f, "Mail", t.Mail)
	s += fmt.Sprintf("\n"+f, "Roles", strings.Join(t.Roles, ","))
	s += "\n\n"
	return s
}
