package client

import (
	"context"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client/state"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
