package user

import (
	"context"
	"log"
	"time"

	"github.com/yusukemisa/registMember/src/infra"

	"go.mercari.io/datastore"
)

type User struct {
	UserName   string
	PassWord   string
	Status     Status
	Registered time.Time
}

type DispUser struct {
	UserName string
	Status   string
}

func (u *User) ToDispUser() *DispUser {
	return &DispUser{
		UserName: u.UserName,
		Status:   StatusMap[u.Status],
	}
}

type Status int32

const (
	ANONYMOUS Status = 0
	GREEN     Status = 1
	YELLOW    Status = 2
	RED       Status = 3
	BLACK     Status = 10
)

var StatusMap = map[Status]string{
	ANONYMOUS: "ANONYMOUS",
	GREEN:     "GREEN",
	YELLOW:    "YELLOW",
	RED:       "RED",
	BLACK:     "BLACK",
}

func (u *User) validate() error {
	return nil
}

func RegistUser(ctx context.Context, cli *infra.Client, u *User) error {
	err := cli.SaveWithTransaction(ctx, func(tx datastore.Transaction) error {
		if err := u.validate(); err != nil {
			return err
		}
		key := cli.Datastore.NameKey("User", u.UserName, nil)
		if _, err := tx.Put(key, u); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func GetUsers(ctx context.Context, cli *infra.Client) ([]*User, error) {
	var users []*User
	query := cli.Datastore.NewQuery("User").Order("Registered")
	if _, err := cli.Datastore.GetAll(ctx, query, &users); err != nil {
		return nil, err
	}
	return users, nil
}
