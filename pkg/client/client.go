package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	chat "github.com/t-kinomura/grpc-chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func runChat(client chat.ChatClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// remote procedure call!!
	stream, err := client.SimpleChat(ctx)
	if err != nil {
		log.Fatalf("%v.SimpleChat(_) = _, %v", client, err)
	}

	// サーバーから受け取ったレスポンスを表示する.
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("Got message: %s", in.Message)
		}
	}()

	// 標準入力から文字列を受け取って送信.
	b := true
	scanner := bufio.NewScanner(os.Stdin)
	for b {
		b = scanner.Scan()
		message := chat.Message{Message: scanner.Text()}
		if err := stream.Send(&message); err != nil {
			log.Fatalf("Failed to send a message: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// サーバーと通信するためのgRPCチャンネルを作成.
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := chat.NewChatClient(conn)

	runChat(client)
}
