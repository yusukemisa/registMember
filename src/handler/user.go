package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/yusukemisa/registMember/src/infra"

	"github.com/yusukemisa/registMember/src/user"
)

var userTemplate = template.Must(template.ParseFiles("template/user.html"))

type UserHandler struct {
	Cli *infra.Client
}

type userTemplateParams struct {
	Notice string
	Users  []*user.DispUser
}

func (h UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = h.handleGet(w, r)
	case "POST":
		err = h.handlePost(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (h UserHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	params := userTemplateParams{}
	users, err := user.GetUsers(ctx, h.Cli)
	if err != nil {
		return err
	}
	for _, v := range users {
		params.Users = append(params.Users, v.ToDispUser())
	}
	return userTemplate.Execute(w, params)
}

func (h UserHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	params := userTemplateParams{}
	userName := r.FormValue("username")
	passWord := r.FormValue("password")

	if userName == "" || passWord == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No name or message provided"
		return userTemplate.Execute(w, params)
	}

	if p, err := hash(passWord); err != nil {
		return err
	} else {
		passWord = p
	}

	u := &user.User{
		UserName:   userName,
		PassWord:   passWord,
		Status:     user.GREEN,
		Registered: time.Now(),
	}

	if err := user.RegistUser(ctx, h.Cli, u); err != nil {
		return err
	}

	users, err := user.GetUsers(ctx, h.Cli)
	if err != nil {
		return err
	}

	for _, v := range users {
		params.Users = append(params.Users, v.ToDispUser())
	}
	params.Notice = fmt.Sprintf("歓迎光臨 %s !", userName)

	return userTemplate.Execute(w, params)
}

func hash(in string) (string, error) {
	s := sha256.New()
	if _, err := s.Write([]byte(in)); err != nil {
		return "", errors.Wrap(err, "failed to write")
	}
	return hex.EncodeToString(s.Sum(nil)), nil
}
