package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"io"
	"log"
	"net"
	"sync"
	"time"

	chat "github.com/t-kinomura/grpc-chat"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
)

type chatServer struct {
	chat.UnimplementedChatServer

	mu sync.Mutex // protects routeNotes
	// TODO
	//routeNotes map[string][]*chat.Message
	messages []*chat.Message
}

func (c *chatServer) SimpleChat(stream chat.Chat_SimpleChatServer) error {
	go SendChat(stream, c)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			c.messages = nil
			return nil
		}
		if err != nil {
			return err
		}

		c.messages = append(c.messages, in)
		fmt.Printf("Recieved Message: %s\n", c.messages[len(c.messages) - 1])
	}
}

func SendChat(stream chat.Chat_SimpleChatServer, c *chatServer) error {
	currentCount := 0
	for {
		if currentCount < len(c.messages) {
			if err := stream.Send(c.messages[currentCount]); err != nil {
				return err
			}

			time.Sleep(100 * time.Microsecond)
			currentCount++
		}
	}
}

func newServer() *chatServer {
	s := &chatServer{messages: make([]*chat.Message, 0)}
	return s
}

func main() {
	flag.Parse()

	// クライアントのリクエストをリッスンするために使用するポートを指定.
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// gRPCサーバーのインスタンスを作成.
	grpcServer := grpc.NewServer(opts...)

	// サービスの実装をgRPCサーバに登録.
	// grpc.pb.go
	chat.RegisterChatServer(grpcServer, newServer())

	// サーバーのServe()をポートの詳細とともに呼び出す.
	// プロセスがキルされるかStop()が呼ばれるまでブロッキング待機.
	grpcServer.Serve(lis)
}
