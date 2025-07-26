package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client/state"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	ErrInvalidResponse  = errors.New("неверный ответ сервера")
	ErrNotAuthenticated = errors.New("не авторизован")
	ErrAuthFailed       = errors.New("аутентификация не удалась")
	ErrDataNotFound     = errors.New("данные не найдены")
)

const (
	httpErrorCodeStart     = 400
	masterPasswordFileName = "master.key"
	configDir              = ".gophkeeper"
	clientTimeoutSeconds   = 10
)

type Client struct {
	State        *state.State
	KeeperClient pb.GoKeeperClient
}

func NewClient(cfg *config.Config) (*Client, error) {
	conn, err := grpc.NewClient(":"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	//defer conn.Close()
	c := pb.NewGoKeeperClient(conn)

	return &Client{
		KeeperClient: c,
		State: &state.State{
			Token: "",
			User:  model.User{},
		},
	}, nil
}

func (c *Client) Login(ctx context.Context, email string, password string) (string, error) {
	token, err := c.KeeperClient.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func (c *Client) authorizationToken(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"token": fmt.Sprintf(`Bearer %s`, c.State.Token),
	}))
}

func (c *Client) GetAllData(ctx context.Context) (*pb.TextList, error) {
	nCtx := c.authorizationToken(ctx)
	data, err := c.KeeperClient.FindAllText(nCtx, &pb.DataRequest{})
	if err != nil {
		return nil, err
	}

	return data, nil
}
