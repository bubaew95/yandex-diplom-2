package main

import (
	"context"
	"fmt"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":3232", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewGoKeeperClient(conn)

	text := pb.TextRequest{
		Text: "Hello World",
	}

	res, err := c.AddText(context.Background(), &text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
