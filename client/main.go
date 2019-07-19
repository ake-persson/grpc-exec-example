package main

import (
	"fmt"
	"os"

	"github.com/mickep76/grpc-exec-example/client/exec"
	"github.com/mickep76/grpc-exec-example/client/info"
	"github.com/mickep76/grpc-exec-example/client/list"
	"github.com/mickep76/grpc-exec-example/client/login"
	"github.com/mickep76/grpc-exec-example/client/renew"
	"github.com/mickep76/grpc-exec-example/client/verify"
)

func usage() {
	fmt.Print(`Usage of runshit [login|verify|renew|info|exec|list] [-h]:

Commands:
  login
        Login for a JWT token.
  verify
        Verify JWT token.
  renew
        Renew JWT token.
  info
        Information about host.
  register
        Register host.
  exec
        Execute command.
  list
        List servers registered in catalog.
`)
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "login":
		login.Cmd(os.Args[2:])
	case "verify":
		verify.Cmd(os.Args[2:])
	case "renew":
		renew.Cmd(os.Args[2:])
	case "info":
		info.Cmd(os.Args[2:])
	case "list":
		list.Cmd(os.Args[2:])
	case "exec":
		exec.Cmd(os.Args[2:])
	default:
		usage()
	}
}
