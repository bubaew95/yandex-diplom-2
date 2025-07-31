package client

import (
	"context"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=GoKeeperClient --filename=servermock_test.go --inpackage
type GoKeeperClient interface {
	Registration(ctx context.Context, in *pb.RegistrationRequest, opts ...grpc.CallOption) (*pb.TokenResponse, error)
	Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.TokenResponse, error)
	Add(ctx context.Context, in *pb.DataRequest, opts ...grpc.CallOption) (*pb.DataResponse, error)
	Edit(ctx context.Context, in *pb.DataEditRequest, opts ...grpc.CallOption) (*pb.DataResponse, error)
	Delete(ctx context.Context, in *pb.IdRequest, opts ...grpc.CallOption) (*pb.SuccessResponse, error)
	FindAll(ctx context.Context, in *pb.EmptyRequest, opts ...grpc.CallOption) (*pb.DataList, error)
}

// State представляет состояние авторизованного пользователя.
//
// Содержит:
//   - Token: строка авторизации (JWT или иная),
//   - User: объект пользователя.
type State struct {
	Token string
	User  model.User
}

type Client struct {
	KeeperClient GoKeeperClient
	State        *State
	Conn         *grpc.ClientConn
}

func NewClient(cfg *config.Config) (*Client, error) {
	conn, err := grpc.NewClient(":"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewGoKeeperClient(conn)

	return &Client{
		Conn:         conn,
		KeeperClient: c,
		State: &State{
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

func (c *Client) GetAllData(ctx context.Context) (*pb.DataList, error) {
	nCtx := c.authorizationToken(ctx)
	data, err := c.KeeperClient.FindAll(nCtx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Add(ctx context.Context, data *model.Data) (bool, error) {
	nCtx := c.authorizationToken(ctx)

	_, err := c.KeeperClient.Add(nCtx, &pb.DataRequest{
		Text: data.Text,
		Type: string(data.Type),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) Edit(ctx context.Context, ID int64, data *model.Data) (bool, error) {
	nCtx := c.authorizationToken(ctx)

	_, err := c.KeeperClient.Edit(nCtx, &pb.DataEditRequest{
		Id:   ID,
		Text: data.Text,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) Delete(ctx context.Context, id *pb.IdRequest) (bool, error) {
	nCtx := c.authorizationToken(ctx)

	res, err := c.KeeperClient.Delete(nCtx, id)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}

func (c *Client) Stop() error {
	return c.Conn.Close()
}
