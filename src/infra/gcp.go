package infra

import (
	"context"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
)

type Client struct {
	Datastore datastore.Client
}

func NewDatastoreClient(ctx context.Context, opts ...datastore.ClientOption) (*Client, error) {
	client, err := clouddatastore.FromContext(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		Datastore: client,
	}, nil
}

func (cli *Client) SaveWithTransaction(ctx context.Context, putch func(tx datastore.Transaction) error) error {
	_, err := cli.Datastore.RunInTransaction(ctx, putch)
	return err
}
