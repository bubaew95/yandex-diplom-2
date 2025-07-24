package main

import (
	"context"
	"fmt"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	//token, err := c.Registration(context.Background(), &pb.RegistrationRequest{
	//	FirstName:  "John",
	//	LastName:   "Doe",
	//	Email:      "sdfdhfsss@doe.com",
	//	Password:   "123456",
	//	RePassword: "123456",
	//})

	token, err := c.Login(context.Background(), &pb.LoginRequest{
		Email:    "sdfdhfsss@doe.com",
		Password: "123456",
	})

	if err != nil {
		log.Fatal(err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"token": fmt.Sprintf(`Bearer %s`, token.Token),
	}))

	//textEdit := pb.TextEditRequest{
	//	Id:   15,
	//	Text: "test pro",
	//}
	//res, err := c.EditText(ctx, &textEdit)

	//text := pb.TextRequest{
	//	Text: "Hello World: ",
	//}
	//_, err = c.AddText(ctx, &text)

	//ID := 15
	//res, err := c.DeleteText(ctx, &pb.IdRequest{Id: int64(ID)})

	res, err := c.FindAllText(ctx, &pb.DataRequest{})

	if err != nil {
		log.Println(res)
	}

	for _, item := range res.List {
		fmt.Println(item)
	}

	fmt.Println(res)
}
