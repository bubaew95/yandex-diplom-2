package main

import (
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client/pages"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := logger.Load(); err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("No .env file found")
	}
}

func main() {
	cfg := config.NewConfig()
	tui, err := pages.NewTUI(cfg)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	if err := tui.Run(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}

/**

// устанавливаем соединение с сервером
//conn, err := grpc.NewClient(":"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
//if err != nil {
//	log.Fatal(err)
//}
//defer conn.Close()
//c := pb.NewGoKeeperClient(conn)
//
////token, err := c.Registration(context.Background(), &pb.RegistrationRequest{
////	FirstName:  "John",
////	LastName:   "Doe",
////	Email:      "sdfdhfsss@doe.com",
////	Password:   "123456",
////	RePassword: "123456",
////})
//
//token, err := c.Login(context.Background(), &pb.LoginRequest{
//	Email:    "sdfdhfsss@doe.com",
//	Password: "123456",
//})
//
//if err != nil {
//	log.Fatal(err)
//}
//
//ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
//	"token": fmt.Sprintf(`Bearer %s`, token.Token),
//}))
//
////textEdit := pb.TextEditRequest{
////	Id:   15,
////	Text: "test pro",
////}
////res, err := c.EditText(ctx, &textEdit)
//
////text := pb.TextRequest{
////	Text: "Hello World: ",
////}
////_, err = c.Add(ctx, &text)
//
////ID := 15
////res, err := c.DeleteText(ctx, &pb.IdRequest{Id: int64(ID)})
//
//res, err := c.FindAllText(ctx, &pb.DataRequest{})
//
//if err != nil {
//	log.Println(res)
//}
//
//for _, item := range res.List {
//	fmt.Println(item)
//}
//
//fmt.Println(res)
*/
