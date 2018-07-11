package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/mickep76/auth/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb_auth "github.com/mickep76/grpc-exec-example/auth"
	"github.com/mickep76/grpc-exec-example/cmd"
	"github.com/mickep76/grpc-exec-example/conf"
	pb_exec "github.com/mickep76/grpc-exec-example/exec"
	"github.com/mickep76/grpc-exec-example/ts"
)

type server struct {
	jwt *jwt.JWTClient
}

func newServer(auth string, ca string) (*server, error) {
	// Get public key from authentication service.
	k, err := pb_auth.GetPublicKey(auth, ca, false)
	if err != nil {
		return nil, err
	}

	s := &server{}

	// Create new jwt client.
	if s.jwt, err = jwt.NewJWTClient(jwt.WithPublicKey(k)); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *server) Exec(in *pb_exec.Command, stream pb_exec.ExecCommand_ExecServer) error {
	if err := s.jwt.AuthorizedGrpc(stream.Context()); err != nil {
		log.Print(err)
		return err
	}

	c := cmd.New(in.Command, in.Arguments, cmd.WithEnv(in.Environment), cmd.WithUser(in.User), cmd.WithGroup(in.Group), cmd.WithDir(in.Directory), cmd.WithRedirect, cmd.WithTimeout(5*time.Second))
	defer c.Kill()

	log.Printf("%s start %s", in.Uuid, c)
	if _, err := c.Start(); err != nil {
		log.Printf("%s %v", in.Uuid, err)
		return err
	}

	line := uint32(0)
	scanner := bufio.NewScanner(c.Stdout())
	for scanner.Scan() {
		now := ts.Now().Timestamp()
		m := &pb_exec.Message{Timestamp: &now, Line: line, Message: scanner.Text()}

		if err := stream.Send(m); err != nil {
			log.Printf("%s %v", in.Uuid, err)
			return err
		}
		line++
	}
	if scanner.Err() != io.EOF && scanner.Err() != nil {
		log.Printf("%s %v", in.Uuid, scanner.Err())
		return scanner.Err()
	}

	if _, err := c.Wait(); err != nil {
		log.Printf("%s %v", in.Uuid, err)
		return err
	}

	log.Printf("%s finished", in.Uuid)

	return nil
}

func main() {
	c := newConfig()
	if err := conf.Load([]string{"/etc/exec-server.toml", "~/.exec-server.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, os.Args[1:], c)

	creds, err := credentials.NewServerTLSFromFile(c.Cert, c.Key)
	if err != nil {
		log.Fatalf("credentials: %v", err)
	}

	lis, err := net.Listen("tcp", c.Bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srvr, err := newServer(c.Auth, c.Ca)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb_exec.RegisterExecCommandServer(s, srvr)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
