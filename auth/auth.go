package auth

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/runshit/tlscfg"
)

func GetPublicKey(addr string, ca string, insecure bool) ([]byte, error) {
	conf, err := tlscfg.NewConfig(ca, "", "", "", insecure)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(conf)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	auth := NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	key, err := auth.GetPublicKey(ctx, &Empty{})
	if err != nil {
		return nil, err
	}

	return key.Pem, nil
}
