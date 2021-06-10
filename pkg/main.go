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
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		//c.mu.Lock()
		//c.routeNotes[key] = append(c.routeNotes[key], in)
		//// Note: this copy prevents blocking other clients while serving this one.
		//// We don't need to do a deep copy, because elements in the slice are
		//// insert-only and never modified.
		//rn := make([]*chat.Message, len(c.routeNotes[key]))
		//copy(rn, c.routeNotes[key])
		//c.mu.Unlock()

		c.messages = append(c.messages, in)

		for _, message := range c.messages {
			if err := stream.Send(message); err != nil {
				return err
			}
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
